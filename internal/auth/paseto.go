package auth

import (
	"time"

	"github.com/o1egl/paseto"
)

var pasetoInstance = paseto.NewV2()

type CustomClaims struct {
	UserID int
	Email  string
	Exp    time.Time
}

// todo: put it in env file, genmerate a more secure key
var symetricKey = []byte("12345678901234567890123456789012")

func GenerateToken(key []byte, userID int, email string, duration time.Duration) (string, error) {
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

func VerifyToken(token string, key []byte) (*CustomClaims, error) {

}
