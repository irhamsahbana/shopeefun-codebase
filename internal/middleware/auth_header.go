package middleware

import (
	"codebase-app/pkg/jwthandler"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func AuthBearer(c *fiber.Ctx) error {
	AccessToken := c.Get("Authorization")
	unauthorizedResponse := fiber.Map{
		"message": "Unauthorized",
		"success": false,
	}

	// If the cookie is not set, return an unauthorized status
	if AccessToken == "" {
		log.Error().Msg("middleware::AuthMiddleware - Unauthorized [Header not set]")
		return c.Status(fiber.StatusUnauthorized).JSON(unauthorizedResponse)
	}

	// remove the Bearer prefix
	if len(AccessToken) > 7 {
		AccessToken = AccessToken[7:]
	}

	// Parse the JWT string and store the result in `claims`
	claims, err := jwthandler.ParseTokenString(AccessToken)
	if err != nil {
		log.Error().Err(err).Any("payload", AccessToken).Msg("middleware::AuthMiddleware - Error while parsing token")
		return c.Status(fiber.StatusUnauthorized).JSON(unauthorizedResponse)
	}

	c.Locals("user_id", claims.UserId)
	c.Locals("role", claims.Role)

	// If the token is valid, pass the request to the next handler
	return c.Next()
}
