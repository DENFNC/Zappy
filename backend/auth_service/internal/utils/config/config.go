package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	AppLog         string `yaml:"app_log" env-default:"dev"`
	PaginateSecret string `yaml:"paginate_secret" env-required:"true"`
	GRPC           ConfigGRPC
	HTTP           ConfigHTTP
	Postgres       ConfigPSQL
	Vault          ConfigVault
}

type ConfigGRPC struct {
	Port       int           `yaml:"port" env-required:"true"`
	Timeout    time.Duration `yaml:"timeout" env-default:"5s"`
	Reflection bool          `yaml:"reflection" env-default:"true"`
}

type ConfigHTTP struct {
	Port int `yaml:"port" env-required:"true"`
}

type ConfigPSQL struct {
	URL string `yaml:"url" env-required:"true"`
}

type ConfigVault struct {
	URL     string        `yaml:"url" env-required:"true"`
	Token   string        `yaml:"token" env-required:"true"`
	AppUUID string        `yaml:"app_uuid" env-required:"true"`
	Issuer  string        `yaml:"issuer" env-required:"true"`
	Expires time.Duration `yaml:"expires" env-default:"5m"`
	KeyName string        `yaml:"key_name" env-required:"true"`
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
