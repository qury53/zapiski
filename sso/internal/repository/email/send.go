package email_smtp_repository

import (
	"context"
	"net/smtp"
)

func (r *SmtpRepository) Send(
	ctx context.Context,
	toEmail string,
	code string,
) error {
	msg := []byte("From: " + r.fromText + "\r\n" +
		"To: " + toEmail + "\r\n" +
		"Subject: Подтверждение Email\r\n" +
		"\r\n" +
		"Ваш токен: " + code)

	return smtp.SendMail(r.addr, r.auth, r.fromText, []string{toEmail}, msg)
}
