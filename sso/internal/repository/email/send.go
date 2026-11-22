package email_smtp_repository

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"mime"
	"net"
	"net/smtp"
	"time"
)

const (
	smtpDialTimeout = 10 * time.Second
	smtpIOTimeout   = 20 * time.Second
)

func (r *SmtpRepository) Send(
	ctx context.Context,
	toEmail string,
	code string,
) error {
	dialer := &net.Dialer{Timeout: smtpDialTimeout}

	var (
		conn net.Conn
		err  error
	)

	// Порт 465 — implicit TLS: сервер не присылает приветствие в открытом
	// виде, поэтому обычный net/smtp там зависает.
	if r.implicitTLS {
		conn, err = tls.DialWithDialer(dialer, "tcp", r.addr, &tls.Config{
			ServerName: r.host,
			MinVersion: tls.VersionTLS12,
		})
	} else {
		conn, err = dialer.DialContext(ctx, "tcp", r.addr)
	}
	if err != nil {
		return fmt.Errorf("smtp dial %s: %w", r.addr, err)
	}

	defer func() { _ = conn.Close() }()

	// Жёсткий дедлайн на весь диалог: без него любой сбой сети вешает RPC.
	_ = conn.SetDeadline(time.Now().Add(smtpIOTimeout))

	// Отмена контекста рвёт соединение, а не ждёт TCP-таймаут.
	stop := context.AfterFunc(ctx, func() { _ = conn.SetDeadline(time.Now()) })
	defer stop()

	client, err := smtp.NewClient(conn, r.host)
	if err != nil {
		return fmt.Errorf("smtp new client: %w", err)
	}

	defer func() { _ = client.Close() }()

	if !r.implicitTLS {
		if ok, _ := client.Extension("STARTTLS"); ok {
			if err := client.StartTLS(&tls.Config{
				ServerName: r.host,
				MinVersion: tls.VersionTLS12,
			}); err != nil {
				return fmt.Errorf("smtp starttls: %w", err)
			}
		}
	}

	if r.auth != nil {
		if err := client.Auth(r.auth); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
	}

	if err := client.Mail(r.fromText); err != nil {
		return fmt.Errorf("smtp mail from: %w", err)
	}

	if err := client.Rcpt(toEmail); err != nil {
		return fmt.Errorf("smtp rcpt to: %w", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}

	if _, err := writer.Write(buildMessage(r.fromText, toEmail, code)); err != nil {
		_ = writer.Close()

		return fmt.Errorf("smtp write body: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("smtp close body: %w", err)
	}

	return client.Quit()
}

// buildMessage собирает письмо с корректными MIME-заголовками: сырой UTF-8
// в Subject и теле многие SMTP-серверы отвергают.
func buildMessage(from, to, code string) []byte {
	subject := mime.QEncoding.Encode("utf-8", "Подтверждение Email")
	body := base64.StdEncoding.EncodeToString([]byte("Ваш токен: " + code))

	return []byte(
		"From: " + from + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"Date: " + time.Now().Format(time.RFC1123Z) + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n" +
			"Content-Transfer-Encoding: base64\r\n" +
			"\r\n" +
			body + "\r\n")
}
