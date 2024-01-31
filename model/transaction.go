package model

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type TransactionModel struct {
	Id        string         `gorm:"primaryKey type:varchar(255)"`
	CreatedAt time.Time      `gorm:"index"`
	UpdatedAt time.Time      `gorm:"index"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (receive *TransactionModel) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewUUID()

	if err == nil && receive.Id == "" {
		receive.Id = id.String()
	} else {
		err = errors.New("can't save invalid data")
		return err
	}
	return nil
}
