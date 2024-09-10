package postgres

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Repository
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("cant connect to db, ERR: %s", err.Error())
	}

	testQueries = NewRepository(testDB)
}
