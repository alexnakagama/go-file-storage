package database

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	IsAdmin        bool      `json:"is_admin"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func CreateUser(db *sql.DB, name, email, hashedPassword string, isAdmin bool) (*User, error) {
	now := time.Now()
	var user User

	query := `
		INSERT INTO users (name, email, hashed_password, is_admin, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6)
		RETURNING id, name, email, is_admin, created_at, updated_at
	`
	err := db.QueryRow(query, name, email, hashedPassword, isAdmin, now, now).
		Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	fmt.Println("User has been created successfully")
	return &user, nil
}

func DeleteUser(db *sql.DB, id int) (*User, error) {
	var user User
	query := `
		DELETE FROM users WHERE id=$1
		RETURNING id, name, email, is_admin, created_at, updated_at
	`

	err := db.QueryRow(query, id).
		Scan(&user.ID, &user.Name, &user.Email, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	fmt.Println("User has been deleted")
	return &user, nil
}
