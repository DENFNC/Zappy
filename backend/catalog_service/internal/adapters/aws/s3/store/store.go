package store

import (
	"context"

	s3client "github.com/DENFNC/Zappy/catalog_service/internal/adapters/aws/s3"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Store struct {
	client *s3client.Client
}

func NewStore(
	client *s3client.Client,
) *Store {
	return &Store{
		client: client,
	}
}

func (store *Store) PresignGet(
	ctx context.Context,
	bucket, key string,
) (string, error) {
	presignedURL, err := store.client.PresignClient.PresignGetObject(
		ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		},
	)
	if err != nil {
		return "", err
	}

	return presignedURL.URL, nil
}

func (store *Store) PresignPut(
	ctx context.Context,
	bucket, key, contentType string,
) (string, error) {
	presignedURL, err := store.client.PresignClient.PresignPutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket:      aws.String(bucket),
			Key:         aws.String(key),
			ContentType: aws.String(contentType),
		},
	)
	if err != nil {
		return "", err
	}

	return presignedURL.URL, nil
}
