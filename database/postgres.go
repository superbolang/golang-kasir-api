package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	// Open DB
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Print("Can not open postgresql")
		return nil, err
	}

	// Test connection
	if err = db.Ping(); err != nil {
		log.Print("Can not connect postgresql")
		return nil, err
	}

	// Connection pool
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(25)
	return db, nil
}
