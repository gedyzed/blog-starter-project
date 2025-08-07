package infrastructure

import (
	"errors"
	"fmt"
	"net/smtp"

	"github.com/gedyzed/blog-starter-project/Infrastructure/config"
	usecases "github.com/gedyzed/blog-starter-project/Usecases"
)

type TokenService struct {
	emailConfig config.EmailConfig
	appUrl string
}


func NewTokenService(emailCfg config.EmailConfig, appUrl string) *TokenService {
	return &TokenService{
		emailConfig: emailCfg,
		appUrl: appUrl,
	}
}



func (ts *TokenService) SendEmail(to []string, subject string, body string) error {

	auth := smtp.PlainAuth(
		"",
		ts.emailConfig.SenderEmail,
		ts.emailConfig.AppPassword,
		ts.emailConfig.SMTPHost,
	)

	if subject == usecases.ResetPasswordEmailSubject {
		body = usecases.ResetPasswordEmailBodyText + ts.appUrl + body
	}
    
	message := "Subject: " + subject + "\r\n\r\n" + body

	err := smtp.SendMail(
		ts.emailConfig.SMTPHost + ":" + ts.emailConfig.SMTPPort,
		auth,
		ts.emailConfig.SenderEmail,
		to,
		[]byte(message),
	)

	if err != nil {
		fmt.Println("email",  err.Error())
		return errors.New("cannot send email.please, check your email")
	}

	return nil
}






