package cmd

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/infrastructure"
	"codebase-app/internal/infrastructure/config"
	"codebase-app/internal/middleware"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Client struct {
	Conn   *websocket.Conn
	UserId string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RunWebsocket(cmd *flag.FlagSet, args []string) {
	var (
		envs       = config.Envs
		flagWsPort = cmd.String("port", "4000", "Websocket port")
		WS_PORT    string
	)

	logLevel, err := zerolog.ParseLevel(envs.App.LogLevel)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	if err := cmd.Parse(args); err != nil {
		log.Fatal().Err(err).Msg("Error while parsing flags")
	}

	if envs.App.WSPort != "" {
		WS_PORT = envs.App.WSPort
	} else {
		WS_PORT = *flagWsPort
	}

	server := &http.Server{
		Addr: ":" + WS_PORT,
	}

	adapter.Adapters.Sync(
		adapter.WithWebsocketServer(server),
	)

	infrastructure.InitializeLogger(envs.App.Environtment, envs.App.LogFileWs, logLevel)

	quit := make(chan os.Signal, 1)

	go func() {
		log.Info().Msgf("Websocket server is running on port %s", WS_PORT)

		http.Handle("/ws", Middleware(http.HandlerFunc(ws), middleware.AuthWs))
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal().Err(err).Msg("Error while running websocket server")
		}
	}()

	shutdownSignals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}
	if runtime.GOOS == "windows" {
		shutdownSignals = []os.Signal{os.Interrupt}
	}

	signal.Notify(quit, shutdownSignals...)
	<-quit
	log.Info().Msg("Websocket server is shutting down ...")

	err = adapter.Adapters.Unsync()
	if err != nil {
		log.Error().Err(err).Msg("Error while unsyncing adapter")
	}

	log.Info().Msg("Websocket server is gracefully stopped")
}

func Middleware(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range middleware {
		h = mw(h)
	}
	return h
}

func ws(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error while upgrading connection")
		return
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error().Err(err).Msg("Error while closing connection")
		}
		log.Info().Msg("Connection closed")
	}()

	// handling close connection
	conn.SetCloseHandler(func(code int, text string) error {
		log.Info().Msgf("Connection closed with code: %d, text: %s", code, text)
		return nil
	})

	// handling pong message, and log the message
	conn.SetPongHandler(func(appData string) error {
		log.Info().Msgf("Received pong: %s", appData)
		return nil
	})

	// handling ping message, and send back pong message
	conn.SetPingHandler(func(appData string) error {
		log.Info().Msgf("Received ping: %s", appData)
		err := conn.WriteControl(websocket.PongMessage, nil, time.Time{})
		if err != nil {
			log.Error().Err(err).Msg("Error while writing pong message")
		}
		return nil
	})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error().Err(err).Msg("Error while reading message")
			}
			break
		}

		log.Info().Msgf("Message received: %s", message)
	}
}
