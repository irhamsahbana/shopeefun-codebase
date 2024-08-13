package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func AuthRole(authorizedRoles []string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		forbiddenResponse := fiber.Map{
			"message": "Terlarang: role anda tidak diizinkan untuk mengakses resource ini",
			"success": false,
		}

		role, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(forbiddenResponse)
		}

		for _, authorizedRole := range authorizedRoles {
			if role == authorizedRole {
				return c.Next()
			}
		}

		payload := struct {
			Role           string   `json:"role"`
			AuthorizedRole []string `json:"authorized_roles"`
		}{
			Role:           role,
			AuthorizedRole: authorizedRoles,
		}

		log.Warn().Any("payload", payload).Msg("middleware::AuthRole - Unauthorized")
		return c.Status(fiber.StatusForbidden).JSON(forbiddenResponse)
	}
}
