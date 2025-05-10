package suite

import (
	"context"
	"net"
	"strconv"
	"testing"

	"github.com/DENFNC/Zappy/catalog_service/internal/utils/config"
	proto "github.com/DENFNC/Zappy/catalog_service/proto/gen/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Cfg                   *config.Config
	CategoryServiceClient proto.CategoryServiceClient
}

const (
	grpcHost = "localhost"
)

// const configPath

// TODO: Потом поработать этот код чтобы ну нормально выглядело
func New(t *testing.T) (context.Context, *Suite) {

	// Функция будет восприниматься как вспомогательная для тестов
	t.Parallel() // Разрешаем параллельный запуск тестов
	t.Helper()

	// Читаем конфиг из файл

	cfg := config.MustLoad("../../config/config_test.yaml")

	// db, err := postgres.New(cfg.Postgres.URL)
	// require.NoError(t, err, "DB connection failed")
	// // defer db.Close()
	// defer db.DB.Close()

	// _, err = db.DB.Exec(context.Background(), `DELETE FROM categories;`)
	// require.NoError(t, err, "failed to clean categories table")

	// Основной родительский контекст
	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	// Когда тесты пройдут, закрываем контекст
	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	grpcAddress := net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))

	// Создаем клиент
	cc, err := grpc.DialContext(context.Background(),
		grpcAddress,
		// Используем insecure-коннект для тестов
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	categoryServiceClient := proto.NewCategoryServiceClient(cc)

	return ctx, &Suite{
		T:                     t,
		Cfg:                   cfg,
		CategoryServiceClient: categoryServiceClient,
	}
}
