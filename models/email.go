package models

import (
	"fmt"
	"kwanjai/config"
	"kwanjai/libraries"
	"math/rand"
	"net/http"
	"net/smtp"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// VerificationEmail model
type VerificationEmail struct {
	User        string
	Email       string
	Key         string `json:"key"`
	UUID        string
	ExpiredDate string
}

// Initialize email
func (email *VerificationEmail) Initialize(user string, emailAddress string) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	email.User = user
	email.Email = emailAddress
	email.Key = fmt.Sprintf("%06d", random.Intn(999999))
	email.UUID = uuid.New().String()
	email.ExpiredDate = time.Now().Add(config.EmailVerficationLifetime).Format(time.RFC3339)
}

// smtpServer data to smtp server
type smtpServer struct {
	host string
	port string
}

// Address URI to smtp server
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

// Send email
func (email *VerificationEmail) Send() (int, string) {
	// Sender data.
	from := "surus.d6101@gmail.com"
	password := config.EmailServicePassword
	to := []string{email.Email}
	smtpServer := smtpServer{host: "smtp.gmail.com", port: "587"}
	verificationLink := fmt.Sprintf("%v/verify/%v/", config.FrontendURL, email.UUID)
	message := fmt.Sprintf("To: %v\r\n"+
		"Subject: verification email.\r\n"+
		"\r\n"+
		"Hi %v.\r\n"+
		"Please verify your email using following link.\r\n"+
		"%v\r\n"+
		"Your verification code is: %v\r\n", to[0], email.User, verificationLink, email.Key)
	auth := smtp.PlainAuth("", from, password, smtpServer.host)
	err := smtp.SendMail(smtpServer.Address(), auth, from, to, []byte(message))
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	return http.StatusOK, "Email sent."
}

// Verify email
func (email *VerificationEmail) Verify() (int, string) {
	firestoreClient, err := libraries.FirebaseApp().Firestore(config.Context)
	defer firestoreClient.Close()
	getEmail, err := firestoreClient.Collection("verificationemail").Doc(email.UUID).Get(config.Context)
	if !getEmail.Exists() {
		if err != nil {
			emailPath := getEmail.Ref.Path
			emailNotExist := status.Errorf(codes.NotFound, "%q not found", emailPath)
			if err.Error() == emailNotExist.Error() {
				return http.StatusBadRequest, "bad verification link."
			}
			return http.StatusInternalServerError, err.Error()
		}
		return http.StatusBadRequest, "bad verification link."
	}
	verificationEmail := new(VerificationEmail)
	err = getEmail.DataTo(&verificationEmail)
	now := time.Now()
	expriredDate, err := time.Parse(time.RFC3339, verificationEmail.ExpiredDate)
	exprired := now.After(expriredDate)
	if exprired {
		_, err = firestoreClient.Collection("verificationemail").Doc(email.UUID).Delete(config.Context)
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}
		return http.StatusBadRequest, "Link is expired."
	}
	if email.Key == verificationEmail.Key {
		_, err = firestoreClient.Collection("users").Doc(verificationEmail.User).Update(config.Context, []firestore.Update{
			{
				Path:  "IsVerified",
				Value: true,
			},
		})
		_, err = firestoreClient.Collection("verificationemail").Doc(email.UUID).Delete(config.Context)
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}
		return http.StatusOK, "Email verified."
	}
	return http.StatusBadRequest, "Key is invalid."
}
