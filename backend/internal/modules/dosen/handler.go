package dosen

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

// GET /api/approval/pending — daftar KRS menunggu persetujuan
func (h *Handler) ListPending(c fiber.Ctx) error {
	userID, _ := c.Locals("userID").(string)
	role, _ := c.Locals("role").(string)
	data, err := h.svc.ListPending(c.Context(), userID, role)
	if err != nil {
		return response.JSON(c, 500, false, err.Error(), nil)
	}
	return response.JSON(c, 200, true, "OK", data)
}

// POST /api/approval/:id/approve — setujui KRS
func (h *Handler) Approve(c fiber.Ctx) error {
	krsID := c.Params("id")
	userID, _ := c.Locals("userID").(string)
	role, _ := c.Locals("role").(string)
	var req ApproveReq
	_ = c.Bind().Body(&req)
	krs, err := h.svc.Approve(c.Context(), userID, role, krsID, req)
	if err != nil {
		return response.JSON(c, 422, false, err.Error(), nil)
	}
	return response.JSON(c, 200, true, "KRS disetujui", krs)
}

// POST /api/approval/:id/reject — tolak KRS dengan catatan
func (h *Handler) Reject(c fiber.Ctx) error {
	krsID := c.Params("id")
	userID, _ := c.Locals("userID").(string)
	role, _ := c.Locals("role").(string)
	var req RejectReq
	if err := c.Bind().Body(&req); err != nil {
		return response.JSON(c, 400, false, "Body tidak valid", nil)
	}
	krs, err := h.svc.Reject(c.Context(), userID, role, krsID, req)
	if err != nil {
		return response.JSON(c, 422, false, err.Error(), nil)
	}
	return response.JSON(c, 200, true, "KRS ditolak", krs)
}
