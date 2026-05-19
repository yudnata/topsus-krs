package cachemeta

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(router fiber.Router, h *Handler) {
	g := router.Group("/cache")
	g.Get("/status", h.Status)
	g.Get("/demo", h.DemoGet)
	g.Post("/demo", h.DemoSet)
	g.Delete("/demo", h.DemoInvalidate)
}
