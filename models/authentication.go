package models

import (
	"kwanjai/libraries"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoginCredential info.
type LoginCredential struct {
	ID       string `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LogoutData to track logout process event.
type LogoutData struct {
	RefreshPassed    bool
	AccessPassed     bool
	User             string
	AccessTokenUUID  string
	RefreshTokenUUID string
	AccessDeleted    bool
	RefreshDeleted   bool
}

type authenticationPerform interface {
	login() (int, string)
}

// Login user.
func Login(perform authenticationPerform) (int, string) {
	return perform.login()
}

// Verify logout credential.
func (logout *LogoutData) Verify(tokenString string, tokenType string, passed chan bool, UUID chan string) {
	if tokenType == "access" {
		logout.AccessPassed, logout.User, logout.AccessTokenUUID, _ = libraries.VerifyToken(tokenString, "access")
		passed <- logout.AccessPassed
		UUID <- logout.AccessTokenUUID
	} else if tokenType == "refresh" {
		logout.RefreshPassed, logout.User, logout.RefreshTokenUUID, _ = libraries.VerifyToken(tokenString, "refresh")
		passed <- logout.RefreshPassed
		UUID <- logout.RefreshTokenUUID
	} else {
		return
	}
}

func (login *LoginCredential) login() (int, string) {
	hashedPassword := ""
	username := ""
	login.ID = strings.ToLower(login.ID)
	getUser, err := libraries.FirestoreFind("users", login.ID)
	if err != nil {
		userPath := getUser.Ref.Path
		userNotExist := status.Errorf(codes.NotFound, "%q not found", userPath)
		if err.Error() == userNotExist.Error() {
			getEmail, err := libraries.FirestoreSearch("users", "Email", "==", login.ID)
			if err != nil {
				return http.StatusInternalServerError, err.Error()
			}
			if len(getEmail) > 0 {
				hashedPassword = getEmail[0].Data()["HashedPassword"].(string)
				username = getEmail[0].Data()["Username"].(string)
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
	_, err = libraries.FirestoreUpdateField("users", username, "IsActive", true)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusOK, username
}
