package token

import (
	"errors"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func getSecretKey() (string, error) {
	config := struct {
		General struct {
			JwtKey string `toml:"jwtKey"`
		}
	}{}
	_, err := toml.DecodeFile("configproduction/general.toml", &config)
	if err != nil {
		return "", err
	}
	return config.General.JwtKey, nil
}

func CreateToken(id, email string) (string, error) {
	key, err := getSecretKey()
	if err != nil {
		return "", err
	}

	claims := UserClaims{
		ID:    id,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}

func VerifyToken(tokenStr string) (*UserClaims, error) {
	key, err := getSecretKey()
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
