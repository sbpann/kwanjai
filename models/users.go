package models

import (
	"kwanjai/config"
	"kwanjai/libraries"
	"net/http"
	"strings"
	"time"
)

// User model.
type User struct {
	Username       string `json:"username" binding:"required,ne=anonymous"`
	Email          string `json:"email" binding:"required,email"`
	Firstname      string `json:"firstname"`
	Lastname       string `json:"lastname"`
	Password       string `json:"password" binding:"required,min=8"`
	HashedPassword string `json:",omitempty"`
	IsSuperUser    bool   `json:"is_superuser"`
	IsVerified     bool   `json:"is_verified"`
	IsActive       bool   `json:"is_active"`
	JoinedDate     string `json:"joined_date"`
}

type userPerform interface {
	createUser() (int, string, *User)
	findUser() (int, string, *User)
}

// Register user method for interface with controller.
func Register(perform userPerform) (int, string, *User) {
	status, message, user := perform.createUser()
	if status != http.StatusCreated || user == nil {
		return status, message, user
	}
	if user.Email == "test@example.com" {
		return http.StatusOK, "Created account successfully.", user
	}
	status, message = user.SendVerificationEmail()
	return status, message, user
}

// Finduser user method for interface with controller.
func Finduser(perform userPerform) (int, string, *User) {
	status, message, user := perform.findUser()
	return status, message, user
}

func (user *User) createUser() (int, string, *User) {
	firestoreClient, err := libraries.FirebaseApp().Firestore(config.Context)
	defer firestoreClient.Close()
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)
	_, _, userFoud := user.findUser()
	if userFoud != nil {
		return http.StatusConflict, "Provided email or username is already registered.", nil
	}
	user.initialize()
	_, err = firestoreClient.Collection("users").Doc(user.Username).Set(config.Context, user)
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

func (user *User) findUser() (int, string, *User) {
	if user.Username == "" && user.Email == "" {
		return http.StatusNotFound, "User not found.", nil
	}
	firestoreClient, err := libraries.FirebaseApp().Firestore(config.Context)
	defer firestoreClient.Close()
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	getUser, err := firestoreClient.Collection("users").Doc(user.Username).Get(config.Context)
	if err != nil {
		findEmail := firestoreClient.Collection("users").Where("Email", "==", user.Email).Documents(config.Context)
		foundEmail, err := findEmail.GetAll()
		if err != nil {
			return http.StatusInternalServerError, err.Error(), nil
		}
		if len(foundEmail) > 0 {
			foundEmail[0].DataTo(&user)
			user.HashedPassword = ""
			return http.StatusOK, "Get user successfully.", user
		}
		return http.StatusNotFound, "User not found.", nil
	}
	err = getUser.DataTo(&user)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	user.HashedPassword = ""
	return http.StatusOK, "Get user successfully.", user
}

// SendVerificationEmail method for user model.
func (user *User) SendVerificationEmail() (int, string) {
	email := new(VerificationEmail)
	email.Initialize(user.Username, user.Email)
	firestoreClient, err := libraries.FirebaseApp().Firestore(config.Context)
	defer firestoreClient.Close()
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	_, err = firestoreClient.Collection("verificationemail").Doc(email.UUID).Set(config.Context, email)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	status, message := email.Send()
	return status, message
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

// MakeAnonymous user
func (user *User) MakeAnonymous() {
	user.Username = "anonymous"
	user.IsSuperUser = false
	user.IsVerified = false
	user.JoinedDate = time.Now().Format(time.RFC3339)
}
