package utils

import "golang.org/x/crypto/bcrypt"

// Contains checks if a string exists in a slice of strings
func Contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return v == str
		}
	}
	return false
}

// HashPassword takes a password string and returns the bcrypt hash of it
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
