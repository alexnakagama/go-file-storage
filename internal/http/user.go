package http

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alexnakagama/go-file-storage/internal/database"
	"github.com/alexnakagama/go-file-storage/utils"
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

		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		user, err := database.CreateUser(db, req.Name, req.Email, hashedPassword)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not create user: %v", err), http.StatusInternalServerError)
			return
		}

		user.HashedPassword = ""

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func LoginHandler(db *sql.DB) http.HandlerFunc {

}

func DeleteUserHandler(db *sql.DB) http.HandlerFunc {

}
