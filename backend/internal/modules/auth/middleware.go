package auth

import (
	"strings"

	"backend/pkg/response"

	"github.com/gofiber/fiber/v3"
)

// RequireAuth — middleware JWT per modul auth (feature-based).
func RequireAuth(svc *Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.JSON(c, 401, false, "Missing authorization header", nil)
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return response.JSON(c, 401, false, "Invalid authorization format", nil)
		}

		userID, err := svc.ValidateToken(parts[1])
		if err != nil {
			return response.JSON(c, 401, false, "Invalid or expired token", nil)
		}

		user, err := svc.repo.FindByID(c.Context(), userID)
		if err != nil {
			return response.JSON(c, 401, false, "User not found", nil)
		}

		c.Locals("user", user)
		c.Locals("userID", userID)
		c.Locals("role", user.Role)
		return c.Next()
	}
}

// RequireRole membatasi route ke role tertentu.
func RequireRole(roles ...string) fiber.Handler {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[strings.ToUpper(r)] = true
	}
	return func(c fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		if !allowed[strings.ToUpper(role)] {
			return response.JSON(c, 403, false, "Forbidden", nil)
		}
		return c.Next()
	}
}
