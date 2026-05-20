package main

import (
	"database/sql"
	"log"
	"net/http"

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

	// opens db connection
	db, err := sql.Open("postgres", "")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// endpoints
	http.HandleFunc("/register", handler.RegisterUserHandler(db))
	http.HandleFunc("/login", handler.LoginUserHandler(db))
	http.Handle("/delete", middleware.AuthMiddleware(http.HandlerFunc(handler.DeleteUserHandler(db))))

	log.Println("Server running in port:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
