package util

import (
	"crypto/tls"
	"net/mail"
	"net/smtp"
	"os"
)

func SendMail(toAddr string, subjectText string, body string) error {
	from, err := mail.ParseAddress(os.Getenv("SMTP_FROM"))

	if err != nil {
		return err
	}

	to := mail.Address{Address: toAddr, Name: ""}

	fromLine := "From: " + from.String() + "\n"
	toLine := "To: " + to.String() + "\n"
	subject := "Subject: " + subjectText + "!\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(fromLine + toLine + subject + mime + body)

	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	auth := smtp.PlainAuth("", os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASS"), host)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	client, err := smtp.Dial(host + ":" + port)
	if err != nil {
		return err
	}

	client.StartTLS(tlsconfig)

	if err = client.Auth(auth); err != nil {
		return err
	}
	if err = client.Mail(from.Address); err != nil {
		return err
	}

	if err = client.Rcpt(to.Address); err != nil {
		return err
	}

	w, err := client.Data()

	if err != nil {
		return err
	}

	_, err = w.Write(msg)

	if err != nil {
		return err
	}

	err = w.Close()

	if err != nil {
		return err
	}

	err = client.Quit()

	if err != nil {
		return err
	}

	return nil

}
