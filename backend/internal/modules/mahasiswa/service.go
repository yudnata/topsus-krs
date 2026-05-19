package mahasiswa

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"backend/internal/cache"
)

type Service struct {
	repo *Repository
	inv  *cache.Invalidator
}

func NewService(repo *Repository, inv *cache.Invalidator) *Service {
	return &Service{repo: repo, inv: inv}
}

// ─── GetCurrentKRS ────────────────────────────────────────────────────────────

// GetCurrentKRS — ambil/buat draft KRS untuk TA aktif, sertakan detail kelas.
func (s *Service) GetCurrentKRS(ctx context.Context, userID string) (*KrsHeader, error) {
	profile, err := s.repo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("profil mahasiswa tidak ditemukan")
	}
	taID, kodeTa, semester, err := s.repo.GetTahunAkademikAktif(ctx)
	if err != nil {
		return nil, errors.New("tidak ada tahun akademik aktif")
	}
	krs, err := s.repo.GetOrCreateDraft(ctx, profile.ID, taID)
	if err != nil {
		return nil, err
	}
	krs.KodeTa = kodeTa
	krs.Semester = semester
	krs.Detail, _ = s.repo.GetDetails(ctx, krs.ID)
	return krs, nil
}

// ─── AddClass ─────────────────────────────────────────────────────────────────

// AddClass — tambah kelas ke draft KRS dengan 5 validasi bisnis.
func (s *Service) AddClass(ctx context.Context, userID, krsID, kelasID string) (*KrsHeader, error) {
	// Pastikan KRS milik mahasiswa ini dan masih draft
	krs, err := s.repo.GetKrsByID(ctx, krsID)
	if err != nil {
		return nil, errors.New("KRS tidak ditemukan")
	}
	profile, err := s.repo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("profil mahasiswa tidak ditemukan")
	}
	if krs.MahasiswaID != profile.ID {
		return nil, errors.New("akses ditolak: bukan KRS Anda")
	}

	// [V5] Kunci KRS jika bukan draft
	if krs.Status != "draft" {
		return nil, fmt.Errorf("KRS tidak dapat diedit (status: %s)", krs.Status)
	}

	// Ambil detail kelas yang ingin ditambahkan
	newKelas, err := s.repo.GetKelasDetail(ctx, kelasID)
	if err != nil {
		return nil, errors.New("kelas tidak ditemukan")
	}

	// [V1] Kapasitas kelas
	if newKelas.Terisi >= newKelas.Kapasitas {
		return nil, errors.New("kapasitas kelas sudah penuh")
	}

	// Ambil kelas yang sudah ada di KRS
	existing, err := s.repo.GetCurrentKelasInKrs(ctx, krsID)
	if err != nil {
		return nil, err
	}

	// [V4] Duplikat kelas / MK
	for _, ex := range existing {
		if ex.ID == kelasID {
			return nil, errors.New("kelas ini sudah ada di KRS Anda")
		}
		if ex.MataKuliahID == newKelas.MataKuliahID {
			return nil, errors.New("mata kuliah ini sudah ada di KRS Anda (kelas berbeda)")
		}
	}

	// [V3] Max SKS
	if krs.TotalSKS+newKelas.SKS > profile.MaxSKS {
		return nil, fmt.Errorf("melebihi batas SKS (%d/%d)", krs.TotalSKS+newKelas.SKS, profile.MaxSKS)
	}

	// [V2] Bentrok jadwal — cek overlap waktu
	if err := checkScheduleConflict(existing, newKelas); err != nil {
		return nil, err
	}

	// Tambahkan ke DB
	if err := s.repo.AddKelas(ctx, krsID, kelasID, newKelas.SKS); err != nil {
		return nil, err
	}

	// Invalidate cache
	dpID := ""
	if profile.DosenPembimbingID != nil {
		dpID = *profile.DosenPembimbingID
	}
	s.inv.OnKrsMutation(ctx, profile.ID, dpID)

	// Return KRS terbaru
	return s.getKrsWithDetail(ctx, krsID)
}

// ─── RemoveClass ──────────────────────────────────────────────────────────────

func (s *Service) RemoveClass(ctx context.Context, userID, krsID, kelasID string) (*KrsHeader, error) {
	krs, err := s.repo.GetKrsByID(ctx, krsID)
	if err != nil {
		return nil, errors.New("KRS tidak ditemukan")
	}
	profile, err := s.repo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("profil mahasiswa tidak ditemukan")
	}
	if krs.MahasiswaID != profile.ID {
		return nil, errors.New("akses ditolak: bukan KRS Anda")
	}
	if krs.Status != "draft" {
		return nil, fmt.Errorf("KRS tidak dapat diedit (status: %s)", krs.Status)
	}

	kd, err := s.repo.GetKelasDetail(ctx, kelasID)
	if err != nil {
		return nil, errors.New("kelas tidak ditemukan")
	}

	if err := s.repo.RemoveKelas(ctx, krsID, kelasID, kd.SKS); err != nil {
		return nil, errors.New("kelas tidak ada di KRS Anda")
	}

	dpID := ""
	if profile.DosenPembimbingID != nil {
		dpID = *profile.DosenPembimbingID
	}
	s.inv.OnKrsMutation(ctx, profile.ID, dpID)

	return s.getKrsWithDetail(ctx, krsID)
}

// ─── SubmitKRS ────────────────────────────────────────────────────────────────

func (s *Service) SubmitKRS(ctx context.Context, userID, krsID string) (*KrsHeader, error) {
	krs, err := s.repo.GetKrsByID(ctx, krsID)
	if err != nil {
		return nil, errors.New("KRS tidak ditemukan")
	}
	profile, err := s.repo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("profil mahasiswa tidak ditemukan")
	}
	if krs.MahasiswaID != profile.ID {
		return nil, errors.New("akses ditolak: bukan KRS Anda")
	}
	if krs.TotalSKS == 0 {
		return nil, errors.New("KRS kosong, tambahkan minimal 1 kelas")
	}
	if err := s.repo.SubmitKrs(ctx, krsID); err != nil {
		return nil, err
	}

	dpID := ""
	if profile.DosenPembimbingID != nil {
		dpID = *profile.DosenPembimbingID
	}
	s.inv.OnKrsMutation(ctx, profile.ID, dpID)

	return s.getKrsWithDetail(ctx, krsID)
}

// ─── GetHistory ───────────────────────────────────────────────────────────────

func (s *Service) GetHistory(ctx context.Context, userID string) ([]KrsHeader, error) {
	profile, err := s.repo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("profil mahasiswa tidak ditemukan")
	}
	return s.repo.GetHistory(ctx, profile.ID)
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func (s *Service) getKrsWithDetail(ctx context.Context, krsID string) (*KrsHeader, error) {
	krs, err := s.repo.GetKrsByID(ctx, krsID)
	if err != nil {
		return nil, err
	}
	krs.Detail, _ = s.repo.GetDetails(ctx, krsID)
	return krs, nil
}

// checkScheduleConflict — cek apakah jadwal kelas baru overlap dengan kelas di existing.
// Format jam: "08:00:00" atau "08:00"
func checkScheduleConflict(existing []KelasDetail, newKelas *KelasDetail) error {
	for _, ex := range existing {
		for _, exJ := range ex.Jadwal {
			for _, newJ := range newKelas.Jadwal {
				if !strings.EqualFold(exJ.Hari, newJ.Hari) {
					continue
				}
				// Overlap: exStart < newEnd && newStart < exEnd
				if timeStr(exJ.JamMulai) < timeStr(newJ.JamSelesai) &&
					timeStr(newJ.JamMulai) < timeStr(exJ.JamSelesai) {
					return fmt.Errorf("bentrok jadwal hari %s (%s–%s)", newJ.Hari, newJ.JamMulai, newJ.JamSelesai)
				}
			}
		}
	}
	return nil
}

// timeStr — normalisasi "08:00:00" → "08:00" untuk perbandingan string.
func timeStr(t string) string {
	if len(t) >= 5 {
		return t[:5]
	}
	return t
}
