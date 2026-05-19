package dosen

import (
	"backend/internal/modules/auth"

	"github.com/gofiber/fiber/v3"
)

// RegisterRoutes mendaftarkan endpoint approval di bawah /api/approval.
// DOSEN dan STAFF boleh akses; otorisasi DPA dicek di service layer.
func RegisterRoutes(api fiber.Router, h *Handler, authSvc *auth.Service) {
	approval := api.Group("/approval",
		auth.RequireAuth(authSvc),
		auth.RequireRole("DOSEN", "STAFF", "ADMIN"),
	)

	approval.Get("/pending", h.ListPending)
	approval.Post("/:id/approve", h.Approve)
	approval.Post("/:id/reject", h.Reject)
}
