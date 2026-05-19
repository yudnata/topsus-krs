package mahasiswa

import (
	"backend/internal/modules/auth"

	"github.com/gofiber/fiber/v3"
)

// RegisterRoutes mendaftarkan endpoint KRS mahasiswa di bawah /api/krs.
// Semua dilindungi JWT + role MAHASISWA.
func RegisterRoutes(api fiber.Router, h *Handler, authSvc *auth.Service) {
	krs := api.Group("/krs",
		auth.RequireAuth(authSvc),
		auth.RequireRole("MAHASISWA"),
	)

	krs.Get("/current", h.GetCurrentKRS)
	krs.Get("/history", h.GetHistory)
	krs.Post("/:id/add-class", h.AddClass)
	krs.Delete("/:id/remove-class", h.RemoveClass)
	krs.Post("/:id/submit", h.SubmitKRS)
}
