package libraries

import (
	"context"
	"errors"
	"fmt"
	"kwanjai/config"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// Token object
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type customClaims struct {
	User string `json:"user"`
	UUID string `json:"uuid,omitempty"`
	*jwt.StandardClaims
}

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
func createToken(tokenType string, username string) (string, error) {
	secretKey, lifetime, err := getSecretKeyAndLifetime(tokenType)
	var tokenUUID string
	if err != nil {
		return "no token created.", err
	}

	tokenUUID = uuid.New().String()
	ctx := context.Background()
	firestoreClient, err := FirebaseApp().Firestore(ctx)
	defer firestoreClient.Close()
	if err != nil {
		return "Cannot create token uuid.", err
	}
	getUserToken, err := firestoreClient.Collection("tokenUUID").Doc(username).Get(ctx)
	if !getUserToken.Exists() {
		_, err = firestoreClient.Collection("tokenUUID").Doc(username).Set(ctx, map[string]interface{}{
			tokenUUID: tokenType,
		})
	}
	_, err = firestoreClient.Collection("tokenUUID").Doc(username).Update(ctx, []firestore.Update{
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

// Initialize token
func (token *Token) Initialize(username string) (int, string) {
	var err error
	token.AccessToken, err = createToken("access", username)
	token.RefreshToken, err = createToken("refresh", username)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusOK, "Token issued."
}

// VerifyToken with a particular type.
func VerifyToken(tokenString string, tokenType string) (bool, string, error) {
	var secretKey string
	if tokenType == "access" {
		secretKey = config.JWTAccessTokenSecretKey
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, "anonymous", errors.New("claim failed")
	}
	tokenUUID, username := claims["uuid"].(string), claims["user"].(string)
	ctx := context.Background()
	firestoreClient, ferr := FirebaseApp().Firestore(ctx)
	defer firestoreClient.Close()
	if ferr != nil {
		return false, "", ferr
	}
	uuidVerification, ferr := firestoreClient.Collection("tokenUUID").Doc(username).Get(ctx)
	if ferr != nil {
		return false, "", ferr
	}
	if uuidVerification.Data()[tokenUUID] == nil {
		return false, "anonymous", errors.New("token is not valid")
	}
	if err.Error() == "Token is expired" {
		_, ferr = firestoreClient.Collection("tokenUUID").Doc(username).Update(ctx, []firestore.Update{
			{
				Path:  tokenUUID,
				Value: firestore.Delete,
			},
		})
		if ferr != nil {
			return false, "", ferr
		}
	}
	if ok && token.Valid {
		return true, username, nil
	}
	return false, "anonymous", err
}
