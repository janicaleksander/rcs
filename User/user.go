package User

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Base information about user in system.
type User struct {
	ID       uuid.UUID
	Email    string
	Password string
	Role     Role // name assigning when adding to unit User(1)-(1)Unit
}

func HashPassword(password string) (string, error) {
	cost := 14
	b, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func DecryptHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

// Roles for each user
// Permissions maybe:
// 0 -> nothing
// 1 -> see things  in unit
// 2 -> see things in all squad
type Role struct {
	name      string
	ruleLevel uint // from 0 to ... describe what each role can do in system
}
