package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func InitDB() *sql.DB {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSLMODE")

	if (user == "") || (dbName == "") {
		log.Fatalf("Error: DB_USER or DB_NAME are not defined in .env file")
	}

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbName, sslMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error opening the database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	fmt.Println("Connection to postgresql established")
	return db
}
