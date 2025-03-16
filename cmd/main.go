package main

import (
	"database/sql"
	"log"

	"github.com/hisshihi/order-of-venhicles-services/internal/config"
	"github.com/hisshihi/order-of-venhicles-services/internal/db"
	api "github.com/hisshihi/order-of-venhicles-services/internal/db/api/handlers"

	_ "github.com/lib/pq"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config file", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start to server:", err)
	}
}
