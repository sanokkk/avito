package main

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sanokkk/avito-shop/internal/config"
)

func main() {
	mustMigrateDown()
}

func mustMigrateDown() {
	cfg := config.MustLoad()
	m, err := migrate.New(
		"file://../../migrations",
		cfg.DbConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Down(); err != nil {
		log.Fatal(err)
	}
}
