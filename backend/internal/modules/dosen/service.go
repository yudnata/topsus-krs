package dosen

import (
	"context"
	"errors"
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

// ─── List Pending ─────────────────────────────────────────────────────────────

// ListPending — DOSEN hanya melihat mahasiswa bimbingannya; STAFF melihat semua.
func (s *Service) ListPending(ctx context.Context, userID, role string) ([]PendingKRS, error) {
	role = strings.ToUpper(role)
	if role == "STAFF" || role == "ADMIN" {
		return s.repo.ListPendingAll(ctx)
	}
	// DOSEN — ambil profil dosen dulu
	profile, err := s.repo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("profil dosen tidak ditemukan")
	}
	return s.repo.ListPendingForDosen(ctx, profile.ID)
}

// ─── Approve ─────────────────────────────────────────────────────────────────

func (s *Service) Approve(ctx context.Context, userID, role, krsID string, req ApproveReq) (*PendingKRS, error) {
	// Pastikan DOSEN hanya bisa ACC mahasiswa bimbingannya
	if err := s.checkDosenAuthorization(ctx, userID, role, krsID); err != nil {
		return nil, err
	}

	reviewerID, err := s.getUserIDFromContext(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Approve(ctx, krsID, reviewerID, req.Catatan); err != nil {
		return nil, err
	}

	// Invalidate cache
	_, mahasiswaID, _ := s.repo.GetDosenPembimbingOfKrs(ctx, krsID)
	s.inv.OnKrsMutation(ctx, mahasiswaID, "")

	return s.repo.GetKrs(ctx, krsID)
}

// ─── Reject ──────────────────────────────────────────────────────────────────

func (s *Service) Reject(ctx context.Context, userID, role, krsID string, req RejectReq) (*PendingKRS, error) {
	if strings.TrimSpace(req.Catatan) == "" {
		return nil, errors.New("catatan penolakan wajib diisi")
	}

	if err := s.checkDosenAuthorization(ctx, userID, role, krsID); err != nil {
		return nil, err
	}

	reviewerID, err := s.getUserIDFromContext(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Reject(ctx, krsID, reviewerID, req.Catatan); err != nil {
		return nil, err
	}

	_, mahasiswaID, _ := s.repo.GetDosenPembimbingOfKrs(ctx, krsID)
	s.inv.OnKrsMutation(ctx, mahasiswaID, "")

	return s.repo.GetKrs(ctx, krsID)
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

// checkDosenAuthorization — DOSEN hanya boleh approve mahasiswa bimbingannya (DPA).
// STAFF/ADMIN boleh semua.
func (s *Service) checkDosenAuthorization(ctx context.Context, userID, role, krsID string) error {
	role = strings.ToUpper(role)
	if role == "STAFF" || role == "ADMIN" {
		return nil
	}
	// Cek apakah DOSEN ini adalah pembimbing mahasiswa pemilik KRS
	dosenPembimbingID, _, err := s.repo.GetDosenPembimbingOfKrs(ctx, krsID)
	if err != nil {
		return errors.New("KRS tidak ditemukan")
	}
	profile, err := s.repo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return errors.New("profil dosen tidak ditemukan")
	}
	if dosenPembimbingID != profile.ID {
		return errors.New("otorisasi ditolak: mahasiswa ini bukan bimbingan Anda")
	}
	return nil
}

// getUserIDFromContext — ambil user ID dari c.Locals("userID") yang sudah di-set middleware.
// Di sini kita pass langsung userID (sudah ada di handler).
func (s *Service) getUserIDFromContext(ctx context.Context, userID string) (string, error) {
	if userID == "" {
		return "", errors.New("user tidak teridentifikasi")
	}
	return userID, nil
}
