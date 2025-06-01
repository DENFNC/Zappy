package s3client

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
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

func (c *Client) EnsureBucketExists(ctx context.Context, bucketNames ...string) error {
	for _, bucketName := range bucketNames {
		_, err := c.API.HeadBucket(ctx, &s3.HeadBucketInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			if isNotFound(err) {
				// Создаём бакет
				_, err := c.API.CreateBucket(ctx, &s3.CreateBucketInput{
					Bucket: aws.String(bucketName),
				})
				if err != nil {
					return fmt.Errorf("failed to create bucket %q: %w", bucketName, err)
				}
				// Ждём, пока бакет станет доступен
				waiter := s3.NewBucketExistsWaiter(c.API)
				waitCtx, cancel := context.WithTimeout(ctx, 1*time.Minute)
				defer cancel()

				if waitErr := waiter.Wait(waitCtx, &s3.HeadBucketInput{
					Bucket: aws.String(bucketName),
				}, 5*time.Second); waitErr != nil {
					return fmt.Errorf("bucket %q создан, но не стал доступен: %w", bucketName, waitErr)
				}
				continue
			}
			if isForbidden(err) {
				return fmt.Errorf("no permissions to access bucket %q: %w", bucketName, err)
			}
			return fmt.Errorf("failed to check bucket %q: %w", bucketName, err)
		}
	}
	return nil
}

func isNotFound(err error) bool {
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		code := apiErr.ErrorCode()
		return code == "NoSuchBucket" || code == "NotFound"
	}
	return false
}

func isForbidden(err error) bool {
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		return apiErr.ErrorCode() == "Forbidden"
	}
	return false
}
