package main

import (
	"context"
	"fmt"
	"time"

	"familyquest-backend/pkg/db"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	HTTP     HTTPConfig     `env:",prefix=HTTP_"`
	Postgres PostgresConfig `env:",prefix=POSTGRES_"`
	Auth     AuthConfig     `env:",prefix=AUTH_"`
}

type HTTPConfig struct {
	Addr            string        `env:"ADDR,default=:8080"`
	ReadTimeout     time.Duration `env:"READ_TIMEOUT,default=5s"`
	WriteTimeout    time.Duration `env:"WRITE_TIMEOUT,default=10s"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT,default=5s"`
}

type PostgresConfig struct {
	Host            string        `env:"HOST,default=127.0.0.1"`
	Port            int           `env:"PORT,default=5432"`
	User            string        `env:"USER,default=familyquest"`
	Password        string        `env:"PASSWORD,default=familyquest"`
	Database        string        `env:"DATABASE,default=familyquest"`
	SSLMode         string        `env:"SSL_MODE,default=disable"`
	MaxConns        int32         `env:"MAX_CONNS,default=10"`
	MinConns        int32         `env:"MIN_CONNS,default=1"`
	ConnMaxLifetime time.Duration `env:"CONN_MAX_LIFETIME,default=30m"`
	ConnMaxIdleTime time.Duration `env:"CONN_MAX_IDLE_TIME,default=5m"`
}

type AuthConfig struct {
	JWTSecret      string        `env:"JWT_SECRET,default=familyquest-local-secret"`
	JWTIssuer      string        `env:"JWT_ISSUER,default=familyquest"`
	AccessTokenTTL time.Duration `env:"ACCESS_TOKEN_TTL,default=24h"`
}

func LoadConfig(ctx context.Context) (Config, error) {
	var cfg Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func NewPostgresConfig(cfg Config) db.Config {
	return db.Config{
		DSN: fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=%s",
			cfg.Postgres.User,
			cfg.Postgres.Password,
			cfg.Postgres.Host,
			cfg.Postgres.Port,
			cfg.Postgres.Database,
			cfg.Postgres.SSLMode,
		),
		MaxConns:        cfg.Postgres.MaxConns,
		MinConns:        cfg.Postgres.MinConns,
		ConnMaxLifetime: cfg.Postgres.ConnMaxLifetime,
		ConnMaxIdleTime: cfg.Postgres.ConnMaxIdleTime,
	}
}
