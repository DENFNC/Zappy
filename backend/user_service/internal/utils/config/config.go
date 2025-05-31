package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		LogType  string     `yaml:"log_type" env:"LOG_TYPE" env-default:"dev"`
		GRPC     ConfigGRPC `yaml:"grpc"`
		Postgres ConfigPSQL `yaml:"postgres"`
	}

	ConfigGRPC struct {
		Port       int           `yaml:"port" env:"GRPC_PORT" env-required:"true"`
		Timeout    time.Duration `yaml:"timeout" env:"GRPC_TIMEOUT"`
		Reflection bool          `yaml:"reflection" env:"GRPC_REFLECTION"`
	}

	ConfigPSQL struct {
		URL string `yaml:"url" env:"POSTGRES_URL" env-required:"true"`
	}
)

func MustLoad(path string) *Config {
	var cfg Config

	if path != "" {
		if _, err := os.Stat(path); err == nil {
			if err := cleanenv.ReadConfig(path, &cfg); err != nil {
				log.Fatalf("failed to read config file: %v", err)
			}
			log.Printf("configuration loaded from file: %s", path)
		} else {
			log.Printf("config file not found at %s, skipping file load", path)
		}
	}

	if err := godotenv.Load(); err == nil {
		log.Println("loaded .env file")
	} else if !os.IsNotExist(err) {
		log.Fatalf("error loading .env file: %v", err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("failed to read configuration from environment: %v", err)
	}

	log.Println("configuration loaded with environment variables precedence")
	return &cfg
}
