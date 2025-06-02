package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	AppLogType     string     `yaml:"app_log" env:"APP_LOG" env-default:"dev"`
	PaginateSecret string     `yaml:"paginate_secret" env:"PAGINATE_SECRET"`
	GRPC           gRPCConfig `yaml:"grpc"`
	HTTP           httpConfig `yaml:"http"`
	Storage        ConfigPSQL `yaml:"storage"`
}

type gRPCConfig struct {
	URL        string        `yaml:"url" env:"GRPC_URL" env-default:"50052"`
	Timeout    time.Duration `yaml:"timeout" env:"GRPC_TIMEOUT" env-default:"10s"`
	Reflection bool          `yaml:"reflection" env:"GRPC_REFLECTION" env-default:"true"`
}

type httpConfig struct {
	URL string `yaml:"url" env:"HTTP_URL" env-default:"8081"`
}

type ConfigPSQL struct {
	URL             string        `yaml:"url" env:"POSTGRES_URL" env-required:"true"`
	MaxRetries      int           `yaml:"max_retries" env:"POSTGRES_MAX_RETRIES" env-default:"3"`
	RetryIntervalMS time.Duration `yaml:"retry_interval_ms" env:"POSTGRES_RETRY_INTERVAL_MS" env-default:"2000ms"`
}

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
