package libraries

import (
	"errors"
	"fmt"
	"kwanjai/config"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	firestoreClient, err := FirebaseApp().Firestore(config.Context)
	defer firestoreClient.Close()
	if err != nil {
		return "Cannot create token uuid.", err
	}
	getUserToken, err := firestoreClient.Collection("tokenUUID").Doc(username).Get(config.Context)
	if !getUserToken.Exists() {
		_, err = firestoreClient.Collection("tokenUUID").Doc(username).Set(config.Context, map[string]interface{}{
			tokenUUID: tokenType,
		})
	}
	_, err = firestoreClient.Collection("tokenUUID").Doc(username).Update(config.Context, []firestore.Update{
		{
			Path:  tokenUUID,
			Value: tokenType,
		},
	})
	if err != nil {
		return "Cannot create token uuid.", err
	}

	claims := &customClaims{
		username,
		tokenUUID,
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(lifetime).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	return signedToken, err
}

// Initialize token method.
func (token *Token) Initialize(username string) (int, string) {
	var passed bool
	tokenStatus := new(tokenStatus)
	go tokenStatus.createToken(username, "access", token)
	go tokenStatus.createToken(username, "refresh", token)
	timeout := time.Now().Add(time.Second * 4)
	timer := time.Now()
	for !passed && timer.Before(timeout) {
		passed = tokenStatus.AccessToken == true && tokenStatus.RefreshToken == true
		timer = time.Now()
	}
	if !passed {
		return http.StatusInternalServerError, "create token error"
	}
	return http.StatusOK, "Token issued."
}

func (tokenStatus *tokenStatus) createToken(username string, tokenType string, token *Token) {
	var err error
	if tokenType == "access" {
		token.AccessToken, err = CreateToken("access", username)
		tokenStatus.AccessToken = err == nil
	} else if tokenType == "refresh" {
		token.RefreshToken, err = CreateToken("refresh", username)
		tokenStatus.RefreshToken = err == nil
	}
}

// VerifyToken function returns token validiation status (bool), username (string), token UUID (string), error.
func VerifyToken(tokenString string, tokenType string) (bool, string, string, error) {
	tokenUUID, valid, err := GetTokenPayload(tokenString, tokenType, "uuid")
	username, _, err := GetTokenPayload(tokenString, tokenType, "user")
	firestoreClient, ferr := FirebaseApp().Firestore(config.Context)
	defer firestoreClient.Close()
	if ferr != nil {
		return false, "anonymous", "", ferr
	}
	uuidVerification, ferr := firestoreClient.Collection("tokenUUID").Doc(username).Get(config.Context)
	if ferr != nil {
		if uuidVerification != nil {
			tokenPath := uuidVerification.Ref.Path
			tokenNotExist := status.Errorf(codes.NotFound, "%q not found", tokenPath)
			if ferr.Error() == tokenNotExist.Error() {
				return false, "anonymous", "", errors.New("token not found")
			}
		}
		return false, "anonymous", "", ferr
	}
	if uuidVerification.Data()[tokenUUID] == nil {
		return false, "anonymous", "", errors.New("token is not valid")
	}
	if err != nil {
		if err.Error() == "Token is expired" {
			ferr = DeleteToken(username, tokenUUID)
			if ferr != nil {
				return false, "anonymous", "", ferr
			}
		}
	}
	if valid {
		return true, username, tokenUUID, nil
	}
	return false, "anonymous", "", err
}

// DeleteToken by username and token uuid.
func DeleteToken(username string, tokenUUID string) error {
	firestoreClient, err := FirebaseApp().Firestore(config.Context)
	defer firestoreClient.Close()
	_, err = firestoreClient.Collection("tokenUUID").Doc(username).Update(config.Context, []firestore.Update{
		{
			Path:  tokenUUID,
			Value: firestore.Delete,
		},
	})
	if err != nil {
		return err
	}
	return err
}
