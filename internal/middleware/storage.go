package middleware

import (
	"codebase-app/internal/infrastructure/config"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ValidateSignedURL(c *fiber.Ctx) error {
	// Parse expiration time and signature from query parameters
	var (
		expiresStr     = c.Query("expires")
		signature      = c.Query("signature")
		ErrUrlNotValid = fiber.Map{
			"success": false,
			"message": "URL not valid",
		}
	)

	// Convert expiration time to int64
	expires, err := strconv.ParseInt(expiresStr, 10, 64)
	if err != nil || time.Now().Unix() > expires {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrUrlNotValid)
	}

	// Recreate the original data and signature
	data := fmt.Sprintf("%s%d", c.BaseURL()+c.Path(), expires)
	h := hmac.New(sha256.New, []byte(config.Envs.Guard.JwtPrivateKey))
	h.Write([]byte(data))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	// Compare the signatures
	if !hmac.Equal([]byte(expectedSignature), []byte(signature)) {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrUrlNotValid)
	}

	return c.Next()
}
