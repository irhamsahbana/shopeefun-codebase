package adapter

import (
	// "log"

	"codebase-app/internal/infrastructure/config"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func WithShopeefunPostgres() Option {
	return func(a *Adapter) {
		dbUser := config.Envs.ShopeefunPostgres.Username
		dbPassword := config.Envs.ShopeefunPostgres.Password
		dbName := config.Envs.ShopeefunPostgres.Database
		dbHost := config.Envs.ShopeefunPostgres.Host
		dbSSLMode := config.Envs.ShopeefunPostgres.SslMode
		dbPort := config.Envs.ShopeefunPostgres.Port

		dbMaxPoolSize := config.Envs.DB.MaxOpenCons
		dbMaxIdleConns := config.Envs.DB.MaxIdleCons
		dbConnMaxLifetime := config.Envs.DB.ConnMaxLifetime

		connectionString := "user=" + dbUser + " password=" + dbPassword + " host=" + dbHost + " port=" + dbPort + " dbname=" + dbName + " sslmode=" + dbSSLMode + " TimeZone=UTC"
		db, err := sqlx.Connect("postgres", connectionString)
		if err != nil {
			log.Fatal().Err(err).Msg("Error connecting to Postgres")
		}

		db.SetMaxOpenConns(dbMaxPoolSize)
		db.SetMaxIdleConns(dbMaxIdleConns)
		db.SetConnMaxLifetime(time.Duration(dbConnMaxLifetime) * time.Second)

		// check connection
		err = db.Ping()
		if err != nil {
			log.Fatal().Err(err).Msg("Error connecting to Digihub Postgres")
		}

		a.ShopeefunPostgres = db
		log.Info().Msg("Digihub Postgres connected")
	}
}
