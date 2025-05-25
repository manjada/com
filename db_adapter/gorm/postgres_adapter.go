package gorm

import (
	"github.com/manjada/com/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresGormAdapter struct {
	db *gorm.DB
}

func NewPostgresGormAdapter(dsn string) (*PostgresGormAdapter, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &PostgresGormAdapter{db: db}, nil
}

func (a *PostgresGormAdapter) Create(data interface{}) error {
	return a.db.Create(data).Error
}

func (a *PostgresGormAdapter) AutoMigrate(data interface{}) error {
	if err := a.db.AutoMigrate(data); err != nil {
		config.Error(err)
		return err
	}
	return nil
}
