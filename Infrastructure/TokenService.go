package infrastructure

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"

	domain "github.com/gedyzed/blog-starter-project/Domain"
)

type TokenSevice struct {}


func NewTokenServices() domain.ITokenService {
	return &TokenSevice{}
}


func (ts *TokenSevice) SendEmail(to []string, subject string, body string) error {


	appPassword := os.Getenv("Mail_APP_PASSWORD")
	address := "smtp.gmail.com"
	auth := smtp.PlainAuth(
		"",
		"gediizeyuu@gmail.com",
		appPassword,
		address,
	)

	senderEmail := "gediizeyuu@gmail.com"
	message := "Subject: " + subject + "\r\n\r\n" + body

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		senderEmail,
		to,
		[]byte(message),
	)

	if err != nil {
		fmt.Println("email",  err.Error())
		return errors.New("cannot send email.please, check your email")
	}

	return nil
}
