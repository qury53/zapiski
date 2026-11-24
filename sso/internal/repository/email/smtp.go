package email_smtp_repository

import (
	"net/smtp"
)

type SmtpRepository struct {
	auth     smtp.Auth
	host     string
	addr     string
	fromText string
}

func NewSmtpRepository(config *Config) *SmtpRepository {
	var auth smtp.Auth
	if config.Username != "" || config.Password != "" {
		auth = smtp.PlainAuth("", config.Username, config.Password, config.Host)
	}

	return &SmtpRepository{
		auth:     auth,
		host:     config.Host,
		addr:     config.Host + ":" + config.Port,
		fromText: config.SenderEmail,
	}
}
