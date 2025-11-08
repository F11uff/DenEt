package config

import (
	"errors"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string `env:"PORT" default:"8080"`
	Host string `env:"HOST" default:"localhost"`
}

type DatabaseConfig struct {
	URL            string        `env:"DATABASE_URL" default:"postgres://postgres:postgresql@localhost:5432/postgres?sslmode=disable"`
	Port           int           `env:"PORT_DB" default:"5432"`
	MaxOpenConns   int           `env:"DB_MAX_OPEN_CONNS" default:"15"`
	MaxIdleConns   int           `env:"DB_MAX_IDLE_CONNS" default:"10"`
	ConnMaxExpired time.Duration `env:"DB_CONN_MAX_EXPIRED" default:"5m"`
}

type JWTConfig struct {
	SecretKey  string        `env:"JWT_SECRET_KEY" default:"your-super-secret-jwt-key"`
	ExpireTime time.Duration `env:"JWT_EXPIRE_TIME" default:"24h"`
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	var conf Config
	if err := env.Parse(&conf); err != nil {
		return nil, errors.New("don't parse .env")
	}

	return &conf, nil
}
