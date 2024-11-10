package svc

import (
	"fmt"
	"github.com/manjada/com/config"
	"github.com/manjada/com/db"
	"github.com/manjada/com/dto"
	repo2 "github.com/manjada/com/repo"
	"gopkg.in/gomail.v2"
	"regexp"
)

var dialer *gomail.Dialer

type EmailServiceInterface interface {
	SendEmail(data map[string]interface{}, to []string, templateKey string, from string, cc []string, attachment string) error
}

type EmailService struct {
	EmailRepo repo2.EmailTemplateRepoInterface
}

func NewEmailService(Db db.DBConnector) EmailServiceInterface {
	if dialer == nil {
		getConfig := config.GetConfig()
		dialer = gomail.NewDialer(
			getConfig.Smtp.Host, getConfig.Smtp.Port, getConfig.Smtp.User, getConfig.Smtp.Password)

	}
	return &EmailService{EmailRepo: repo2.NewEmailTemplateRepo(Db)}
}

func (e EmailService) SendEmail(data map[string]interface{}, to []string, templateKey string, from string, cc []string, attachment string) error {
	emailTemplate := e.EmailRepo.GetEmailTemplateByKey(templateKey)
	if emailTemplate == nil {
		return dto.ErrorParse(fmt.Errorf("email template not found"))
	}
	body := e.parseEmailBody(emailTemplate.Body, data)
	fmt.Println("Sending email to", to, "with body", body)

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	if cc != nil {
		m.SetHeader("Cc", cc...)
	}

	m.SetHeader("Subject", emailTemplate.Subject)
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
