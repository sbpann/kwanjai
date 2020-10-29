package models

import (
	"context"
	"kwanjai/libraries"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// User model
type User struct {
	Username       string `form:"username" json:"username" binding:"required,ne=anonymous"`
	Email          string `form:"email" json:"email" binding:"required,email"`
	Firstname      string `form:"firstname" json:"firstname"`
	Lastname       string `form:"lastname" json:"lastname"`
	Password       string `form:"password" json:"password" binding:"required,min=8"`
	HashedPassword string `json:",omitempty"`
	IsSuperUser    bool   `json:"is_superuser"`
	IsVerified     bool   `json:"is_verified"`
	IsActive       bool   `json:"is_active"`
	JoinedDate     string `json:"joined_date"`
}

// LoginCredential info
type LoginCredential struct {
	ID       string `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type userPerform interface {
	createUser() (int, string, *User)
	findUser(username string) (int, string, *User)
}

type authenticationPerform interface {
	login() (int, string)
}

// Login user
func Login(perform authenticationPerform) (int, string) {
	return perform.login()
}

// Register user
func Register(perform userPerform) (int, string, *User) {
	status, message, user := perform.createUser()
	if status != http.StatusCreated || user == nil {
		return status, message, user
	}
	// status, message = perform.sendVerificationEmail()
	return status, message, user
}

func (user *User) createUser() (int, string, *User) {
	ctx := context.Background()
	firestoreClient, err := libraries.FirebaseApp().Firestore(ctx)
	defer firestoreClient.Close()
	user.Username = strings.ToLower(user.Username)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	getUser, err := firestoreClient.Collection("users").Doc(user.Username).Get(ctx)
	if err != nil {
		userPath := getUser.Ref.Path
		userNotExist := status.Errorf(codes.NotFound, "%q not found", userPath)
		if err.Error() == userNotExist.Error() {
			err = nil
		}
	}
	if getUser.Exists() {
		return http.StatusConflict, "User already exist.", nil
	}
	findEmail := firestoreClient.Collection("users").Where("Email", "==", user.Email).Documents(ctx)
	foundEmail, err := findEmail.GetAll()
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	if len(foundEmail) > 0 {
		return http.StatusConflict, "There is registered user with this email.", nil
	}
	user.initialize()
	_, err = firestoreClient.Collection("users").Doc(user.Username).Set(ctx, user)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	user.HashedPassword = ""
	return http.StatusCreated, "User created successfully.", user
}

func (user *User) login() (int, string) {
	login := new(LoginCredential)
	login.ID = user.Username
	login.Password = user.Password
	return login.login()
}

func (user *User) findUser(username string) (int, string, *User) {
	if username == "" {
		return http.StatusNotFound, "User not found.", nil
	}
	ctx := context.Background()
	firestoreClient, err := libraries.FirebaseApp().Firestore(ctx)
	defer firestoreClient.Close()
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	getUser, err := firestoreClient.Collection("users").Doc(username).Get(ctx)
	if !getUser.Exists() {
		return http.StatusNotFound, "User not found.", nil
	}
	err = getUser.DataTo(&user)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	user.HashedPassword = ""
	return http.StatusOK, "Get user successfully.", user
}

func (user *User) sendVerificationEmail() (int, string) {
	email := new(VerificationEmail)
	email.Initialize(user.Username, user.Email)
	ctx := context.Background()
	firestoreClient, err := libraries.FirebaseApp().Firestore(ctx)
	defer firestoreClient.Close()
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	_, err = firestoreClient.Collection("verificationemail").Doc(email.User).Set(ctx, email)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	status, message := email.Send()
	return status, message
}

func (login *LoginCredential) login() (int, string) {
	hashedPassword := ""
	username := ""
	login.ID = strings.ToLower(login.ID)
	ctx := context.Background()
	firestoreClient, err := libraries.FirebaseApp().Firestore(ctx)
	defer firestoreClient.Close()
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	getUser, err := firestoreClient.Collection("users").Doc(login.ID).Get(ctx)
	if err != nil {
		userPath := getUser.Ref.Path
		userNotExist := status.Errorf(codes.NotFound, "%q not found", userPath)
		if err.Error() == userNotExist.Error() {
			findEmail := firestoreClient.Collection("users").Where("Email", "==", login.ID).Documents(ctx)
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
	_, err = firestoreClient.Collection("users").Doc(username).Update(ctx, []firestore.Update{{Path: "IsActive", Value: true}})
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusOK, username
}

// HashPassword before register
func (user *User) HashPassword() {
	hashedpassword, _ := libraries.HashPassword(user.Password)
	user.Password = "password_is_created"
	user.HashedPassword = hashedpassword
}

func (user *User) initialize() {
	user.IsSuperUser = false
	user.IsVerified = false
	user.JoinedDate = time.Now().Format(time.RFC3339)
}
