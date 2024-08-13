package adapter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func WithRestServer(app *fiber.App) Option {
	log.Info().Msg("Rest server connected")
	return func(a *Adapter) {
		a.RestServer = app
	}
}
