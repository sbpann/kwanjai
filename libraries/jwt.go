package libraries

import (
	"errors"
	"fmt"
	"kwanjai/config"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// Token object.
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Use for tracking token issue status.
type tokenStatus struct {
	AccessToken  bool
	RefreshToken bool
}

type customClaims struct {
	User string `json:"user"`
	UUID string `json:"uuid,omitempty"`
	*jwt.StandardClaims
}

// getSecretKeyAndLifetime function returns secret key (string), token lifetime (time.Duration), error.
func getSecretKeyAndLifetime(tokenType string) (string, time.Duration, error) {
	if tokenType == "access" {
		return config.JWTAccessTokenSecretKey, config.JWTAccessTokenLifetime, nil
	} else if tokenType == "refresh" {
		return config.JWTRefreshTokenSecretKey, config.JWTRefreshTokenLifetime, nil
	} else {
		err := errors.New("no token type provide")
		return "", time.Second, err
	}
}

// GetTokenPayload function returns payload value (string), token validiation status (bool), error.
func GetTokenPayload(tokenString string, tokenType string, field string) (string, bool, error) {
	secretKey, _, err := getSecretKeyAndLifetime(tokenType)
	if err != nil {
		return "", false, err
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", false, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", false, errors.New("claim failed")
	}
	return claims[field].(string), token.Valid, nil
}

// CreateToken returns token (string) and error.
func CreateToken(tokenType string, username string) (string, error) {
	secretKey, lifetime, err := getSecretKeyAndLifetime(tokenType)
	var tokenUUID string
	if err != nil {
		return "no token created.", err
	}

	tokenUUID = uuid.New().String()
	now := time.Now().Truncate(time.Millisecond)

	_, err = FirestoreCreatedOrSet("tokenUUID", tokenUUID,
		map[string]interface{}{
			"user":   username,
			"expire": now,
		})
	if err != nil {
		return "Cannot store token data.", err
	}
	claims := &customClaims{
		username,
		tokenUUID,
		&jwt.StandardClaims{
			ExpiresAt: now.Add(lifetime).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	return signedToken, err
}

// Initialize token method.
func (token *Token) Initialize(username string) (int, string) {

	tokenStatus := new(tokenStatus)
	accessPassed := make(chan bool)
	refreshPassed := make(chan bool)
	go tokenStatus.createToken(username, "access", token, accessPassed)
	go tokenStatus.createToken(username, "refresh", token, refreshPassed)
	passed := true == <-accessPassed && true == <-refreshPassed
	if !passed {
		return http.StatusInternalServerError, "create token error"
	}
	return http.StatusOK, "Token issued."
}

func (tokenStatus *tokenStatus) createToken(username string, tokenType string, token *Token, passed chan bool) {
	var err error
	if tokenType == "access" {
		token.AccessToken, err = CreateToken("access", username)
		passed <- err == nil
	} else if tokenType == "refresh" {
		token.RefreshToken, err = CreateToken("refresh", username)
		passed <- err == nil
	}
}

// VerifyToken function returns token validiation status (bool), username (string), token UUID (string), error.
func VerifyToken(tokenString string, tokenType string) (bool, string, string, error) {
	tokenUUID, valid, err := GetTokenPayload(tokenString, tokenType, "uuid")
	username, _, err := GetTokenPayload(tokenString, tokenType, "user")
	if err != nil {
		if err.Error() == "Token is expired" {
			FirestoreDelete("tokenUUID", tokenUUID)
		}
		return false, "anonymous", "", err
	}
	uuidVerification, err := FirestoreFind("tokenUUID", tokenUUID)
	if !uuidVerification.Exists() {
		return false, "anonymous", "", err
	}
	if valid {
		return true, username, tokenUUID, nil
	}
	return false, "anonymous", "", err
}
