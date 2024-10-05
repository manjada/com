package db

import "gorm.io/gorm"

type DBConnector interface {
	Connect() error
	GetDB() *gorm.DB
}
