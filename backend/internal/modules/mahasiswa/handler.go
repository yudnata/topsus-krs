package mahasiswa

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

// GET /api/krs/current — KRS aktif mahasiswa (draft atau terkini)
func (h *Handler) GetCurrentKRS(c fiber.Ctx) error {
	userID, _ := c.Locals("userID").(string)
	krs, err := h.svc.GetCurrentKRS(c.Context(), userID)
	if err != nil {
		return response.JSON(c, 404, false, err.Error(), nil)
	}
	return response.JSON(c, 200, true, "OK", krs)
}

// POST /api/krs/:id/add-class — tambah kelas ke draft KRS
func (h *Handler) AddClass(c fiber.Ctx) error {
	krsID := c.Params("id")
	userID, _ := c.Locals("userID").(string)
	var req AddClassReq
	if err := c.Bind().Body(&req); err != nil || req.KelasID == "" {
		return response.JSON(c, 400, false, "kelas_id wajib diisi", nil)
	}
	krs, err := h.svc.AddClass(c.Context(), userID, krsID, req.KelasID)
	if err != nil {
		return response.JSON(c, 422, false, err.Error(), nil)
	}
	return response.JSON(c, 200, true, "Kelas ditambahkan", krs)
}

// DELETE /api/krs/:id/remove-class — hapus kelas dari draft KRS
func (h *Handler) RemoveClass(c fiber.Ctx) error {
	krsID := c.Params("id")
	userID, _ := c.Locals("userID").(string)
	var req RemoveClassReq
	if err := c.Bind().Body(&req); err != nil || req.KelasID == "" {
		return response.JSON(c, 400, false, "kelas_id wajib diisi", nil)
	}
	krs, err := h.svc.RemoveClass(c.Context(), userID, krsID, req.KelasID)
	if err != nil {
		return response.JSON(c, 422, false, err.Error(), nil)
	}
	return response.JSON(c, 200, true, "Kelas dihapus", krs)
}

// POST /api/krs/:id/submit — ajukan KRS (draft → diajukan)
func (h *Handler) SubmitKRS(c fiber.Ctx) error {
	krsID := c.Params("id")
	userID, _ := c.Locals("userID").(string)
	krs, err := h.svc.SubmitKRS(c.Context(), userID, krsID)
	if err != nil {
		return response.JSON(c, 422, false, err.Error(), nil)
	}
	return response.JSON(c, 200, true, "KRS berhasil diajukan", krs)
}

// GET /api/krs/history — riwayat semua KRS mahasiswa
func (h *Handler) GetHistory(c fiber.Ctx) error {
	userID, _ := c.Locals("userID").(string)
	history, err := h.svc.GetHistory(c.Context(), userID)
	if err != nil {
		return response.JSON(c, 500, false, err.Error(), nil)
	}
	return response.JSON(c, 200, true, "OK", history)
}
