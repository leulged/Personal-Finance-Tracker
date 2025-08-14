package utils

import (
	"errors"
	"regexp"
)

func IsEmailValid(email string) error {
	if email == "" {
		return errors.New("email is required")
	}
	
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	
	return nil
}

