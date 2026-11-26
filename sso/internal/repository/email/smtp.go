package email_smtp_repository

import (
	"net/smtp"
)

type SmtpRepository struct {
	auth        smtp.Auth
	host        string
	addr        string
	fromText    string
	implicitTLS bool
}

func NewSmtpRepository(config *Config) *SmtpRepository {
	return &SmtpRepository{
		auth:        smtp.PlainAuth("", config.Username, config.Password, config.Host),
		host:        config.Host,
		addr:        config.Host + ":" + config.Port,
		fromText:    config.SenderEmail,
		implicitTLS: config.Port == "465",
	}
}
