package config

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Env        string     `mapstructure:"env" validate:"required"`
	HttpServer HttpServer `mapstructure:"http_server" validate:"required"`
	Redis      Redis      `mapstructure:"redis" validate:"required"`
	Postgres   Postgres   `mapstructure:"postgres" validate:"required"`
}

type HttpServer struct {
	Port string `mapstructure:"port" validate:"required"`
}

type Redis struct {
	Addr     string `mapstructure:"addr" validate:"required"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type Postgres struct {
	DSN string `mapstructure:"dsn" validate:"required"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config File does not exist: %s", configPath)
	}

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		log.Fatalf("Config validation failed: %v", err)
	}

	return &cfg
}
