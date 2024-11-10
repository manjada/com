package svc

import (
	"encoding/json"
	"github.com/manjada/com/db"
	"github.com/manjada/com/db/repo"
	"github.com/manjada/com/dto"
	"github.com/manjada/com/repo"
)

type AuditLogger interface {
	LogCreate(userId, username, ipAddress, module string, detail interface{}) error
	LogRead(userId, username, ipAddress, module string, detail interface{}) error
	LogUpdate(userId, username, ipAddress, module string, detail interface{}) error
	LogDelete(userId, username, ipAddress, module string, detail interface{}) error
	LogApproved(userId, username, ipAddress, module string, detail interface{}) error
	LogRejected(userId, username, ipAddress, module string, detail interface{}) error
}

func NewAuditLogService(DB db.DBConnector) AuditLogger {
	return &AuditLogService{Db: repo.NewBaseRepo(DB)}
}

type AuditLogService struct {
	Db repo.BaseRepoGorm
}

func (a *AuditLogService) LogCreate(userId, username, ipAddress, module string, detail interface{}) error {
	return a.log(userId, username, ipAddress, module, dto.AUDIT_ACTION_CREATE, detail)
}

func (a *AuditLogService) LogRead(userId, username, ipAddress, module string, detail interface{}) error {
	return a.log(userId, username, ipAddress, module, dto.AUDIT_ACTION_READ, detail)
}

func (a *AuditLogService) LogUpdate(userId, username, ipAddress, module string, detail interface{}) error {
	return a.log(userId, username, ipAddress, module, dto.AUDIT_ACTION_UPDATE, detail)
}

func (a *AuditLogService) LogDelete(userId, username, ipAddress, module string, detail interface{}) error {
	return a.log(userId, username, ipAddress, module, dto.AUDIT_ACTION_DELETE, detail)
}

func (a *AuditLogService) LogApproved(userId, username, ipAddress, module string, detail interface{}) error {
	return a.log(userId, username, ipAddress, module, dto.AUDIT_ACTION_APPROVE, detail)
}

func (a *AuditLogService) LogRejected(userId, username, ipAddress, module string, detail interface{}) error {
	return a.log(userId, username, ipAddress, module, dto.AUDIT_ACTION_REJECT, detail)
}

func (a *AuditLogService) log(userId, username, ipAddress, module, action string, data interface{}) error {
	detail, err := json.Marshal(data)
	if err != nil {
		return err
	}
	auditLog := repo.AuditLog{
		UserId:    userId,
		Username:  username,
		IpAddress: ipAddress,
		Module:    module,
		Action:    action,
		Detail:    string(detail),
	}
	return a.Db.Create(&auditLog).DbRepo.Error
}
