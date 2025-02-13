package main

import (
	"log"

	"github.com/go-pg/pg/extra/pgdebug"
	"github.com/go-pg/pg/v10"
	"github.com/sanokkk/avito-shop/internal/app"
	"github.com/sanokkk/avito-shop/internal/config"
)

func main() {
	cfg := config.MustLoad()
	dbOptions, err := pg.ParseURL(cfg.DbConnectionString)
	if err != nil {
		log.Fatal("Ошибка при открытии коннекта к БД", err)
	}

	db := pg.Connect(dbOptions)
	db.AddQueryHook(pgdebug.DebugHook{
		// Print all queries.
		Verbose: true,
	})
	app := app.CreateServer(db)

	app.Start()
}
