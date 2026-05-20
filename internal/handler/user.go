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

type RegisterResponse struct {
	Message string         `json:"message"`
	User    *database.User `json:"user"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string         `json:"message"`
	User    *database.User `json:"user"`
	Token   string         `json:"token"`
}

type DeleteRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DeleteUserResponse struct {
	Message string         `json:"message"`
	User    *database.User `json:"user"`
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

		resp := RegisterResponse{
			Message: "User created successfully. Verify your email before logging in",
			User:    nil,
		}

		json.NewEncoder(w).Encode(resp)
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

		if !user.EmailVerified {
			http.Error(w, "You need to verify your email", http.StatusUnauthorized)
			return
		}

		token, err := auth.GenerateToken(user.ID, user.Email, time.Hour)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		user.HashedPassword = ""

		resp := LoginResponse{
			Message: "Login successfull",
			User:    user,
			Token:   token,
		}

		w.Header().Set("Content-Type", "applicaction/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func DeleteUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
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

		resp := DeleteUserResponse{
			Message: "User deleted successfully",
			User:    deletedUser,
		}

		json.NewEncoder(w).Encode(resp)
	}
}

func VerifyEmailHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "Token missing", http.StatusBadRequest)
			return
		}

		user, err := database.FindUserByVerificationToken(db, token)
		if err != nil || user == nil {
			http.Error(w, "Invalid or expired token", http.StatusBadRequest)
			return
		}

		err = database.MarkEmailAsVerified(db, user.ID)
		if err != nil {
			http.Error(w, "Error verifying email", http.StatusInternalServerError)
			return
		}

		fmt.Println(w, "Email verified, now you can log in!")
	}
}
