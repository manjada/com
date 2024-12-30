package repo

import (
	"github.com/manjada/com/config"
	"github.com/manjada/com/db"
	"github.com/manjada/com/db/repo"
)

type EmailTemplate struct {
	TransactionModel
	Name        string `gorm:"column:name;varchar(255)"`
	Subject     string `gorm:"column:subject;varchar(255)"`
	Body        string `gorm:"column:body"`
	IsActive    bool   `gorm:"column:is_active"`
	TemplateKey string `gorm:"column:template_key;varchar(255)"`
	ClientId    string `gorm:"column:client_id;varchar(255)"`
}

type EmailTemplateRepo struct {
	Db repo.BaseRepoGorm
}

func (e EmailTemplateRepo) GetEmailTemplateByKeyAndClientId(templateKey string, clientId string) *EmailTemplate {
	//TODO implement me
	var emailTemplate EmailTemplate
	if err := e.Db.Where(`template_key = ? and client_id = ?`, templateKey, clientId).First(&emailTemplate).DbRepo.Error; err != nil {
		config.Error(err)
		return nil
	}
	return &emailTemplate
}

type EmailTemplateRepoInterface interface {
	GetEmailTemplateByKeyAndClientId(templateKey string, clientId string) *EmailTemplate
}

func NewEmailTemplateRepo(Db db.DBConnector) EmailTemplateRepoInterface {
	return EmailTemplateRepo{Db: repo.NewBaseRepo(Db)}
}
