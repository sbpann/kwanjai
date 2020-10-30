package libraries

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword function generates salted and hashed password, returns password (string) and error.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(bytes), err
}

// CheckPasswordHash function compares raw password and hash password, returns validiation status (bool).
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
