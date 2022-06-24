package checking

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User is the structure which holds one user from the database.
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"-"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) PasswordMatches(password *string) (bool, error) {
	// decode password
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(*password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, err
	} else if err != nil {
		return false, err
	}

	return true, nil
}
