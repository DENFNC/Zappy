package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct{}

func MustLoad(path string) *Config {
	if path == "" {
		panic("path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("file not found")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("error reading config")
	}

	return &cfg
}
