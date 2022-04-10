package models

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// ListUsers is simply a list of users
type ListUsers []User

// User is the data model of the user shared by processor and manager
type User struct {
	Country   string `json:"country,omitempty"`
	Email     string `json:"email,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	NickName  string `json:"nickName,omitempty"`
	Password  string `json:"password,omitempty"`
	LastName  string `json:"second_name,omitempty"`
	ID        string `json:"id,omitempty"`
}

// ValidateUserInput is an example of how to perform validation of model.
func (m *User) ValidateUserInput() error {
	if m.ID == "" {
		return errors.New("user id is required")
	}

	if _, err := uuid.Parse(m.ID); err != nil {
		return fmt.Errorf("uid cannot be parsed: %w", err)
	}

	return nil
}
