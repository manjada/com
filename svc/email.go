package svc

import (
	"fmt"
	"github.com/manjada/com/config"
	"gopkg.in/gomail.v2"
	"regexp"
)

var dialer *gomail.Dialer

type EmailServiceInterface interface {
	SendEmail(data map[string]interface{}, to []string, templateKey string, from string, cc []string, attachment string, clientId string) error
}

type EmailService struct {
}

func NewEmailService() *EmailService {
	if dialer == nil {
		getConfig := config.GetConfig()
		dialer = gomail.NewDialer(
			getConfig.Smtp.Host, getConfig.Smtp.Port, getConfig.Smtp.User, getConfig.Smtp.Password)

	}
	return &EmailService{}
}

func (e EmailService) SendEmail(data map[string]interface{}, to []string, from string, cc []string, attachment string, body string, subject string) error {

	body = e.parseEmailBody(body, data)
	fmt.Println("Sending email to", to, "with body", body)

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	if cc != nil {
		m.SetHeader("Cc", cc...)
	}

	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	if attachment != "" {
		m.Attach(attachment)
	}

	if err := dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func (e EmailService) parseEmailBody(body string, data map[string]interface{}) string {
	re := regexp.MustCompile(`\${{(\w+)}}`)
	return re.ReplaceAllStringFunc(body, func(match string) string {
		key := re.FindStringSubmatch(match)[1]
		if val, ok := data[key]; ok {
			return val.(string)
		}
		return match
	})
}
