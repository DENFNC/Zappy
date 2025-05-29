package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AppLogType     string `yaml:"app_log" env-default:"dev"`
	PaginateSecret string `yaml:"paginate_secret"`
	GRPC           gRPCConfig
	HTTP           httpConfig
	Storage        StorageConfig
}

type gRPCConfig struct {
	Port       int           `yaml:"port" env-default:"50052"`
	Timeout    time.Duration `yaml:"timeout" env-default:"10s"`
	Reflection bool          `yaml:"reflection" env-default:"true"`
}

type httpConfig struct {
	Port int `yaml:"port" env-default:"8081"`
}

type StorageConfig struct {
	URL string `yaml:"url" env-required:"true"`
}

func MustLoad(path string) *Config {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}
