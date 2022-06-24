package postgres

import (
	"time"
)

// User is the structure which holds one user from the database.
type User struct {
	ID        int
	Email     string
	FirstName string
	LastName  string
	Password  string
	Active    int
	CreatedAt time.Time
	UpdatedAt time.Time
}
