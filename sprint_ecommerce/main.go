package main

import (
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ramadhan1445sprint/sprint_ecommerce/config"
	"github.com/ramadhan1445sprint/sprint_ecommerce/pkg/database"
	"github.com/ramadhan1445sprint/sprint_ecommerce/server"
)

func main() {
	// config
	config.LoadConfig(".env")

	// db connection
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// generate new server instance
	s := server.NewServer(db)
	s.RegisterRoute()

	// run server
	s.Run()
}
