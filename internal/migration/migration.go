package migration

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sanokkk/avito-shop/internal/config"
)

func MustMigrate() {
	cfg := config.MustLoad()
	var migrationsPath string
	if os.Getenv("ENV") == "docker" {
		migrationsPath = "migrations"
	} else {
		migrationsPath = "//../../migrations"
	}
	m, err := migrate.New(
		fmt.Sprintf("file:%s", migrationsPath),
		cfg.DbConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Выполнены все миграции")
			return
		}

		log.Fatal(err)
	}
}
