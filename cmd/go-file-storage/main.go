package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/alexnakagama/go-file-storage/internal/auth"
	"github.com/alexnakagama/go-file-storage/internal/handler"
	"github.com/alexnakagama/go-file-storage/internal/middleware"
)

func main() {
	// reads the .env file
	godotenv.Load()
	auth.InitPaseto()

	connStr := os.Getenv("DATABASE_URL")

	// opens db connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// endpoints
	http.HandleFunc("/register", handler.RegisterUserHandler(db))
	http.HandleFunc("/login", handler.LoginUserHandler(db))
	http.HandleFunc("/verify-email", handler.VerifyEmailHandler(db))

	// protected endpoints
	http.Handle("/delete", middleware.AuthMiddleware(http.HandlerFunc(handler.DeleteUserHandler(db))))
	http.Handle("/files", middleware.AuthMiddleware(handler.AddFileHandler(db, "./uploads")))

	log.Println("Server running in port:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
