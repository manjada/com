package svc

import (
	"encoding/json"
	"github.com/manjada/com/db"
	"github.com/manjada/com/db/repo"
	"github.com/manjada/com/dto"
	repo2 "github.com/manjada/com/repo"
)

type AuditLogger interface {
	LogCreate(auth dto.AccessDetail, module string, detail interface{}) error
	LogRead(auth dto.AccessDetail, module string, detail interface{}) error
	LogUpdate(auth dto.AccessDetail, module string, detail interface{}) error
	LogDelete(auth dto.AccessDetail, module string, detail interface{}) error
	LogApproved(auth dto.AccessDetail, module string, detail interface{}) error
	LogRejected(auth dto.AccessDetail, module string, detail interface{}) error
}

func NewAuditLogService(DB db.DBConnector) AuditLogger {
	return &AuditLogService{Db: repo.NewBaseRepo(DB)}
}

type AuditLogService struct {
	Db repo.BaseRepoGorm
}

func (a *AuditLogService) LogCreate(auth dto.AccessDetail, module string, detail interface{}) error {
	return a.log(auth, module, dto.AUDIT_ACTION_CREATE, detail)
}

func (a *AuditLogService) LogRead(auth dto.AccessDetail, module string, detail interface{}) error {
	return a.log(auth, module, dto.AUDIT_ACTION_READ, detail)
}

func (a *AuditLogService) LogUpdate(auth dto.AccessDetail, module string, detail interface{}) error {
	return a.log(auth, module, dto.AUDIT_ACTION_UPDATE, detail)
}

func (a *AuditLogService) LogDelete(auth dto.AccessDetail, module string, detail interface{}) error {
	return a.log(auth, module, dto.AUDIT_ACTION_DELETE, detail)
}

func (a *AuditLogService) LogApproved(auth dto.AccessDetail, module string, detail interface{}) error {
	return a.log(auth, module, dto.AUDIT_ACTION_APPROVE, detail)
}

func (a *AuditLogService) LogRejected(auth dto.AccessDetail, module string, detail interface{}) error {
	return a.log(auth, module, dto.AUDIT_ACTION_REJECT, detail)
}

func (a *AuditLogService) log(auth dto.AccessDetail, module, action string, data interface{}) error {
	detail, err := json.Marshal(data)
	if err != nil {
		return err
	}
	auditLog := repo2.AuditLog{
		UserId:    auth.UserId,
		Name:      auth.Name,
		IpAddress: auth.IpAddress,
		ClientId:  auth.ClientId,
		Module:    module,
		Action:    action,
		Detail:    string(detail),
	}
	return a.Db.Create(&auditLog).DbRepo.Error
}
