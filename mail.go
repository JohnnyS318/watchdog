package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

const htmlBody = `<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
		<title>Watchdog: Monitor update</title>
	</head>
	<body>
		<p>Hello %s</p>
		<h2>The website (%s) your are watching has changed!</h2>
		<p>Best regards, your Watchdog</p>
	</body>
</html>`

func CreateHTML() string {
	return fmt.Sprintf(htmlBody, Config.Receiver, Config.Watchtarget)
}

func SendMail(client *mail.SMTPClient) error {
	email := mail.NewMSG()
	email.SetFrom("From  Jonas Schneider <jonas.max.schneider@gmail.com>").AddTo(Config.Receiver...).SetSubject("Monitor update").SetBody(mail.TextHTML, CreateHTML())

	if email.Error != nil {
		return email.Error
	}

	return email.Send(client)
}

func RegisterMailClient() (*mail.SMTPClient, error) {
	server := mail.NewSMTPClient()

	server.Host = os.Getenv(EMAIL_HOST)
	port, err := strconv.Atoi(os.Getenv(EMAIL_PORT))
	if err != nil {
		return nil, err
	}
	server.Port = port
	server.Username = os.Getenv(EMAIL_USERNAME)
	server.Password = os.Getenv(EMAIL_PASSWORD)

	encryption, err := strconv.Atoi(os.Getenv(EMAIL_ENCRYPTION))
	if err != nil {
		return nil, err
	}
	server.Encryption = mail.Encryption(encryption)

	server.KeepAlive = true
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		return nil, err
	}

	return client, nil
}
