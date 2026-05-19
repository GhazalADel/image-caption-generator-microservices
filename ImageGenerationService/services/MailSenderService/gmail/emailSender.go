package gmail

import (
	"ImageGenerationService/config"
	"fmt"
	"net/smtp"
)

func SendCustomEmail(subject string, url string, receiver string) error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}
	auth := smtp.PlainAuth(
		"",
		cfg.Username,
		cfg.Password,
		cfg.Host,
	)

	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\nURL: %s", cfg.Username, receiver, subject, url)

	err = smtp.SendMail(
		cfg.Address,
		auth,
		cfg.Username,
		[]string{receiver},
		[]byte(msg),
	)
	if err != nil {
		return err
	}

	return err
}
