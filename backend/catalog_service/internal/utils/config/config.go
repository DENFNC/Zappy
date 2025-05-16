package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		LogType        string `yaml:"log_type" env-default:"dev"`
		PaginateSecret string `yaml:"paginate_secret" env-required:"true"`
		GRPC           ConfigGRPC
		HTTP           ConfigHTTP
		Postgres       ConfigPSQL
		ObjectStore    BucketCfg `yaml:"object_store"`
	}

	ConfigGRPC struct {
		Port       int           `yaml:"port" env-required:"true"`
		Timeout    time.Duration `yaml:"timeout"`
		Reflection bool          `yaml:"reflection"`
	}

	ConfigHTTP struct {
		Port int `yaml:"port" env-required:"true"`
	}

	ConfigPSQL struct {
		URL string `yaml:"url" env-required:"true"`
	}

	BucketCfg struct {
		ObjectOrigin  string `yaml:"object_origin" env-default:"http://localhost:9000"`
		ImageBucket   string `yaml:"aws_bucket_image" env-required:"true"`
		StagingBucket string `yaml:"aws_staging_bucket" env-required:"true"`
		Buckets       map[string]struct {
			Name      string   `yaml:"name"`
			MimeTypes []string `yaml:"mime_types"`
			Path      string   `yaml:"path"`
		} `yaml:"buckets"`
	}
)

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
