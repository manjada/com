package gorm

import (
	"github.com/manjada/com/config"
	_interface "github.com/manjada/com/db_adapter/interface"
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
	if config.GetConfig().DbConfig.Debug {
		db = db.Debug()
	}
	return &PostgresGormAdapter{db: db}, nil
}

func (a *PostgresGormAdapter) AutoMigrate(data interface{}) error {
	if err := a.db.AutoMigrate(data); err != nil {
		config.Error(err)
		return err
	}
	return nil
}

func (a *PostgresGormAdapter) Where(query interface{}, args ...interface{}) _interface.DBAdapter {
	a.db = a.db.Where(query, args...)
	return a
}

func (a *PostgresGormAdapter) resetDB() {
	a.db = a.db.Session(&gorm.Session{})
}

func (a *PostgresGormAdapter) First(dest interface{}) error {
	// Execute the query on the current *gorm.DB instance
	err := a.db.First(dest).Error
	// Reset the *gorm.DB instance after the operation
	a.resetDB()
	return err
}

func (a *PostgresGormAdapter) Create(data interface{}) error {
	return a.db.Create(data).Error
}
