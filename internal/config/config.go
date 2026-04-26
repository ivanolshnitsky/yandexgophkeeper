package config

import (
	"flag"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type AppConfig struct {
	ServerAddress string `mapstructure:"server_address"`
	BaseURL       string `mapstructure:"base_url"`
	LogLevel      string `mapstructure:"log_level"`
	SecretKey     string `mapstructure:"secret_key"`
}

// New загружает конфиг с приоритетом:
// defaults < config file < flags < ENV
func New() *AppConfig {
	const (
		defaultServerAddress = "localhost:8080"
		defaultBaseURL       = "http://localhost:8080"
		defaultLogLevel      = "info"
		defaultSecretKey     = "super-secret-key"
	)

	fs := flag.NewFlagSet("gophkeeper", flag.ContinueOnError)

	flagServer := fs.String("a", "", "server address")
	flagBase := fs.String("b", "", "base url")
	flagLog := fs.String("l", "", "log level")

	var configPath string
	fs.StringVar(&configPath, "c", "", "config file path")

	_ = fs.Parse(os.Args[1:])

	if envConfig := os.Getenv("CONFIG"); envConfig != "" {
		configPath = envConfig
	}

	v := viper.New()

	// defaults
	v.SetDefault("server_address", defaultServerAddress)
	v.SetDefault("base_url", defaultBaseURL)
	v.SetDefault("log_level", defaultLogLevel)
	v.SetDefault("secret_key", defaultSecretKey)

	// config file
	if configPath != "" {
		v.SetConfigFile(configPath)
		_ = v.ReadInConfig()
	}

	// ENV
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// flags
	if *flagServer != "" {
		v.Set("server_address", *flagServer)
	}
	if *flagBase != "" {
		v.Set("base_url", *flagBase)
	}
	if *flagLog != "" {
		v.Set("log_level", *flagLog)
	}

	cfg := &AppConfig{}
	if err := v.Unmarshal(cfg); err != nil {
		panic(err)
	}

	return cfg
}
