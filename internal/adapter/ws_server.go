package adapter

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

// WithVenamonGolog sets up a new Venamon Golog bot using the provided token and settings.
//
// It assigns the newly created bot to the Adapter's VenamonGolog field.
// When using this option for hooks or commands, make sure to check if the bot is not nil.
func WithWebsocketServer(s *http.Server) Option {
	log.Info().Msg("Websocket server is running")
	return func(a *Adapter) {
		a.WsServer = s
	}
}
