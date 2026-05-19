package http

import (
	"database/sql"
	"net/http"
)

type RegisterRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

	}
}

func LoginHandler(db *sql.DB) http.HandlerFunc {

}

func DeleteUserHandler(db *sql.DB) http.HandlerFunc {

}
