package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type AppConfig struct {
	ServerAddress string `mapstructure:"server_address"`
	ClientURL     string `mapstructure:"client_url"`
	LogLevel      string `mapstructure:"log_level"`
	SecretKey     string `mapstructure:"secret_key"`

	StorageType string `mapstructure:"storage_type"`
	DatabaseDSN string `mapstructure:"database_dsn"`
}

func New() *AppConfig {
	v := viper.New()

	v.SetDefault("server_address", "localhost:8080")
	v.SetDefault("client_url", "http://localhost:8080")
	v.SetDefault("log_level", "info")
	v.SetDefault("secret_key", "super-secret-key")

	v.SetDefault("storage_type", "memory")
	v.SetDefault("database_dsn", "postgres://user:pass@localhost:5432/gophkeeper?sslmode=disable")

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if cfg := os.Getenv("CONFIG"); cfg != "" {
		v.SetConfigFile(cfg)
		_ = v.ReadInConfig()
	}

	var c AppConfig
	if err := v.Unmarshal(&c); err != nil {
		panic(err)
	}

	return &c
}
