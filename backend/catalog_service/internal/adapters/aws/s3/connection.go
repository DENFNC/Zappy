package s3client

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	API            *s3.Client
	PresignClient  *s3.PresignClient
	presignExpires time.Duration
}

type Option func(*clientOptions)

type clientOptions struct {
	region        string
	endpoint      string
	credsProvider aws.CredentialsProvider
	presignExpiry time.Duration
}

func WithRegion(region string) Option {
	return func(o *clientOptions) {
		o.region = region
	}
}

func WithEndpoint(endpoint string) Option {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

func WithCredentials(provider aws.CredentialsProvider) Option {
	return func(o *clientOptions) {
		o.credsProvider = provider
	}
}

func WithPresignExpiry(d time.Duration) Option {
	return func(o *clientOptions) {
		o.presignExpiry = d
	}
}

func NewClient(ctx context.Context, opts ...Option) (*Client, error) {
	co := &clientOptions{
		region:        "us-east-1",
		credsProvider: nil,
		presignExpiry: 15 * time.Minute,
	}
	for _, opt := range opts {
		opt(co)
	}

	loaderOpts := []func(*config.LoadOptions) error{
		config.WithRegion(co.region),
	}
	if co.credsProvider != nil {
		loaderOpts = append(loaderOpts, config.WithCredentialsProvider(co.credsProvider))
	}

	cfg, err := config.LoadDefaultConfig(ctx, loaderOpts...)
	if err != nil {
		return nil, err
	}

	apiOpts := []func(*s3.Options){}
	if co.endpoint != "" {
		apiOpts = append(apiOpts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(co.endpoint)
			o.UsePathStyle = true
		})
	}

	api := s3.NewFromConfig(cfg, apiOpts...)
	presigner := s3.NewPresignClient(api, func(po *s3.PresignOptions) {
		po.Expires = co.presignExpiry
	})

	return &Client{
		API:            api,
		PresignClient:  presigner,
		presignExpires: co.presignExpiry,
	}, nil
}
