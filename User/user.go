package User

import "github.com/google/uuid"

// Base information about user in system.
type User struct {
	ID       uuid.UUID
	Email    string
	Password string
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
