package auth

import (
	"errors"
	"os"
	"time"

	"github.com/o1egl/paseto"
)

type CustomClaims struct {
	UserID int       `json:"user_id"`
	Email  string    `json:"email"`
	Exp    time.Time `json:"exp"`
}

var symetricKey []byte

func InitPaseto() {
	key := os.Getenv("PASETO_SYMETRIC_KEY")
	if len(key) != 32 {
		panic("PASETO_SYMETRIC_KEY must have 32 characters")
	}
	symetricKey = []byte(key)
}

func GenerateToken(userID int, email string, duration time.Duration) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		UserID: userID,
		Email:  email,
		Exp:    now.Add(duration),
	}

	token, err := paseto.NewV2().Encrypt(symetricKey, claims, nil)
	if err != nil {
		return "", err
	}

	return token, nil
}

func VerifyToken(token string) (*CustomClaims, error) {
	var newClaims CustomClaims

	err := paseto.NewV2().Decrypt(token, symetricKey, &newClaims, nil)
	if err != nil {
		return nil, err
	}

	if time.Now().After(newClaims.Exp) {
		return nil, errors.New("Expired token")
	}

	return &newClaims, nil
}
