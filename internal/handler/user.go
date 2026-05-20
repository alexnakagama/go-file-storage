package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/alexnakagama/go-file-storage/internal/auth"
	"github.com/alexnakagama/go-file-storage/internal/database"
	"github.com/alexnakagama/go-file-storage/internal/middleware"
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

type DeleteRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterUserHandler(db *sql.DB) http.HandlerFunc {
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

		verificationToken := uuid.New().String()

		user, err := database.CreateUser(
			db,
			req.Name,
			req.Email,
			hashedPassword,
			verificationToken,
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not create user: %v", err), http.StatusInternalServerError)
			return
		}

		verificationURL := fmt.Sprintf("http://localhost:8080/verify-email?token=%s", verificationToken)
		fmt.Println("Verify your email visiting:", verificationURL)
		// Here send email in prod

		user.HashedPassword = ""

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "User created successfully. Verify your email before logging in",
		})
	}
}

func LoginUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
		}

		if req.Email == "" || req.Password == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		user, err := database.SearchUser(db, req.Email, req.Password)
		if err == sql.ErrNoRows || user == nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		token, err := auth.GenerateToken(user.ID, user.Email, time.Hour)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		user.HashedPassword = ""

		resp := struct {
			User  *database.User `json:"user"`
			Token string         `json:"token"`
		}{
			User:  user,
			Token: token,
		}

		w.Header().Set("Content-Type", "applicaction/json")
		json.NewEncoder(w).Encode(&resp)
	}
}

func DeleteUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		claims, ok := middleware.GetUserClaims(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		deletedUser, err := database.DeleteUser(db, claims.UserID)
		if err != nil {
			http.Error(w, "Failed to delete user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(deletedUser)
	}
}
