package repo

import (
	"errors"
	"fmt"
	"github.com/oklog/ulid"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type TransactionModel struct {
	Id         string         `gorm:"primaryKey type:varchar(255)"`
	CreatedAt  time.Time      `gorm:"index"`
	UpdatedAt  time.Time      `gorm:"index"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	ModuleName string         `gorm:"-"`
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

func (receive *TransactionModel) AfterCreate(tx *gorm.DB) error {
	if receive.ModuleName != "" {
		var IsNeedApproval bool
		tx.Select("is_need_approval").Where("menu_code = ?", receive.ModuleName).Table(`module_menus`).Scan(&IsNeedApproval)
		if IsNeedApproval {
			approval := Approval{
				ApprovalName:     "Pending",
				ApproveBy:        "",
				Status:           "Pending",
				TransactionId:    receive.Id,
				ModuleName:       receive.ModuleName,
				Data:             "",
				ApprovalDate:     nil,
				RejectDate:       nil,
				PendingDate:      nil,
				ApprovalDuration: 0,
				Description:      "",
				DelegationId:     "",
				DelegationName:   "",
			}
			tx.Create(&approval)
		}
		fmt.Printf("Record created by module: %s\n", receive.ModuleName)
	} else {
		fmt.Println("Module name not set")
	}
	return nil
}

func (receive *TransactionModel) AfterUpdate(tx *gorm.DB) error {
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
