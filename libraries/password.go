package libraries

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword is the function to generated salted and hashed password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(bytes), err
}

// CheckPasswordHash is the function to compare raw password and hash password.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
