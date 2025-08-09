package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func CreateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
			"exp":   time.Now().Add(time.Hour * 12).Unix(),
		})
	tokenStr, err := token.SignedString(os.Getenv("JWT_KEY"))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func VerifyToken(tokenStr string) (bool, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return os.Getenv("JWT_KEY"), nil
	})
	if err != nil {
		return false, err
	}
	if !token.Valid {
		return false, errors.New("jwt token expired")
	}
	return true, nil
}
