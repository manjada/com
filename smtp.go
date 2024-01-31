package mjd

import (
	"gopkg.in/gomail.v2"
)

var dialer *gomail.Dialer

type Mail struct {
	From       string
	To         []string
	Cc         []string
	Subject    string
	Body       string
	Attachment string
}

type MailInterface interface {
	SendMail() error
}

func NewSmtpDialer() *Mail {
	if dialer == nil {
		config := GetConfig()
		dialer = gomail.NewDialer(
			config.Smtp.Host, config.Smtp.Port, config.Smtp.User, config.Smtp.Password)
	}
	return &Mail{}
}

func (data Mail) SendMail() error {

	m := gomail.NewMessage()
	m.SetHeader("From", data.From)
	m.SetHeader("To", data.To...)
	if data.Cc != nil {
		m.SetHeader("Cc", data.Cc...)
	}

	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", data.Body)
	if data.Attachment != "" {
		m.Attach(data.Attachment)
	}

	if err := dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
