package models

import (
	"context"
	"fmt"
	"kwanjai/libraries"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// User model
type User struct {
	Username       string `form:"username" json:"username" binding:"required"`
	Email          string `form:"email" json:"email" binding:"required,email"`
	Firstname      string `form:"firstname" json:"firstname"`
	Lastname       string `form:"lastname" json:"lastname"`
	Password       string `form:"password" ,json:"password" binding:"required,min=8"`
	HashedPassword string
	IsSuperUser    bool `json:"is_superuser"`
	IsVerified     bool `json:"is_verified"`
	JoinedDate     string
}

// LoginCredential info
type LoginCredential struct {
	ID       string `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginOption struct {
	fromRegister bool
}

type userPerform interface {
	createUser() (int, string)
	login(loginOption *loginOption) (int, string)
	sendVerificationEmail() (int, string)
}

type authenticationPerform interface {
	login(loginOption *loginOption) (int, string)
}

// Login user
func Login(perform authenticationPerform) (int, string) {
	option := loginOption{fromRegister: false}
	return perform.login(&option)
}

// Register user
func Register(perform userPerform) (int, string) {
	status, detail := perform.createUser()
	if status != http.StatusCreated {
		return status, detail
	}
	// status, detail = perform.sendVerificationEmail()
	// if status != http.StatusOK {
	// 	return status, detail
	// }
	option := loginOption{fromRegister: true}
	return perform.login(&option)
}

func (user *User) createUser() (int, string) {
	ctx := context.Background()
	firestoreClient, err := libraries.FirebaseApp().Firestore(ctx)
	defer firestoreClient.Close()
	user.Username = strings.ToLower(user.Username)
	if err != nil {
		return http.StatusInternalServerError, "Cannot initialize Firestore."
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
		return http.StatusConflict, "User already exist."
	}
	findEmail := firestoreClient.Collection("users").Where("Email", "==", user.Email).Documents(ctx)
	foundEmail, err := findEmail.GetAll()
	if err != nil {
		return http.StatusInternalServerError, "Cannot access Firestore."
	}
	if len(foundEmail) > 0 {
		return http.StatusConflict, "There is registered user with this email."
	}
	user.initialize()
	_, err = firestoreClient.Collection("users").Doc(user.Username).Set(ctx, user)
	if err != nil {
		return http.StatusInternalServerError, "Cannot create user."
	}
	return http.StatusCreated, "User created successfully."
}

func (user *User) login(loginOption *loginOption) (int, string) {
	login := new(LoginCredential)
	login.ID = user.Username
	login.Password = user.Password
	return login.login(loginOption)
}

func (user *User) sendVerificationEmail() (int, string) {
	email := new(VerificationEmail)
	email.Initialize(user.Username, user.Email)
	ctx := context.Background()
	firestoreClient, err := libraries.FirebaseApp().Firestore(ctx)
	defer firestoreClient.Close()
	if err != nil {
		return http.StatusInternalServerError, "Cannot access Firestore."
	}
	_, err = firestoreClient.Collection("verificationemail").Doc(email.User).Set(ctx, email)
	if err != nil {
		return http.StatusInternalServerError, "Cannot access Firestore."
	}
	email.Send()
	return http.StatusOK, "Email sent."
}

func (login *LoginCredential) login(loginOption *loginOption) (int, string) {
	hashedPassword := ""
	if loginOption.fromRegister {
		return http.StatusCreated, "User created successfully."
	}
	login.ID = strings.ToLower(login.ID)
	ctx := context.Background()
	firestoreClient, err := libraries.FirebaseApp().Firestore(ctx)
	defer firestoreClient.Close()
	if err != nil {
		return http.StatusInternalServerError, "Cannot initialize Firestore."
	}
	getUser, err := firestoreClient.Collection("users").Doc(login.ID).Get(ctx)
	if err != nil {
		userPath := getUser.Ref.Path
		userNotExist := status.Errorf(codes.NotFound, "%q not found", userPath)
		if err.Error() == userNotExist.Error() {
			findEmail := firestoreClient.Collection("users").Where("Email", "==", login.ID).Documents(ctx)
			foundEmail, err := findEmail.GetAll()
			if err != nil {
				return http.StatusInternalServerError, "Cannot access Firestore."
			}
			if len(foundEmail) > 0 {
				hashedPassword = foundEmail[0].Data()["HashedPassword"].(string)
			} else {
				return http.StatusBadRequest, "Cannot login with provided credential."
			}
		} else {
			return http.StatusInternalServerError, "Cannot access Firestore."
		}
	} else {
		hashedPassword = getUser.Data()["HashedPassword"].(string)
	}
	fmt.Println(hashedPassword)
	passwordPass := libraries.CheckPasswordHash(login.Password, hashedPassword)
	if !passwordPass && !loginOption.fromRegister {
		return http.StatusBadRequest, "Cannot login with provided credential."
	}
	return http.StatusOK, "Logged in successfully."
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
