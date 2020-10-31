package libraries

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestInvalidPassword(t *testing.T) {
	password := "password"
	hashedpassword, _ := HashPassword(password)
	assert.Equal(t, CheckPasswordHash("wrongpassword", hashedpassword), false)
}
