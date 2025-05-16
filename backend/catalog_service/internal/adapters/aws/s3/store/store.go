package store

import (
	"context"
	"fmt"
	"io"

	s3client "github.com/DENFNC/Zappy/catalog_service/internal/adapters/aws/s3"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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

func (store *Store) GetObject(
	ctx context.Context,
	bucket, key string,
) (io.ReadCloser, error) {
	object, err := store.client.API.GetObject(ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		},
	)
	if err != nil {
		return nil, err
	}

	return object.Body, nil
}

func (store *Store) GetObjectRange(
	ctx context.Context,
	bucket, key string,
	byteRange string,
) (io.ReadCloser, error) {
	object, err := store.client.API.GetObject(ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Range:  aws.String(byteRange),
		},
	)
	if err != nil {
		return nil, err
	}

	return object.Body, nil
}

func (store *Store) CopyObject(
	ctx context.Context,
	dstBucket, dstKey, src string,
) error {
	_, err := store.client.API.CopyObject(
		ctx,
		&s3.CopyObjectInput{
			Bucket:            aws.String(dstBucket),
			Key:               aws.String(dstKey),
			CopySource:        aws.String(src),
			MetadataDirective: types.MetadataDirectiveCopy,
		},
	)
	return err
}

func (store *Store) DeleteObject(
	ctx context.Context,
	bucket, key string,
) error {
	_, err := store.client.API.DeleteObject(
		ctx,
		&s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		},
	)
	if err != nil {
		return err
	}

	fmt.Println(bucket)
	fmt.Println(key)

	return nil
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
