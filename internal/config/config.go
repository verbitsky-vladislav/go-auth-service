package config

import (
	"auth-microservice/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"sync"
)

type application struct {
	PORT                    string `envconfig:"APP_PORT" required:"true"`
	DOMAIN                  string `envconfig:"APP_DOMAIN" required:"true"`
	URL                     string `envconfig:"APP_URL" required:"true"`
	VERIFICATION_PASSPHRASE string `envconfig:"APP_VERIFICATION_PASSPHRASE" required:"true"`
	FRONTEND_URL            string `envconfig:"APP_FRONTEND_URL" required:"true"`
}

type jwt struct {
	ACCESS_LIFE_TIME  int    `envconfig:"JWT_ACCESS_LIFE_TIME" required:"true"`
	REFRESH_LIFE_TIME int    `envconfig:"JWT_REFRESH_LIFE_TIME" required:"true"`
	SECRET_KEY        string `envconfig:"JWT_SECRET_KEY" required:"true"`
}

type google struct {
	REDIRECT_URL  string `envconfig:"GOOGLE_REDIRECT_URL" required:"true"`
	CLIENT_ID     string `envconfig:"GOOGLE_CLIENT_ID" required:"true"`
	CLIENT_SECRET string `envconfig:"GOOGLE_CLIENT_SECRET" required:"true"`
}

type yandex struct {
	REDIRECT_URL  string `envconfig:"YANDEX_REDIRECT_URL" required:"true"`
	CLIENT_ID     string `envconfig:"YANDEX_CLIENT_ID" required:"true"`
	CLIENT_SECRET string `envconfig:"YANDEX_CLIENT_SECRET" required:"true"`
}

type database struct {
	PORT     string `envconfig:"POSTGRES_PORT_EXTERNAL" required:"true"`
	HOST     string `envconfig:"POSTGRES_HOST" required:"true"`
	USER     string `envconfig:"POSTGRES_USER" required:"true"`
	PASSWORD string `envconfig:"POSTGRES_PASSWORD" required:"true"`
	DATABASE string `envconfig:"POSTGRES_DATABASE" required:"true"`
}

type redis struct {
	PASSWORD string `envconfig:"REDIS_PASSWORD" required:"true"`
	HOST     string `envconfig:"REDIS_HOST" required:"true"`
	PORT     string `envconfig:"REDIS_EXTERNAL_PORT" required:"true"`
}

type mailer struct {
	USERNAME string `envconfig:"MAILER_USERNAME" required:"true"`
	PASSWORD string `envconfig:"MAILER_PASSWORD" required:"true"`
}

type Config struct {
	Application application
	Jwt         jwt
	Google      google
	Yandex      yandex
	Database    database
	Redis       redis
	Mailer      mailer
}

var (
	once sync.Once
	cfg  Config
)

func New() *Config {
	once.Do(func() {
		path := ".env"
		err := godotenv.Load(path)
		if err = envconfig.Process("", &cfg); err != nil {
			logger.Panic(err, "Error loading env vars")
		}
	})
	return &cfg
}
