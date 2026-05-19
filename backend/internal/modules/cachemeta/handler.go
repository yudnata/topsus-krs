package cachemeta

import (
	"time"

	"backend/internal/cache"
	"backend/pkg/response"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	cache *cache.Client
	inv   *cache.Invalidator
}

func NewHandler(c *cache.Client, inv *cache.Invalidator) *Handler {
	return &Handler{cache: c, inv: inv}
}

func (h *Handler) Status(c fiber.Ctx) error {
	status := "disabled"
	if h.cache.Enabled() {
		if err := h.cache.Ping(c.Context()); err != nil {
			status = "error"
		} else {
			status = "ok"
		}
	}
	return response.JSON(c, 200, true, "OK", fiber.Map{
		"redis": status,
	})
}

// DemoGet menguji cache hit/miss — header X-Cache: HIT | MISS
func (h *Handler) DemoGet(c fiber.Ctx) error {
	key := c.Query("key", "cache:demo")
	val, hit, err := h.cache.Get(c.Context(), key)
	if err != nil {
		return response.JSON(c, 500, false, err.Error(), nil)
	}
	if hit {
		c.Set("X-Cache", "HIT")
		return response.JSON(c, 200, true, "Cache hit", fiber.Map{"key": key, "value": val})
	}
	c.Set("X-Cache", "MISS")
	return response.JSON(c, 200, true, "Cache miss", fiber.Map{"key": key})
}

func (h *Handler) DemoSet(c fiber.Ctx) error {
	key := c.Query("key", "cache:demo")
	value := c.Query("value", "hello-krs")
	ttl := time.Duration(cache.TTLDosenMahasiswaList) * time.Second
	if err := h.cache.Set(c.Context(), key, value, ttl); err != nil {
		return response.JSON(c, 500, false, err.Error(), nil)
	}
	return response.JSON(c, 200, true, "Cached", fiber.Map{"key": key, "ttl_seconds": cache.TTLDosenMahasiswaList})
}

func (h *Handler) DemoInvalidate(c fiber.Ctx) error {
	key := c.Query("key", "cache:demo")
	if err := h.cache.Del(c.Context(), key); err != nil {
		return response.JSON(c, 500, false, err.Error(), nil)
	}
	return response.JSON(c, 200, true, "Invalidated", fiber.Map{"key": key})
}
