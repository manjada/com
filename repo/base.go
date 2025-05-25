package repo

import (
	"errors"
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
