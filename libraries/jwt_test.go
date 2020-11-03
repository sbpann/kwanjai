package libraries

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
)

func TestVerifyNonTokenText(t *testing.T) {
	_, _, err := GetTokenPayload("", "access", "")
	assert.Equal(t, err.Error(), "token contains an invalid number of segments")
}

func TestVerifyInvalidSigningMethod(t *testing.T) {
	_, _, err := GetTokenPayload("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.POstGetfAytaZS82wHcjoTyoqhMyxXiWdR7Nn7A29DNSl0EiXLdwJ6xC6AfgZWF1bOsS_TuYI3OG85AmiExREkrS6tDfTQ2B3WXlrr-wp5AokiRbz3_oB4OxG-W9KcEEbDRcZc0nH3L7LzYptiy1PtAylQGxHTWZXtGz4ht0bAecBgmpdgXMguEIcoqPJ1n3pIWk_dUZegpqx0Lka21H6XxUTxiy8OcaarA8zdnPUnV6AmNP3ecFawIFYdvJB_cm-GvpCSbr8G8y_Mllj8f4x9nBH8pQux89_6gUY618iYv7tuPWBFfEbLxtF2pZS6YC1aSfLQxeNe8djT9YjpvRZA", "access", "")
	assert.Equal(t, strings.HasPrefix(err.Error(), "Unexpected signing"), true)
}

func TestVerifyInvalidSignature(t *testing.T) {
	_, _, err := GetTokenPayload("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.Et9HFtf9R3GEMA0IICOfFMVXY7kkTX1wr4qCyhIf58U", "access", "")
	assert.Equal(t, err.Error(), "signature is invalid")
	_, _, err = GetTokenPayload("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.Et9HFtf9R3GEMA0IICOfFMVXY7kkTX1wr4qCyhIf58U", "refresh", "")
	assert.Equal(t, err.Error(), "signature is invalid")
}

func TestExpiredToken(t *testing.T) {
	claims := &customClaims{
		"test",
		uuid.New().String(),
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Second).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("secret"))
	assert.Equal(t, err, nil)
	time.Sleep(2 * time.Second)
	_, err = jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	assert.Equal(t, err.Error(), "Token is expired")
}
