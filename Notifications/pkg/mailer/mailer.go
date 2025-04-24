package mailer

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type Mailer interface {
	SendEmail(to, subject, body string) error
}

type Config struct {
	User string
	Pass string
	Port int
	Host string
}

type mailer struct {
	conf Config
}

func New(conf Config) Mailer {
	return &mailer{
		conf: conf,
	}
}

func (m *mailer) SendEmail(to, subject, body string) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", m.conf.User)
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", subject)
	mail.SetBody("text/plain", body)

	d := gomail.NewDialer(m.conf.Host, m.conf.Port, m.conf.User, m.conf.Pass)
	if err := d.DialAndSend(mail); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
