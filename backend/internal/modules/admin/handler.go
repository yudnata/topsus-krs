package admin

import (
	"backend/pkg/response"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// ─── Program Studi ───────────────────────────────────────────────────────────

func (h *Handler) ListProdi(c fiber.Ctx) error {
	data, err := h.svc.ListProdi(c.Context())
	if err != nil {
		return response.JSON(c, 500, false, err.Error(), nil)
	}
	return response.JSON(c, 200, true, "OK", data)
}

// ─── Tahun Akademik ──────────────────────────────────────────────────────────

func (h *Handler) ListTahunAkademik(c fiber.Ctx) error {
	data, err := h.svc.ListTahunAkademik(c.Context())
	if err != nil {
		return response.JSON(c, 500, false, err.Error(), nil)
	}
	return response.JSON(c, 200, true, "OK", data)
}

func (h *Handler) GetTahunAkademikAktif(c fiber.Ctx) error {
	data, err := h.svc.GetTahunAkademikAktif(c.Context())
	if err != nil {
		return response.JSON(c, 404, false, "Tidak ada tahun akademik aktif", nil)
	}
	return response.JSON(c, 200, true, "OK", data)
}

// ─── Mata Kuliah ─────────────────────────────────────────────────────────────

func (h *Handler) ListMataKuliah(c fiber.Ctx) error {
	data, err := h.svc.ListMataKuliah(c.Context())
	if err != nil {
		return response.JSON(c, 500, false, err.Error(), nil)
	}
	return response.JSON(c, 200, true, "OK", data)
}

func (h *Handler) CreateMataKuliah(c fiber.Ctx) error {
	var req CreateMataKuliahReq
	if err := c.Bind().Body(&req); err != nil {
		return response.JSON(c, 400, false, "Body tidak valid", nil)
	}
	data, err := h.svc.CreateMataKuliah(c.Context(), req)
	if err != nil {
		return response.JSON(c, 422, false, err.Error(), nil)
	}
	return response.JSON(c, 201, true, "Mata kuliah dibuat", data)
}

// ─── Kelas ───────────────────────────────────────────────────────────────────

func (h *Handler) ListKelas(c fiber.Ctx) error {
	taID := c.Query("tahun_akademik_id", "")
	data, err := h.svc.ListKelas(c.Context(), taID)
	if err != nil {
		return response.JSON(c, 500, false, err.Error(), nil)
	}
	return response.JSON(c, 200, true, "OK", data)
}

func (h *Handler) CreateKelas(c fiber.Ctx) error {
	var req CreateKelasReq
	if err := c.Bind().Body(&req); err != nil {
		return response.JSON(c, 400, false, "Body tidak valid", nil)
	}
	data, err := h.svc.CreateKelas(c.Context(), req)
	if err != nil {
		return response.JSON(c, 422, false, err.Error(), nil)
	}
	return response.JSON(c, 201, true, "Kelas dibuat", data)
}

// ─── Dosen ───────────────────────────────────────────────────────────────────

func (h *Handler) ListDosen(c fiber.Ctx) error {
	data, err := h.svc.ListDosen(c.Context())
	if err != nil {
		return response.JSON(c, 500, false, err.Error(), nil)
	}
	return response.JSON(c, 200, true, "OK", data)
}

func (h *Handler) GetDosen(c fiber.Ctx) error {
	id := c.Params("id")
	data, err := h.svc.GetDosen(c.Context(), id)
	if err != nil {
		return response.JSON(c, 404, false, "Dosen tidak ditemukan", nil)
	}
	return response.JSON(c, 200, true, "OK", data)
}

func (h *Handler) CreateDosen(c fiber.Ctx) error {
	var req CreateDosenReq
	if err := c.Bind().Body(&req); err != nil {
		return response.JSON(c, 400, false, "Body tidak valid", nil)
	}
	data, err := h.svc.CreateDosen(c.Context(), req)
	if err != nil {
		return response.JSON(c, 409, false, err.Error(), nil)
	}
	return response.JSON(c, 201, true, "Dosen dibuat", data)
}

// ─── Mahasiswa ───────────────────────────────────────────────────────────────

func (h *Handler) ListMahasiswa(c fiber.Ctx) error {
	data, err := h.svc.ListMahasiswa(c.Context())
	if err != nil {
		return response.JSON(c, 500, false, err.Error(), nil)
	}
	return response.JSON(c, 200, true, "OK", data)
}

func (h *Handler) GetMahasiswa(c fiber.Ctx) error {
	id := c.Params("id")
	data, err := h.svc.GetMahasiswa(c.Context(), id)
	if err != nil {
		return response.JSON(c, 404, false, "Mahasiswa tidak ditemukan", nil)
	}
	return response.JSON(c, 200, true, "OK", data)
}

func (h *Handler) CreateMahasiswa(c fiber.Ctx) error {
	var req CreateMahasiswaReq
	if err := c.Bind().Body(&req); err != nil {
		return response.JSON(c, 400, false, "Body tidak valid", nil)
	}
	data, err := h.svc.CreateMahasiswa(c.Context(), req)
	if err != nil {
		return response.JSON(c, 409, false, err.Error(), nil)
	}
	return response.JSON(c, 201, true, "Mahasiswa dibuat", data)
}

func (h *Handler) PatchDPA(c fiber.Ctx) error {
	mahasiswaID := c.Params("id")
	var req PatchDPAReq
	if err := c.Bind().Body(&req); err != nil {
		return response.JSON(c, 400, false, "Body tidak valid", nil)
	}
	if req.DosenPembimbingID == "" {
		return response.JSON(c, 400, false, "dosen_pembimbing_id wajib diisi", nil)
	}
	if err := h.svc.PatchDPA(c.Context(), mahasiswaID, req); err != nil {
		return response.JSON(c, 422, false, err.Error(), nil)
	}
	return response.JSON(c, 200, true, "DPA berhasil diperbarui", nil)
}
