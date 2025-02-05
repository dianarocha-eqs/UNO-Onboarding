package utils

import (
	"errors"
	"regexp"
)

// Checks if the email is well constructed
func ValidateEmail(email string) error {
	var checkemail = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !checkemail.MatchString(email) || len(email) < 5 {
		return errors.New("invalid email format")
	}
	return nil
}

// Checks if the phone is well constructed
func ValidatePhone(phone string) error {
	var checknumber = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	if !checknumber.MatchString(phone) {
		return errors.New("invalid phone number format")
	}
	return nil
}
