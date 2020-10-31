package libraries

import (
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
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
