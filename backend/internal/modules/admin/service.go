package admin

import (
	"context"
	"errors"
	"strings"

	"backend/internal/cache"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
	inv  *cache.Invalidator
}

func NewService(repo *Repository, inv *cache.Invalidator) *Service {
	return &Service{repo: repo, inv: inv}
}

// ─── Prodi ───────────────────────────────────────────────────────────────────

func (s *Service) ListProdi(ctx context.Context) ([]ProgramStudi, error) {
	return s.repo.ListProdi(ctx)
}

// ─── Tahun Akademik ──────────────────────────────────────────────────────────

func (s *Service) ListTahunAkademik(ctx context.Context) ([]TahunAkademik, error) {
	return s.repo.ListTahunAkademik(ctx)
}

func (s *Service) GetTahunAkademikAktif(ctx context.Context) (*TahunAkademik, error) {
	return s.repo.GetTahunAkademikAktif(ctx)
}

// ─── Mata Kuliah ─────────────────────────────────────────────────────────────

func (s *Service) ListMataKuliah(ctx context.Context) ([]MataKuliah, error) {
	return s.repo.ListMataKuliah(ctx)
}

func (s *Service) CreateMataKuliah(ctx context.Context, req CreateMataKuliahReq) (*MataKuliah, error) {
	req.KodeMK = strings.TrimSpace(req.KodeMK)
	req.NamaMK = strings.TrimSpace(req.NamaMK)
	if req.KodeMK == "" || req.NamaMK == "" || req.SKS <= 0 {
		return nil, errors.New("kode_mk, nama_mk, dan sks wajib diisi (sks > 0)")
	}
	return s.repo.CreateMataKuliah(ctx, req)
}

// ─── Kelas ───────────────────────────────────────────────────────────────────

func (s *Service) ListKelas(ctx context.Context, taID string) ([]Kelas, error) {
	return s.repo.ListKelas(ctx, taID)
}

func (s *Service) CreateKelas(ctx context.Context, req CreateKelasReq) (*Kelas, error) {
	if req.MataKuliahID == "" || req.TahunAkademikID == "" || req.NamaKelas == "" {
		return nil, errors.New("mata_kuliah_id, tahun_akademik_id, nama_kelas wajib diisi")
	}
	if req.Kapasitas <= 0 {
		req.Kapasitas = 40
	}
	return s.repo.CreateKelas(ctx, req)
}

// ─── Dosen ───────────────────────────────────────────────────────────────────

func (s *Service) ListDosen(ctx context.Context) ([]Dosen, error) {
	return s.repo.ListDosen(ctx)
}

func (s *Service) GetDosen(ctx context.Context, id string) (*Dosen, error) {
	return s.repo.GetDosen(ctx, id)
}

func (s *Service) CreateDosen(ctx context.Context, req CreateDosenReq) (*Dosen, error) {
	req.NIP = strings.TrimSpace(req.NIP)
	req.Nama = strings.TrimSpace(req.Nama)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.NIP == "" || req.Nama == "" || req.Email == "" || req.Password == "" {
		return nil, errors.New("nip, nama, email, password wajib diisi")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	d, err := s.repo.CreateDosenWithUser(ctx, req.NIP, req.Nama, req.Email, string(hash))
	if err != nil {
		return nil, err
	}
	s.inv.OnDosenWrite(ctx, d.ID)
	return d, nil
}

// ─── Mahasiswa ───────────────────────────────────────────────────────────────

func (s *Service) ListMahasiswa(ctx context.Context) ([]Mahasiswa, error) {
	return s.repo.ListMahasiswa(ctx)
}

func (s *Service) GetMahasiswa(ctx context.Context, id string) (*Mahasiswa, error) {
	return s.repo.GetMahasiswa(ctx, id)
}

func (s *Service) CreateMahasiswa(ctx context.Context, req CreateMahasiswaReq) (*Mahasiswa, error) {
	req.NIM = strings.TrimSpace(req.NIM)
	req.Nama = strings.TrimSpace(req.Nama)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.NIM == "" || req.Nama == "" || req.ProdiID == "" || req.Email == "" || req.Password == "" {
		return nil, errors.New("nim, nama, prodi_id, email, password wajib diisi")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	m, err := s.repo.CreateMahasiswaWithUser(ctx, req, string(hash))
	if err != nil {
		return nil, err
	}
	dpID := ""
	if m.DosenPembimbingID != nil {
		dpID = *m.DosenPembimbingID
	}
	s.inv.OnMahasiswaWrite(ctx, m.ID, dpID)
	return m, nil
}

// PatchDPA — ADMIN set dosen pembimbing akademik mahasiswa.
func (s *Service) PatchDPA(ctx context.Context, mahasiswaID string, req PatchDPAReq) error {
	oldDosenID, err := s.repo.GetOldDosenPembimbing(ctx, mahasiswaID)
	if err != nil {
		return err
	}
	_, err = s.repo.PatchDPA(ctx, mahasiswaID, req.DosenPembimbingID)
	if err != nil {
		return err
	}
	s.inv.OnDpaChange(ctx, oldDosenID, req.DosenPembimbingID, mahasiswaID)
	return nil
}
