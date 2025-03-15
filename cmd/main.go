package main

import (
	"database/sql"
	"log"

	"github.com/hisshihi/order-of-venhicles-services/internal/db"
	api "github.com/hisshihi/order-of-venhicles-services/internal/db/api/handlers"

	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/order_services?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start to server:", err)
	}
}
