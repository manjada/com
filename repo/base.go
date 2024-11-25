package repo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/manjada/com/config"
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
	err = receive.buildApprovalTransaction(tx, tableName, clientId.(string), dataBin)
	if err != nil && err.Error() != "approval not found" {
		config.Error(err)
		return err
	}
	return nil
}

func (receive *TransactionModel) buildApprovalTransaction(tx *gorm.DB, tableName string, clientId string, dataBin []byte) error {
	var dataApproval []map[string]interface{}
	err := tx.Table("approvals").
		Select(`"approvals".id, "approval_details".approval_by, "approval_details".approval_name, "approval_details".client_id`).
		Joins(`join approval_details on "approvals".id = "approval_details".approval_id`).
		Where(`"approvals".module_menu_code = ? and "approval_details".client_id = ?`, tableName, clientId).
		Scan(&dataApproval).Error
	if err != nil {
		return err
	}

	if len(dataApproval) == 0 {
		return errors.New("approval not found")
	}
	err = tx.Table("approval_transactions").Create(map[string]interface{}{
		"id":             receive.generateUlid(),
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
			"id":                      receive.generateUlid(),
			"approval_transaction_id": datas["id"],
			"approval_by":             datas["approval_by"],
			"approval_by_name":        datas["approval_by_name"],
		}).Error

		if err != nil {
			return err
		}
	}
	return nil
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
