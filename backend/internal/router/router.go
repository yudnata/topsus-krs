package router

import (
	"backend/internal/cache"
	"backend/internal/modules/admin"
	"backend/internal/modules/auth"
	"backend/internal/modules/cachemeta"
	"backend/internal/modules/dosen"
	"backend/internal/modules/mahasiswa"
	"backend/pkg/response"

	"github.com/gofiber/fiber/v3"
)

// Setup mendaftarkan semua modul feature-based di bawah /api.
func Setup(
	app *fiber.App,
	authH *auth.Handler, authSvc *auth.Service,
	adminH *admin.Handler,
	mhsH *mahasiswa.Handler,
	dosenH *dosen.Handler,
	redis *cache.Client,
	inv *cache.Invalidator,
) {
	api := app.Group("/api")

	api.Get("/health", func(c fiber.Ctx) error {
		return healthHandler(c, redis)
	})

	// Slice 1–3
	auth.RegisterRoutes(api, authH, authSvc)
	cachemeta.RegisterRoutes(api, cachemeta.NewHandler(redis, inv))

	// Slice 4
	admin.RegisterRoutes(api, adminH, authSvc)
	mahasiswa.RegisterRoutes(api, mhsH, authSvc)
	dosen.RegisterRoutes(api, dosenH, authSvc)
}

func healthHandler(c fiber.Ctx, redis *cache.Client) error {
	redisStatus := "disabled"
	if redis.Enabled() {
		if err := redis.Ping(c.Context()); err != nil {
			redisStatus = "error"
		} else {
			redisStatus = "ok"
		}
	}
	return response.JSON(c, 200, true, "OK", fiber.Map{
		"status":  "ok",
		"service": "krs-api",
		"redis":   redisStatus,
	})
}
