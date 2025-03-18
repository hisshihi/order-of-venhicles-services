package sqlc

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/hisshihi/order-of-venhicles-services/internal/config"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config file", err)
	}
	testDB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
