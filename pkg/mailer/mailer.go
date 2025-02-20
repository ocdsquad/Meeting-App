package mailer

import (
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"gopkg.in/gomail.v2"
)

type Mailer struct {
	dialer *gomail.Dialer
	sender string
	Logger echo.Logger
}

func NewMailer(logger echo.Logger) *Mailer {
	mailPort, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		logger.Fatal(err)
	}
	mailHost := os.Getenv("MAIL_HOST")
	mailUsername := os.Getenv("MAIL_USERNAME")
	mailPassword := os.Getenv("MAIL_PASSWORD")
	mailSender := os.Getenv("MAIL_SENDER")

	dialer := gomail.NewDialer(mailHost, mailPort, mailUsername, mailPassword)
	return &Mailer{
		dialer: dialer,
		sender: mailSender,
	}
}

func (m *Mailer) SendMail(to, subject, body string) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", m.sender)
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", subject)
	mail.SetBody("text/plain", body)

	if err := m.dialer.DialAndSend(mail); err != nil {
		return err
	}

	return nil
}
