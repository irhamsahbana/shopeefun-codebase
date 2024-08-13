package adapter

import (
	"fmt"
	"net/http"
	"strings"

	// import "codebase-app/internal/pkg/validator"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var (
	Adapters *Adapter
)

type Option func(adapter *Adapter)

type Validator interface {
	Validate(i any) error
}

type Adapter struct {
	// Driving Adapters
	RestServer *fiber.App
	WsServer   *http.Server

	//Driven Adapters
	ShopeefunPostgres *sqlx.DB
	Validator         Validator // *validator.Validator
	ShopeefunStorage  *s3.Client
}

func (a *Adapter) Sync(opts ...Option) {
	for o := range opts {
		opt := opts[o]
		opt(a)
	}
}

func (a *Adapter) Unsync() error {
	var errs []string

	if a.RestServer != nil {
		if err := a.RestServer.Shutdown(); err != nil {
			errs = append(errs, err.Error())
		}
		log.Info().Msg("Rest server disconnected")
	}

	if a.WsServer != nil {
		if err := a.WsServer.Close(); err != nil {
			errs = append(errs, err.Error())
		}
		log.Info().Msg("Ws server disconnected")
	}

	if a.ShopeefunPostgres != nil {
		if err := a.ShopeefunPostgres.Close(); err != nil {
			errs = append(errs, err.Error())
		}
		log.Info().Msg("Digihub Postgres disconnected")
	}

	if len(errs) > 0 {
		err := fmt.Errorf(strings.Join(errs, "\n"))
		log.Error().Msgf("Error while disconnecting adapters: %v", err)
		return err
	}

	return nil
}
