package models

import (
	"fmt"
	"kwanjai/config"
	"kwanjai/libraries"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoginCredential info
type LoginCredential struct {
	ID       string `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LogoutData struct {
	RefreshPass      bool
	AccessPass       bool
	User             string
	AccessTokenUUID  string
	RefreshTokenUUID string
	AccessDeleted    bool
	RefreshDeleted   bool
	Err              error
}

type authenticationPerform interface {
	login() (int, string)
}

// Login user
func Login(perform authenticationPerform) (int, string) {
	return perform.login()
}

func GetUserSession(username string) (int, string) {
	firestoreClient, err := libraries.FirebaseApp().Firestore(config.Context)
	defer firestoreClient.Close()
	uuidVerification, err := firestoreClient.Collection("tokenUUID").Doc(username).Get(config.Context)
	fmt.Println(uuidVerification.Data(), err)
	return http.StatusOK, "get user successfully"
}

func (logout *LogoutData) Verify(tokenString string, tokenType string) {
	if tokenType == "access" {
		logout.AccessPass, logout.User, logout.AccessTokenUUID, logout.Err = libraries.VerifyToken(tokenString, "access")
	} else if tokenType == "refresh" {
		logout.RefreshPass, logout.User, logout.RefreshTokenUUID, logout.Err = libraries.VerifyToken(tokenString, "refresh")
	} else {
		return
	}
}

func (login *LoginCredential) login() (int, string) {
	hashedPassword := ""
	username := ""
	login.ID = strings.ToLower(login.ID)
	firestoreClient, err := libraries.FirebaseApp().Firestore(config.Context)
	defer firestoreClient.Close()
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	getUser, err := firestoreClient.Collection("users").Doc(login.ID).Get(config.Context)
	if err != nil {
		userPath := getUser.Ref.Path
		userNotExist := status.Errorf(codes.NotFound, "%q not found", userPath)
		if err.Error() == userNotExist.Error() {
			findEmail := firestoreClient.Collection("users").Where("Email", "==", login.ID).Documents(config.Context)
			foundEmail, err := findEmail.GetAll()
			if err != nil {
				return http.StatusInternalServerError, err.Error()
			}
			if len(foundEmail) > 0 {
				hashedPassword = foundEmail[0].Data()["HashedPassword"].(string)
				username = foundEmail[0].Data()["Username"].(string)
			} else {
				return http.StatusBadRequest, "Cannot login with provided credential."
			}
		} else {
			return http.StatusInternalServerError, err.Error()
		}
	} else {
		hashedPassword = getUser.Data()["HashedPassword"].(string)
		username = getUser.Data()["Username"].(string)
	}
	passwordPass := libraries.CheckPasswordHash(login.Password, hashedPassword)
	if !passwordPass {
		return http.StatusBadRequest, "Cannot login with provided credential."
	}
	_, err = firestoreClient.Collection("users").Doc(username).Update(config.Context, []firestore.Update{{Path: "IsActive", Value: true}})
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusOK, username
}
