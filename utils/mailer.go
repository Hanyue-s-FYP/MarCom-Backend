package utils

import (
	"fmt"
	"log/slog"
	"net/smtp"
)

func SendMail(to, subject, body string) error {
	config := GetConfig()
	from := config.SMTP_EMAIL
	pass := config.SMTP_PW
	host := config.SMTP_ADDR
	port := config.SMTP_PORT

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail(fmt.Sprintf("%s:%s", host, port),
		smtp.PlainAuth("", from, pass, host),
		from, []string{to}, []byte(msg))

	if err != nil {
		slog.Error(fmt.Sprintf("smtp error: %s", err))
		return err
	}
	slog.Info("Successfully sended to " + to)
	return nil
}
