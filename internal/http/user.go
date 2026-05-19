package http

import (
	"database/sql"
	"encoding/json"
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
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid json", http.StatusBadRequest)
			return
		}

		if req.Name == "" || req.Email == "" || req.Password == "" || req.ConfirmPassword == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		if req.Password != req.ConfirmPassword {
			http.Error(w, "Passwords doesnt match", http.StatusBadRequest)
			return
		}

		// todo hash password
	}
}

func LoginHandler(db *sql.DB) http.HandlerFunc {

}

func DeleteUserHandler(db *sql.DB) http.HandlerFunc {

}
