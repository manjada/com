package repo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/manjada/com/config"
	"github.com/manjada/com/svc"
	"github.com/oklog/ulid"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

var excludeTable = []string{"module_menus", "roles", "approvals", "approval_details", "approval_transactions", "approval_transaction_details"}

type TransactionModel struct {
	Id        string         `gorm:"primaryKey type:varchar(255)"`
	CreatedAt time.Time      `gorm:"index"`
	UpdatedAt time.Time      `gorm:"index"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (receive *TransactionModel) BeforeUpdate(tx *gorm.DB) error {
	receive.UpdatedAt = time.Now()
	return nil
}

func (receive *TransactionModel) AfterCreate(tx *gorm.DB) error {
	tableName := tx.Statement.Table
	for _, data := range excludeTable {
		if data == tableName {
			return nil
		}
	}
	var IsNeedApproval bool
	tx.Table("module_menus").Where("menu_code = ?", tableName).Select("is_need_approval").Scan(&IsNeedApproval)
	if !IsNeedApproval {
		config.Info(fmt.Sprintf("Approval not needed for table: %s", tableName))
		return nil
	}
	config.Info(fmt.Sprintf("After Create Table Name: %s and check approval", tableName))
	var data map[string]interface{}
	tx.Table(tableName).Where("id = ?", receive.Id).Scan(&data)
	var dataBin []byte
	dataBin, err := json.Marshal(&data)
	if err != nil {
		config.Error(err)
		return err
	}

	clientId, ok := data["client_id"]
	if !ok {
		config.Error(errors.New("client_id not found"))
		return errors.New("client_id not found")
	}
	err, emails := receive.buildApprovalTransaction(tx, tableName, clientId.(string), dataBin)
	if err != nil && err.Error() != "approval not found" {
		config.Error(err)
		return err
	}

	var bodyMail string
	tx.Table(EmailTemplate{}.TableName()).Where("template_key = ?", "approval_"+tableName).Select("body").Scan(&bodyMail)
	if bodyMail == "" {
		config.Info("email template not found")
	}

	err = svc.NewEmailService().SendEmail(data, emails, tableName+"_approval", nil, "", bodyMail, clientId.(string))
	if err != nil {
		config.Error(err)
		return err
	}
	return nil
}

func (receive *TransactionModel) buildApprovalTransaction(tx *gorm.DB, tableName string, clientId string, dataBin []byte) (error, []string) {
	var dataApproval []map[string]interface{}
	err := tx.Table("approvals").
		Select(`"approvals".id, "approval_details".approval_by, "approval_details".approval_name, "approval_details".client_id,
"approvals".module_menu_name, "approvals".module_menu_code`).
		Joins(`join approval_details on "approvals".id = "approval_details".approval_id`).
		Where(`"approvals".module_menu_code = ? and "approval_details".client_id = ?`, tableName, clientId).
		Scan(&dataApproval).Error
	if err != nil {
		return err, nil
	}

	if len(dataApproval) == 0 {
		return errors.New("approval not found"), nil
	}
	var approvalTransactionID string
	err = tx.Raw(`
    INSERT INTO approval_transactions (id, created_at, updated_at, approval_id, client_id, module_code, module_name, status, reference_id, total_approval, type, data)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    RETURNING id
`, receive.generateUlid().String(), time.Now(), time.Now(), dataApproval[0]["id"], dataApproval[0]["client_id"], dataApproval[0]["module_menu_code"], dataApproval[0]["module_menu_name"], "Pending", receive.Id, len(dataApproval), dataApproval[0]["type"], string(dataBin)).Scan(&approvalTransactionID).Error
	if err != nil {
		return err, nil
	}

	emails := []string{}
	for _, datas := range dataApproval {
		var approvalTransactionDetailID string
		err = tx.Raw(`
    INSERT INTO approval_transaction_details (id, created_at, updated_at, approval_transaction_id, approval_by, approval_by_name)
    VALUES (?, ?, ?, ?, ?, ?)
    RETURNING id
`, receive.generateUlid().String(), time.Now(), time.Now(), approvalTransactionID, datas["approval_by"], datas["approval_by_name"]).Scan(&approvalTransactionDetailID).Error

		if err != nil {
			return err, nil
		}
		var email string
		tx.Table("users").Select("email").Where("id = ?", datas["approval_by"]).Scan(&email)
		emails = append(emails, email)
	}
	return nil, emails
}

func (receive *TransactionModel) BeforeCreate(tx *gorm.DB) error {
	id := receive.generateUlid()

	if receive.Id == "" {
		receive.Id = id.String()
	} else {
		return errors.New("can't save invalid data")
	}
	return nil
}

func (receive *TransactionModel) generateUlid() ulid.ULID {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id
}
