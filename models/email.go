package models

import (
	"fmt"
	"kwanjai/config"
	"kwanjai/libraries"
	"math/rand"
	"net/http"
	"net/smtp"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// VerificationEmail model.
type VerificationEmail struct {
	User        string
	Email       string `json:"email"`
	Key         string `json:"key"`
	UUID        string
	ExpiredDate time.Time
}

// Initialize email objects.
func (email *VerificationEmail) Initialize(user string, emailAddress string) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	email.User = user
	email.Email = emailAddress
	email.Key = fmt.Sprintf("%06d", random.Intn(999999))
	email.UUID = uuid.New().String()
	email.ExpiredDate = time.Now().Add(config.EmailVerficationLifetime)
}

// smtpServer data.
type smtpServer struct {
	host string
	port string
}

func (s *smtpServer) address() string {
	return s.host + ":" + s.port
}

// Send method for VerificationEmail object.
func (email *VerificationEmail) Send() (int, string) {
	// Sender data.
	from := "surus.d6101@gmail.com"
	password := config.EmailServicePassword
	to := []string{email.Email}
	smtpServer := smtpServer{host: "smtp.gmail.com", port: "587"}
	verificationLink := fmt.Sprintf("%v/verify_email/%v/", config.FrontendURL, email.UUID)
	message := fmt.Sprintf("To: %v\r\n"+
		"Subject: verification email.\r\n"+
		"\r\n"+
		"Hi %v.\r\n"+
		"Please verify your email using following link.\r\n"+
		"%v\r\n"+
		"Your verification code is: %v\r\n", to[0], email.User, verificationLink, email.Key)
	auth := smtp.PlainAuth("", from, password, smtpServer.host)
	err := smtp.SendMail(smtpServer.address(), auth, from, to, []byte(message))
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusOK, "Email sent."
}

// Verify method for VerificationEmail object.
// The method set user to be verified if verification is completed.
// If the email is expired, the method delete the email in database.
func (email *VerificationEmail) Verify() (int, string) {
	if email.UUID == "" {
		return http.StatusBadRequest, "Bad verification link."
	}
	getEmail, err := libraries.FirestoreFind("verificationemail", email.UUID)
	if !getEmail.Exists() {
		if err != nil {
			emailPath := getEmail.Ref.Path
			emailNotExist := status.Errorf(codes.NotFound, "%q not found", emailPath)
			if err.Error() == emailNotExist.Error() {
				return http.StatusBadRequest, "Bad verification link."
			}
			return http.StatusInternalServerError, err.Error()
		}
		return http.StatusBadRequest, "Bad verification link."
	}
	verificationEmail := new(VerificationEmail)
	err = getEmail.DataTo(verificationEmail)
	now := time.Now()
	expriredDate := verificationEmail.ExpiredDate
	exprired := now.After(expriredDate)
	if exprired {
		_, err = libraries.FirestoreDelete("verificationemail", email.UUID)
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}
		return http.StatusBadRequest, "Link is expired."
	}
	if email.Key == verificationEmail.Key {
		_, err = libraries.FirestoreUpdateField("users", verificationEmail.User, "IsVerified", true)
		_, err = libraries.FirestoreDelete("verificationemail", email.UUID)
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}
		return http.StatusOK, "Email verified."
	}
	return http.StatusBadRequest, "Key is invalid."
}
