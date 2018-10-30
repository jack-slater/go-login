package handlers

import (
	"errors"
	"fmt"
)

type UserDTO struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	Login     string `json:"login,omitempty"`
	Password  string `json:"password,omitempty"`
}

func (u *UserDTO) valid() error {
	fields := []string{u.FirstName, u.LastName, u.Email, u.Login, u.Password}
	for _, field := range fields {
		if err := validateField(field); err != nil {
			return err
		}
	}
	return nil
}

func validateField(fieldName string) error {
	if len(fieldName) == 0 {
		return errors.New(fmt.Sprintf("%v is empty, all values must be complete", fieldName))
	}
	return nil
}
