package models

import (
	"fmt"
	"gin-sandbox/config"
	"math/rand"
	"net/smtp"
	"time"

	"github.com/google/uuid"
)

// VerificationEmail model
type VerificationEmail struct {
	User        string
	Email       string
	Key         string
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
	email.ExpiredDate = time.Now().AddDate(0, 0, 1).Format(time.RFC3339)
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
func (email *VerificationEmail) Send() error {
	// Sender data.
	from := "surus.d6101@gmail.com"
	password := "secret"
	to := []string{email.Email}
	smtpServer := smtpServer{host: "smtp.gmail.com", port: "587"}
	verificationLink := fmt.Sprintf("%v/verify/%v/", config.FrontendURL, email.UUID)
	message := []byte(fmt.Sprintf("To: %v\r\n"+
		"Subject: verification email.\r\n"+
		"\r\n"+
		"Hi %v.\r\n"+
		"Please verify your email using following link.\r\n"+
		"%v\r\n"+
		"Your verification code is: %v\r\n", email.User, to[0], verificationLink, email.Key))
	auth := smtp.PlainAuth("", from, password, smtpServer.host)
	err := smtp.SendMail(smtpServer.Address(), auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Email Sent!")
	return nil
}
