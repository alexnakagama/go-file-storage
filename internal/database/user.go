package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/alexedwards/argon2id"
)

type User struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	IsAdmin        bool      `json:"is_admin"`
	EmailVerified  bool      `json:"email_verified"`
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

func CreateUser(db *sql.DB, name, email, hashedPassword, veritifactionToken string) (*User, error) {
	now := time.Now()
	var user User

	query := `
        INSERT INTO users (name, email, hashed_password, is_admin, email_verified, created_at, updated_at)
        VALUES($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, name, email, is_admin, email_verified, created_at, updated_at
    `
	err := db.QueryRow(query, name, email, hashedPassword, false, false, now, now).
		Scan(&user.ID, &user.Name, &user.Email, &user.IsAdmin, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt)

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

func SearchUser(db *sql.DB, email, password string) (*User, error) {
	var user User
	query := `
		SELECT id, name, email, hashed_password, is_admin, created_at, updated_at
		FROM users WHERE email = $1
	`
	err := db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.HashedPassword,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	} else if err != nil {
		return nil, err
	}

	match, err := argon2id.ComparePasswordAndHash(password, user.HashedPassword)
	if err != nil {
		return nil, err
	}
	if !match {
		fmt.Printf("Incorrect password or email")
		return nil, nil
	}

	user.HashedPassword = ""
	return &user, nil
}
