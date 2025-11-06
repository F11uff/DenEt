package config

import "time"

type Config struct {
	Server ServerConfig
	Database DatabaseConfig
	Logger LoggerConfig
	JWT JWTConfig
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logger   LoggerConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string `env:"PORT" default:"8080"`
	Host string `env:"HOST" default:"localhost"`
}

type DatabaseConfig struct {
	URL            string        `env:"DATABASE_URL" default:"postgres://user:password@localhost:5432/user_rewards?sslmode=disable"`
	MaxOpenConns   int           `env:"DB_MAX_OPEN_CONNS" default:"15"`
	MaxIdleConns   int           `env:"DB_MAX_IDLE_CONNS" default:"10"`
	ConnMaxExpired time.Duration `env:"DB_CONN_MAX_EXPIRED" default:"5m"`
}

type LoggerConfig struct {
	LogLevel  string `env:"LOG_LEVEL" default:"info"`
	LogFormat string `env:"LOG_FORMAT" default:"json"`
}

type JWTConfig struct {
	SecretKey   string        `env:"JWT_SECRET_KEY" default:"your-default-secret-key"`
	ExpireTime  time.Duration `env:"JWT_EXPIRE_TIME" default:"24h"`
}

func LoadConfig(path string) (*Config, error) {

}

