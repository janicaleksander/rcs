package User

import (
	"golang.org/x/crypto/bcrypt"
)

// Base information about user in system.
//this is in proto
/*
type User struct {
	ID        uuid.UUID
	Email     string
	Password  string
	RuleLevel int
}
*/
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
	return err == nil
}

// TODO maybe add role to PersonToUnit
type Role struct {
	name      string
	ruleLevel uint // from 0 to ... describe what each role can do in system
}
