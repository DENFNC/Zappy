package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogType  string `yaml:"log_type" env-default:"dev"`
	GRPC     ConfigGRPC
	Postgres ConfigPSQL
	// Vault    ConfigVault
}

type ConfigGRPC struct {
	Port int `yaml:"port" env-required:"true"`
}

type ConfigPSQL struct {
	URL string `yaml:"url" env-required:"true"`
}

// type ConfigVault struct {
// 	URL     string        `yaml:"url" env-required:"true"`
// 	Token   string        `yaml:"token" env-required:"true"`
// 	AppUUID string        `yaml:"app_uuid" env-required:"true"`
// 	Issuer  string        `yaml:"issuer" env-required:"true"`
// 	Expires time.Duration `yaml:"expires" env-default:"5m"`
// 	KeyName string        `yaml:"key_name" env-required:"true"`
// }

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
