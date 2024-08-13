package config

import (
	"codebase-app/pkg/config"
	"sync"

	"github.com/rs/zerolog/log"
)

var (
	Envs *Config // Envs is global vars Config.
	once sync.Once
)

type Config struct {
	App struct {
		Name                    string `env:"APP_NAME"`
		Environtment            string `env:"APP_ENV" env-default:"production"`
		BaseURL                 string `env:"APP_BASE_URL" env-default:"http://localhost:3000"`
		Port                    string `env:"APP_PORT"`
		WSPort                  string `env:"WS_PORT"`
		LogLevel                string `env:"APP_LOG_LEVEL" env-default:"debug"`
		LogFile                 string `env:"APP_LOG_FILE" env-default:"./logs/app.log"`
		LogFileWs               string `env:"APP_LOG_FILE_WS" env-default:"./logs/ws.log"`
		LocalStoragePublicPath  string `env:"LOCAL_STORAGE_PUBLIC_PATH" env-default:"./storage/public"`
		LocalStoragePrivatePath string `env:"LOCAL_STORAGE_PRIVATE_PATH" env-default:"./storage/private"`
	}
	DB struct {
		ConnectionTimeout int `env:"DB_CONN_TIMEOUT" env-default:"30" env-description:"database timeout in seconds"`
		MaxOpenCons       int `env:"DB_MAX_OPEN_CONS" env-default:"20" env-description:"database max open conn in seconds"`
		MaxIdleCons       int `env:"DB_MAX_IdLE_CONS" env-default:"20" env-description:"database max idle conn in seconds"`
		ConnMaxLifetime   int `env:"DB_CONN_MAX_LIFETIME" env-default:"0" env-description:"database conn max lifetime in seconds"`
	}
	Guard struct {
		JwtPrivateKey   string `env:"JWT_PRIVATE_KEY"`
		JwtPrivateKeyWs string `env:"JWT_PRIVATE_KEY_WS"`
		JwtWsExp        int    `env:"JWT_WS_EXP" env-default:"10"` // 10 seconds
	}
	ShopeefunPostgres struct {
		Host     string `env:"SHOPEEFUN_POSTGRES_HOST" env-default:"localhost"`
		Port     string `env:"SHOPEEFUN_POSTGRES_PORT" env-default:"5432"`
		Username string `env:"SHOPEEFUN_POSTGRES_USER" env-default:"postgres"`
		Password string `env:"SHOPEEFUN_POSTGRES_PASSWORD" env-default:"postgres"`
		Database string `env:"SHOPEEFUN_POSTGRES_DB" env-default:"venatronics"`
		SslMode  string `env:"SHOPEEFUN_POSTGRES_SSL_MODE" env-default:"disable"`
	}
	ShopeefunStorage struct {
		Key      string `env:"SHOPEEFUN_STORAGE_KEY"`
		Secret   string `env:"SHOPEEFUN_STORAGE_SECRET"`
		Endpoint string `env:"SHOPEEFUN_STORAGE_ENDPOINT"`
		Region   string `env:"SHOPEEFUN_STORAGE_REGION"`
		Bucket   string `env:"SHOPEEFUN_STORAGE_BUCKET"`
	}
	Oauth struct {
		Google struct {
			ClientId     string `env:"GOOGLE_CLIENT_ID"`
			ClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
			RedirectURL  string `env:"GOOGLE_REDIRECT_URL"`
		}
	}
}

// Option is Configure type return func.
type Option = func(c *Configure) error

// Configure is the data struct.
type Configure struct {
	path     string
	filename string
}

// Configuration create instance.
func Configuration(opts ...Option) *Configure {
	c := &Configure{}

	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			panic(err)
		}
	}
	return c
}

// Initialize will create instance of Configure.
func (c *Configure) Initialize() {
	once.Do(func() {
		Envs = &Config{}
		if err := config.Load(config.Opts{
			Config:    Envs,
			Paths:     []string{c.path},
			Filenames: []string{c.filename},
		}); err != nil {
			log.Fatal().Err(err).Msg("get config error")
		}
	})
}

// WithPath will assign to field path Configure.
func WithPath(path string) Option {
	return func(c *Configure) error {
		c.path = path
		return nil
	}
}

// WithFilename will assign to field name Configure.
func WithFilename(name string) Option {
	return func(c *Configure) error {
		c.filename = name
		return nil
	}
}
