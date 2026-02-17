package database

import (
	"database/sql"
	"unitrip/internal/config"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func NewDB(config *config.Config) *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DatabaseHost, config.DatabasePort, config.DatabaseUser, config.DatabasePassword, config.DatabaseName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Failed to open database: %v\n", err)
		return nil
	}

	if err := db.Ping(); err != nil {
		log.Printf("Failed to ping database: %v\n", err)
		return nil
	}

	log.Println("Database connected successfully")
	return db
}
