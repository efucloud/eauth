package config

import (
	"gopkg.in/gomail.v2"
)

func NewMailClient() (client *gomail.Dialer) {
	if ApplicationConfig.Email == nil {
		return
	}
	email := ApplicationConfig.Email
	if len(email.Username) > 0 && len(email.Password) > 0 {
		client = gomail.NewDialer(email.SmtpServer, email.SmtpPort, email.Username, email.Password)
	} else if len(email.Username) > 0 {
		client = &gomail.Dialer{
			Host:     email.SmtpServer,
			Port:     email.SmtpPort,
			Username: email.Username,
		}
	}
	return
}
