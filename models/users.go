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

type loginOption struct {
	fromRegister bool
}

// UserPerform action on model
type userPerform interface {
	createUser() (int, string)
	login(loginOption *loginOption) (int, string)
}

// Login user
func Login(perform userPerform) (int, string) {
	option := loginOption{fromRegister: false}
	return perform.login(&option)
}

// Register user
func Register(perform userPerform) (int, string) {
	status, detail := perform.createUser()
	if status != http.StatusCreated {
		return status, detail
	}
	option := loginOption{fromRegister: true}
	return perform.login(&option)
}

func (user *User) createUser() (int, string) {
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

func (user *User) login(loginOption *loginOption) (int, string) {
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
	if loginOption.fromRegister {
		return http.StatusCreated, "User created successfully."
	}
	return http.StatusOK, "Logged in successfully."
}
