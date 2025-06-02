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
		LogType        string      `yaml:"log_type" env:"LOG_TYPE" env-default:"dev"`
		PaginateSecret string      `yaml:"paginate_secret" env:"PAGINATE_SECRET" env-required:"true"`
		GRPC           ConfigGRPC  `yaml:"grpc"`
		HTTP           ConfigHTTP  `yaml:"http"`
		Postgres       ConfigPSQL  `yaml:"postgres"`
		Redis          ConfigRedis `yaml:"redis"`
		ObjectStore    BucketCfg   `yaml:"object_store"`
	}

	ConfigGRPC struct {
		URL        string        `yaml:"url" env:"GRPC_URL" env-required:"true"`
		Timeout    time.Duration `yaml:"timeout" env:"GRPC_TIMEOUT"`
		Reflection bool          `yaml:"reflection" env:"GRPC_REFLECTION"`
	}

	ConfigHTTP struct {
		URL string `yaml:"url" env:"HTTP_URL" env-required:"true"`
	}

	ConfigPSQL struct {
		URL             string        `yaml:"url" env:"POSTGRES_URL" env-required:"true"`
		MaxRetries      int           `yaml:"max_retries" env:"POSTGRES_MAX_RETRIES" env-required:"true"`
		RetryIntervalMS time.Duration `yaml:"retry_interval_ms" env:"POSTGRES_RETRY_INTERVAL_MS" env-required:"true"`
	}

	ConfigRedis struct {
		URL             string        `yaml:"url" env:"REDIS_URL" env-required:"true"`
		Password        string        `yaml:"password" env:"REDIS_PASSWORD" env-default:""`
		DB              int           `yaml:"db" env:"REDIS_DB" env-required:"true"`
		MaxRetries      int           `yaml:"max_retries" env:"REDIS_MAX_RETRIES" env-required:"true"`
		RetryIntervalMS time.Duration `yaml:"retry_interval_ms" env:"REDIS_RETRY_INTERVAL_MS" env-required:"true"`
	}

	BucketCfg struct {
		ObjectOrigin    string        `yaml:"object_origin" env:"OBJECT_STORE_OBJECT_ORIGIN"`
		AccessKey       string        `yaml:"access_key" env:"OBJECT_STORE_ACCESS_KEY" env-required:"true"`
		SecretKey       string        `yaml:"secret_key" env:"OBJECT_STORE_SECRET_KEY" env-required:"true"`
		ImageBucket     string        `yaml:"aws_bucket_image" env:"OBJECT_STORE_AWS_BUCKET_IMAGE" env-required:"true"`
		StagingBucket   string        `yaml:"aws_staging_bucket" env:"OBJECT_STORE_AWS_STAGING_BUCKET" env-required:"true"`
		MaxRetries      int           `yaml:"max_retries" env:"OBJECT_STORE_MAX_RETRIES" env-required:"true"`
		RetryIntervalMS time.Duration `yaml:"retry_interval_ms" env:"OBJECT_STORE_RETRY_INTERVAL_MS" env-required:"true"`
		Buckets         map[string]struct {
			Name      string   `yaml:"name" env:"NAME"`
			MimeTypes []string `yaml:"mime_types" env:"MIME_TYPES"`
			Path      string   `yaml:"path" env:"PATH"`
		} `yaml:"buckets"`
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
