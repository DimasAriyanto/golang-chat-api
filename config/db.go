package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB(config Config) *sql.DB {
	dsn := config.GetDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Database is not responding: %v", err)
	}

	log.Println("Connected to MySQL database successfully")
	return db
}