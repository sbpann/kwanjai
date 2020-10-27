package models

import (
	"context"
	"gin-sandbox/libraries"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// User model
type User struct {
	Username  string `json:"username"`
	Firstname string `json:"fisrtname"`
	Lastname  string `json:"lastname"`
	Password  string `json:"password"`
}

// CreateUser method for usermodel
func (user *User) CreateUser() (int, string) {
	ctx := context.Background()
	firebaseClient, err := libraries.FirebaseApp().Firestore(ctx)
	user.Username = strings.ToLower(user.Username)
	if err != nil {
		return http.StatusInternalServerError, "Cannot initialize Firestore."
	}
	getUser, err := firebaseClient.Collection("users").Doc(user.Username).Get(ctx)
	if err != nil {
		userPath := getUser.Ref.Path
		userNotExist := status.Errorf(codes.NotFound, "%q not found", userPath)
		if err.Error() == userNotExist.Error() {
			err = nil
		}
	}
	if getUser.Exists() {
		return http.StatusConflict, "User already exist."
	}
	_, err = firebaseClient.Collection("users").Doc(user.Username).Set(ctx, user)
	firebaseClient.Close()
	if err != nil {
		return http.StatusInternalServerError, "Cannot create user."
	}
	return http.StatusCreated, "User created successfully."
}

// Login user
func (user *User) Login() (int, string) {
	ctx := context.Background()
	firebaseClient, err := libraries.FirebaseApp().Firestore(ctx)
	user.Username = strings.ToLower(user.Username)
	if err != nil {
		return http.StatusInternalServerError, "Cannot initialize Firestore."
	}
	getUser, err := firebaseClient.Collection("users").Doc(user.Username).Get(ctx)
	firebaseClient.Close()
	if err != nil {
		userPath := getUser.Ref.Path
		userNotExist := status.Errorf(codes.NotFound, "%q not found", userPath)
		if err.Error() == userNotExist.Error() {
			return http.StatusBadRequest, "Cannot login with provided credential."
		}
		return http.StatusInternalServerError, "Cannot access Firestore."

	}
	hashedPassword := getUser.Data()["Password"]
	passwordPass := libraries.CheckPasswordHash(user.Password, hashedPassword.(string))
	if !passwordPass {
		return http.StatusBadRequest, "Cannot login with provided credential."
	}
	return http.StatusOK, "Logged in successfully."
}
