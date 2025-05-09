package main

import (
	"github.com/DENFNC/Zappy/order_service/internal/adapters/sql/postgres"
	"github.com/DENFNC/Zappy/order_service/internal/utils/config"
)

func main() {
	cfg := config.MustLoad("./config/config.yaml")
	// logger, err := logger.New(cfg.LogType)
	// if err != nil {
	// 	panic(err)
	// }

	dbpool, err := postgres.New(cfg.Postgres.URL)
	if err != nil {
		panic(err)
	}
}
