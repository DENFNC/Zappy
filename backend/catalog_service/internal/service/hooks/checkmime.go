package hooks

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/DENFNC/Zappy/catalog_service/internal/domain/models"
	"github.com/DENFNC/Zappy/catalog_service/internal/utils/config"
)

type Hooks struct {
	ProductImageRepo
	ObjectStorage
	KeyValueStorage
	log *slog.Logger
	cfg *config.Config
}

func New(
	storage ProductImageRepo,
	log *slog.Logger,
	objStorage ObjectStorage,
	kvStorage KeyValueStorage,
	cfg *config.Config,
) *Hooks {
	return &Hooks{
		ProductImageRepo: storage,
		ObjectStorage:    objStorage,
		KeyValueStorage:  kvStorage,
		log:              log,
		cfg:              cfg,
	}
}

func (svc *Hooks) CheckMime(
	ctx context.Context,
	bucket, key string,
	byteRange string,
) error {
	const op = "service.hooks.Hooks.CheckMime"

	log := svc.log.With("op", op)

	buffer := make([]byte, 512)
	object, err := svc.ObjectStorage.GetObjectRange(ctx, bucket, key, byteRange)
	if err != nil {
		return err
	}
	n, err := io.ReadFull(object, buffer)
	if err != nil && err != io.ErrUnexpectedEOF {
		return err
	}
	mimeType := http.DetectContentType(buffer[:n])

	cfg, ok := svc.cfg.ObjectStore.Buckets[bucket]
	if !ok {
		log.Warn(
			"bucket not found",
			slog.String("bucket", bucket),
			slog.String("key", key),
		)
		return svc.markAsFailed(ctx, bucket, key, log, "bucket not found")
	}

	if !slices.Contains(cfg.MimeTypes, mimeType) {
		log.Debug(
			"invalid mime type",
			slog.String("bucket", bucket),
			slog.String("key", key),
			slog.String("MIME", mimeType),
		)
		return svc.markAsFailed(ctx, bucket, key, log, "invalid mime type")
	}

	fields, err := svc.KeyValueStorage.HGetAll(ctx, key)
	if err != nil {
		return err
	}

	srcBucket := fields["bucket"]
	productID := fields["product_id"]
	alt := fields["alt"]

	path := cfg.Path
	url := keyBuilder(bucket, path, key)

	srcKey := svc.cfg.ObjectStore.Buckets[srcBucket].Path + key

	if err := svc.ObjectStorage.CopyObject(ctx, srcBucket, srcKey, url); err != nil {
		return err
	}
	if err := svc.ObjectStorage.DeleteObject(ctx, bucket, key); err != nil {
		return err
	}
	if err := svc.KeyValueStorage.HSet(ctx, key, map[string]any{"status": "success"}); err != nil {
		return err
	}
	if _, err := svc.ProductImageRepo.Create(ctx, &models.ProductImage{
		ImageID:   key,
		ProductID: productID,
		URL: urlBuilder(
			svc.cfg.ObjectStore.ObjectOrigin,
			svc.cfg.ObjectStore.Buckets[srcBucket].Name,
			svc.cfg.ObjectStore.Buckets[srcBucket].Path,
			key,
		),
		ALT:       alt,
		ObjectKey: svc.cfg.ObjectStore.Buckets[srcBucket].Path + key,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}); err != nil {
		return err
	}

	return nil
}

func keyBuilder(bucket, path, key string) string {
	cleanPath := strings.Trim(path, "/")
	cleanKey := strings.TrimLeft(key, "/")

	var full string
	if cleanPath != "" {
		full = cleanPath + "/" + cleanKey
	} else {
		full = cleanKey
	}

	segments := strings.Split(full, "/")

	var b strings.Builder
	b.WriteString(bucket)
	b.WriteByte('/')

	for i, seg := range segments {
		b.WriteString(url.PathEscape(seg))
		if i < len(segments)-1 {
			b.WriteByte('/')
		}
	}

	return b.String()
}

func urlBuilder(origin, bucket, path, key string) string {
	var b strings.Builder
	b.WriteString(origin)
	b.WriteString(bucket)
	b.WriteString(path)
	b.WriteString(key)
	return b.String()
}

func (svc *Hooks) markAsFailed(
	ctx context.Context,
	bucket, key string,
	log *slog.Logger,
	reason string,
) error {
	if err := svc.ObjectStorage.DeleteObject(ctx, bucket, key); err != nil {
		log.Warn(
			"failed to delete object",
			slog.String("error", err.Error()),
		)
	}
	if err := svc.KeyValueStorage.HSet(ctx, key, map[string]any{"status": "failed"}); err != nil {
		log.Warn(
			"failed to set status",
			slog.String("error", err.Error()),
		)
	}
	return errors.New(reason)
}
