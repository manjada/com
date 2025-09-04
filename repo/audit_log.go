package repo

import (
	"github.com/manjada/com/dto"
)

type AuditLog struct {
	TransactionModel
	UserId    string
	Name      string
	IpAddress string
	ClientId  string
	Module    string
	Action    string
	Detail    string
}

type AuditLogger interface {
	LogCreate(auth dto.AccessDetail, module string, detail interface{}) error
	LogRead(auth dto.AccessDetail, module string, detail interface{}) error
	LogUpdate(auth dto.AccessDetail, module string, detail interface{}) error
	LogDelete(auth dto.AccessDetail, module string, detail interface{}) error
	LogApproved(auth dto.AccessDetail, module string, detail interface{}) error
	LogRejected(auth dto.AccessDetail, module string, detail interface{}) error
}

/*func NewAuditLogService(DB db_adapter.DBConnector) AuditLogger {
	return &AuditLogService{Db: repo.NewBaseRepo(DB)}
}*/

type AuditLogService struct {
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
	return nil
	/*detail, err := json.Marshal(data)
	if err != nil {
		return err
	}
	auditLog := AuditLog{
		UserId:    auth.UserId,
		Name:      auth.Name,
		IpAddress: auth.IpAddress,
		Module:    module,
		Action:    action,
		Detail:    string(detail),
	}
	return a.Db.Create(&auditLog).DbRepo.Error*/
}
