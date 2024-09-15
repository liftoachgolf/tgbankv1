package postgrestg

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // подключение драйвера для PostgreSQL
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

func NewPostgresDb() (*sql.DB, error) {
	log.Println("Opening database connection...")
	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	log.Println("Pinging database...")
	err = db.Ping() // Проверка соединения с базой данных
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}
