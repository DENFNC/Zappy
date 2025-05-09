package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogType  string `yaml:"log_type" env-default:"dev"`
	GRPC     ConfigGRPC
	HTTP     ConfigHTTP
	Postgres ConfigPSQL
}

type ConfigGRPC struct {
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout"`
}

type ConfigHTTP struct {
	Port int `yaml:"port" env-required:"true"`
}

type ConfigPSQL struct {
	URL string `yaml:"url" env-required:"true"`
}

func MustLoad(path string) *Config {
	if path == "" {
		panic("path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("file not found")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		fmt.Println(err)
		panic("error reading config")
	}

	return &cfg
}
