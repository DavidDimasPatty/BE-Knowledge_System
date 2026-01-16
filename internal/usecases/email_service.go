package usecases

import (
	"be-knowledge/configs"

	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailService interface {
	SendEmail(to, subject, body string) error
}

type emailService struct {
	config *configs.Config
}

func NewEmailService(cfg *configs.Config) EmailService {
	return &emailService{config: cfg}
}

func (e *emailService) SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", e.config.SMTPFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	port, err := strconv.Atoi(e.config.SMTPPort)
	if err != nil {
		return err
	}

	d := gomail.NewDialer(
		e.config.SMTPHost,
		port,
		e.config.SMTPUser,
		e.config.SMTPPass,
	)

	return d.DialAndSend(m)
}
