package model

import (
	"errors"
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
