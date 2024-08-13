package middleware

import (
	"codebase-app/pkg/jwthandler"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func AuthMiddleware(c *fiber.Ctx) error {
	// Get the access_token cookie
	cookie := c.Cookies("access_token")

	// If the cookie is not set, return an unauthorized status
	if cookie == "" {
		log.Error().Msg("middleware::AuthMiddleware - Unauthorized [Cookie not set]")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"success": false,
		})
	}

	// Parse the JWT string and store the result in `claims`
	claims, err := jwthandler.ParseTokenString(cookie)
	if err != nil {
		log.Error().Err(err).Msg("middleware::AuthMiddleware - Error while parsing token")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request",
			"success": false,
		})
	}

	c.Locals("user_id", claims.UserId)

	// If the token is valid, pass the request to the next handler
	return c.Next()
}
