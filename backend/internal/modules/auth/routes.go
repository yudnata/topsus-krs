package auth

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(router fiber.Router, h *Handler, svc *Service) {
	auth := router.Group("/auth")
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
	auth.Get("/profile", RequireAuth(svc), h.Profile)
}
