package checking

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/djedjethai/authentication/pkg/http/dto"
	"github.com/djedjethai/authentication/pkg/storage/postgres"
	"net/http"
)

type Service interface {
	GetByEmail(*dto.RequestPayload) (*User, error)
}

// type Repository interface {
// 	GetAll() ([]*User, error)
// 	GetByEmail(email string) (*User, error)
// 	GetOne(id int) (*User, error)
// 	PasswordMatches(plainText string, u User) (bool, error)
// }

type service struct {
	// r postgres.Repository
	r postgres.Repository
}

// func NewService(r postgres.Repository) Service {
func NewService(r postgres.Repository) Service {
	return &service{r}
}

func (s *service) GetByEmail(rp *dto.RequestPayload) (*User, error) {
	var user User

	u, err := s.r.GetByEmail(rp.Email)
	if err != nil {
		return nil, err
	}

	// convert user from db to user
	user.ID = u.ID
	user.Email = u.Email
	user.FirstName = u.FirstName
	user.LastName = u.LastName
	user.Password = u.Password
	user.Active = u.Active
	user.CreatedAt = u.CreatedAt
	user.UpdatedAt = u.UpdatedAt

	valid, err := user.PasswordMatches(&rp.Password)
	if err != nil || !valid {
		return nil, err
	}

	// log the request
	err = logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil || !valid {
		return nil, err
	}

	return &user, nil
}

func logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
