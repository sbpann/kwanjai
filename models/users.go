package models

import (
	"kwanjai/libraries"
	"net/http"
	"strings"
	"time"
)

// User model.
type User struct {
	Username           string    `json:"username" binding:"required,ne=anonymous"`
	Email              string    `json:"email" binding:"required,email"`
	Firstname          string    `json:"firstname"`
	Lastname           string    `json:"lastname"`
	Password           string    `json:"password" binding:"required,min=8"`
	HashedPassword     string    `json:",omitempty"`
	IsSuperUser        bool      `json:"is_superuser"`
	IsVerified         bool      `json:"is_verified"`
	IsActive           bool      `json:"is_active"`
	JoinedDate         time.Time `json:"joined_date"`
	ProfilePicture     string    `json:"profile_picture"`
	Plan               string    `json:"plan"`
	Projects           int       `json:"projects"`
	CustomerID         string    `json:",omitempty"`
	SubscriptionID     string    `json:",omitempty"`
	DateOfSubscription int       `json:"date_of_subscription"`
}

// Register user method.
func (user *User) Register() (int, string, *User) {
	status, message, user := user.createUser()
	if status != http.StatusCreated || user == nil {
		return status, message, user
	}
	if user.Email == "test1@example.com" || user.Email == "test2@example.com" {
		return http.StatusOK, "Created account successfully.", user
	}
	status, message = user.SendVerificationEmail()
	return status, message, user
}

// Finduser user method.
func (user *User) Finduser() (int, string, *User) {
	status, message, user := user.findUser()
	return status, message, user
}

func (user *User) findUser() (int, string, *User) {
	if user.Username == "" && user.Email == "" {
		return http.StatusNotFound, "User not found.", nil
	}
	getUser, err := libraries.FirestoreFind("users", user.Username)
	if err != nil {
		getEmail, err := libraries.FirestoreSearch("users", "Email", "==", user.Email)
		if err != nil {
			return http.StatusInternalServerError, err.Error(), nil
		}
		if len(getEmail) > 0 {
			getEmail[0].DataTo(user)
			user.HashedPassword = ""
			return http.StatusOK, "Get user successfully.", user
		}
		return http.StatusNotFound, "User not found.", nil
	}
	err = getUser.DataTo(user)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	user.HashedPassword = ""
	return http.StatusOK, "Get user successfully.", user
}

func (user *User) createUser() (int, string, *User) {
	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)
	_, _, userFoud := user.findUser()
	if userFoud != nil {
		return http.StatusConflict, "Provided email or username is already registered.", nil
	}
	user.initialize()
	_, err := libraries.FirestoreCreateOrSet("users", user.Username, user)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	user.HashedPassword = ""
	return http.StatusCreated, "User created successfully.", user
}

// SendVerificationEmail method for user model.
func (user *User) SendVerificationEmail() (int, string) {
	email := new(VerificationEmail)
	email.Initialize(user.Username, user.Email)
	reference, _, err := libraries.FirestoreAdd("verificationEmail", email)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	email.ID = reference.ID
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
	user.Plan = "Starter"
	user.IsSuperUser = false
	user.IsVerified = false
	user.JoinedDate = time.Now().Truncate(time.Millisecond)
	user.ProfilePicture = "https://storage.googleapis.com/kwanjai-a3803.appspot.com/anonymous.png"
}

// MakeAnonymous user
func (user *User) MakeAnonymous() {
	user.Username = "anonymous"
	user.IsSuperUser = false
	user.IsVerified = false
	user.IsActive = false
	user.JoinedDate = time.Now()
}
