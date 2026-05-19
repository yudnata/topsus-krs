package main

import (
	"log"
	"strings"

	"backend/internal/cache"
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/modules/admin"
	"backend/internal/modules/auth"
	"backend/internal/modules/dosen"
	"backend/internal/modules/mahasiswa"
	"backend/internal/router"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func main() {
	cfg := config.Load()
	db := database.Connect(cfg.DatabaseURL)
	database.Migrate(db)

	redisClient := cache.New(cfg.RedisURL)
	defer redisClient.Close()
	invalidator := cache.NewInvalidator(redisClient)

	// Auth
	authRepo := auth.NewRepository(db)
	authSvc := auth.NewService(authRepo, cfg)
	authHand := auth.NewHandler(authSvc)

	// Admin (master data)
	adminRepo := admin.NewRepository(db)
	adminSvc := admin.NewService(adminRepo, invalidator)
	adminHand := admin.NewHandler(adminSvc)

	// Mahasiswa (KRS flow)
	mhsRepo := mahasiswa.NewRepository(db)
	mhsSvc := mahasiswa.NewService(mhsRepo, invalidator)
	mhsHand := mahasiswa.NewHandler(mhsSvc)

	// Dosen (approval KRS)
	dosenRepo := dosen.NewRepository(db)
	dosenSvc := dosen.NewService(dosenRepo, invalidator)
	dosenHand := dosen.NewHandler(dosenSvc)

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: splitOrigins(cfg.CORSOrigin),
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	router.Setup(app, authHand, authSvc, adminHand, mhsHand, dosenHand, redisClient, invalidator)

	log.Fatal(app.Listen(":" + cfg.Port))
}

func splitOrigins(origins string) []string {
	parts := strings.Split(origins, ",")
	out := make([]string, 0, len(parts))
	for _, o := range parts {
		if trimmed := strings.TrimSpace(o); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	if len(out) == 0 {
		return []string{"http://localhost:5173"}
	}
	return out
}

