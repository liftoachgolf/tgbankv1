package pkg

import (
	"database/sql"
	"fmt"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/tgbotdb?sslmode=disable"
)

func NewPostgresDb() (*sql.DB, error) {
	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		return nil, err
	}
	err = db.Ping() // Проверка соединения с базой данных
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil

}
