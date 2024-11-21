package repo

import (
	"github.com/manjada/com/config"
	"github.com/manjada/com/db"
)

type EmailTemplate struct {
	TransactionModel
	Name        string `json:"name" gorm:"column:name"`
	Subject     string `json:"subject" gorm:"column:subject"`
	Body        string `json:"body" gorm:"column:body"`
	IsActive    bool   `json:"is_active" gorm:"column:is_active"`
	TemplateKey string `json:"template_key" gorm:"column:template_key"`
}

type EmailTemplateRepo struct {
	Db BaseRepoGorm
}

func (e EmailTemplateRepo) GetEmailTemplateByKey(templateKey string) *EmailTemplate {
	//TODO implement me
	var emailTemplate EmailTemplate
	if err := e.Db.Where(`template_key = ?`, templateKey).First(&emailTemplate).DbRepo.Error; err != nil {
		config.Error(err)
		return nil
	}
	return &emailTemplate
}

type EmailTemplateRepoInterface interface {
	GetEmailTemplateByKey(templateKey string) *EmailTemplate
}

func NewEmailTemplateRepo(Db db.DBConnector) EmailTemplateRepoInterface {
	return EmailTemplateRepo{Db: NewBaseRepo(Db)}
}
