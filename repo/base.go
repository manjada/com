package repo

import (
	"encoding/json"
	"errors"
	"github.com/manjada/com/dto"
	"github.com/oklog/ulid"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

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
	var data map[string]interface{}
	tx.Table(tableName).Where("id = ?", receive.Id).First(&data)
	var dataBin []byte
	dataBin, err := json.Marshal(&data)
	if err != nil {
		dto.ErrorDb(err)
		return err
	}
	err = receive.buildApprovalTransaction(tx, tableName, data, dataBin)
	if err != nil {
		dto.ErrorDb(err)
		return err
	}
	return nil
}

func (receive *TransactionModel) buildApprovalTransaction(tx *gorm.DB, tableName string, data map[string]interface{}, dataBin []byte) error {
	var dataApproval []map[string]interface{}
	err := tx.Table("approvals").
		Select(`"approvals".id, "approval_details".approval_by, "approval_details".approval_name, "approval_details".client_id`).
		Joins(`left join "approval_details".approval_id = "approvals".id`).
		Where(`"approvals".module_menu_code = ? and "approval_details".client_id = ?`, tableName, data["client_id"].(string)).
		Find(&dataApproval).Error
	if err != nil {
		return err
	}

	err = tx.Table("approval_transactions").Create(map[string]interface{}{
		"approval_id":    dataApproval[0]["id"],
		"client_id":      dataApproval[0]["client_id"],
		"module_code":    dataApproval[0]["module_code"],
		"status":         "Pending",
		"reference_id":   receive.Id,
		"total_approval": len(dataApproval),
		"type":           dataApproval[0]["type"],
		"data":           string(dataBin),
	}).Error
	if err != nil {
		return err
	}

	for _, datas := range dataApproval {
		err = tx.Table("approval_transaction_details").Create(map[string]interface{}{
			"approval_transaction_id": datas["id"],
			"approval_by":             data["approval_by"],
			"approval_by_name":        data["approval_by_name"],
		}).Error

		if err != nil {
			return err
		}
	}
	return nil
}

func (receive *TransactionModel) BeforeCreate(tx *gorm.DB) error {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)

	if receive.Id == "" {
		receive.Id = id.String()
	} else {
		return errors.New("can't save invalid data")
	}
	return nil
}
