package config_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/DENFNC/Zappy/catalog_service/internal/utils/config"
	"github.com/stretchr/testify/assert"
)

func writeTempFile(t *testing.T, content string) string {
	tmpDir := os.TempDir()
	f, err := ioutil.TempFile(tmpDir, "cfg-*.yaml")
	assert.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(content)
	assert.NoError(t, err)
	return f.Name()
}

func TestMustLoad_PanicOnEmptyPath(t *testing.T) {
	t.Run("empty path", func(t *testing.T) {
		assert.PanicsWithValue(t, "path is empty", func() {
			config.MustLoad("")
		})
	})
}

func TestMustLoad_PanicOnNotExist(t *testing.T) {
	t.Run("file does not exist", func(t *testing.T) {
		assert.PanicsWithValue(t, "file not found", func() {
			config.MustLoad("/no/such/file.yaml")
		})
	})
}

func TestMustLoad_PanicOnInvalidConfig(t *testing.T) {
	path := writeTempFile(t, "invalid: [ : ]")
	assert.PanicsWithValue(t, "error reading config", func() {
		config.MustLoad(path)
	})
}

func TestMustLoad_Success(t *testing.T) {
	tmpl := `
log_type: prod
paginate_secret: topsecret
grpc:
  port: 50051
  timeout: 2s
  reflection: true
http:
  port: 8080
postgres:
  url: postgres://user:pass@localhost:5432/db
object_store:
  aws_bucket_image: img-bucket
  aws_staging_bucket: staging-bucket
  buckets:
    thumb:
      name: thumb-bucket
      mime_types: ["image/png","image/jpeg"]
      path: /thumbs
`
	path := writeTempFile(t, tmpl)
	cfg := config.MustLoad(path)

	assert.Equal(t, "prod", cfg.LogType)
	assert.Equal(t, "topsecret", cfg.PaginateSecret)
	assert.Equal(t, 50051, cfg.GRPC.Port)
	assert.Equal(t, 2*time.Second, cfg.GRPC.Timeout)
	assert.True(t, cfg.GRPC.Reflection)
	assert.Equal(t, 8080, cfg.HTTP.Port)
	assert.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.Postgres.URL)

	assert.Equal(t, "http://localhost:9000", cfg.ObjectStore.ObjectOrigin)
	assert.Equal(t, "img-bucket", cfg.ObjectStore.ImageBucket)
	assert.Equal(t, "staging-bucket", cfg.ObjectStore.StagingBucket)

	b, ok := cfg.ObjectStore.Buckets["thumb"]
	assert.True(t, ok)
	assert.Equal(t, "thumb-bucket", b.Name)
	assert.ElementsMatch(t, []string{"image/png", "image/jpeg"}, b.MimeTypes)
	assert.Equal(t, "/thumbs", b.Path)
}
