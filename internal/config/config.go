package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type AppConfig struct {
	ServerAddress string        `mapstructure:"server_address"`
	ClientURL     string        `mapstructure:"client_url"`
	LogLevel      string        `mapstructure:"log_level"`
	JWTSecret     string        `mapstructure:"jwt_secret"`
	CryptoKey     string        `mapstructure:"crypto_key"`
	TokenTTL      time.Duration `mapstructure:"token_ttl"`
	DatabaseDSN   string        `mapstructure:"database_dsn"`
}

func New() (*AppConfig, error) {
	v := viper.New()

	v.SetDefault("server_address", "localhost:8080")
	v.SetDefault("client_url", "http://localhost:8080")
	v.SetDefault("log_level", "info")
	v.SetDefault("token_ttl", "1h")

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var c AppConfig

	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("config unmarshal: %w", err)
	}

	if c.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	if c.CryptoKey == "" {
		return nil, fmt.Errorf("CRYPTO_KEY is required")
	}

	if c.DatabaseDSN == "" {
		return nil, fmt.Errorf("DATABASE_DSN is required")
	}

	return &c, nil
}
