package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/DENFNC/Zappy/internal/config"
	"github.com/DENFNC/Zappy/internal/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // postgres driver
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const emptyValue = 0

func main() {
	const op = "migrate"

	// Create logger (assuming your logger.New returns a logger with Info and Error methods)
	log, err := logger.New("dev")
	if err != nil {
		panic(err)
	}
	log = log.With("op", op)

	// Global flags
	envPath := flag.String("env-path", "", "Path to the configuration file")
	migPath := flag.String("mig-path", "", "Path to the migrations folder")
	version := flag.Uint("mig-ver", emptyValue, "Target migration version (if needed)")
	flag.Parse()

	// Check if a command is provided as the first positional argument
	if flag.NArg() < 1 {
		fmt.Println("Command required. Available commands: drop, down, up, version")
		os.Exit(1)
	}

	command := flag.Arg(0)

	// Load configuration
	cfg := config.MustLoad(*envPath)

	// Create migrator instance
	migrator, err := migrate.New(
		fmt.Sprintf("file://%s", *migPath),
		cfg.Postgres.URL,
	)
	if err != nil {
		fmt.Println(err)
		panic("Failed to create migrator")
	}

	// Execute command using switch
	switch command {
	case "drop":
		log.Info("Dropping all migrations")
		if err := migrator.Drop(); err != nil {
			log.Error("Failed to drop migrations", "error", err)
			os.Exit(1)
		}
	case "down":
		log.Info("Rolling back migrations")
		if err := migrator.Down(); err != nil {
			log.Error("Failed to roll back migrations", "error", err)
			os.Exit(1)
		}
	case "up":
		log.Info("Applying migrations")
		if err := migrator.Up(); err != nil {
			log.Error("Failed to apply migrations", "error", err)
			os.Exit(1)
		}
	case "version":
		if *version == emptyValue {
			fmt.Println("For the 'version' command, please provide -mig-ver flag, e.g., -mig-ver=2022051801")
			os.Exit(1)
		}
		log.Info("Switching to migration version", "version", strconv.Itoa(int(*version)))
		if err := migrator.Migrate(*version); err != nil {
			log.Error("Failed to switch migration version", "error", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown command: %s. Available commands: drop, down, up, version\n", command)
		os.Exit(1)
	}
}
