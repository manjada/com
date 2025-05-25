package db_adapter

import (
	"github.com/manjada/com/config"
	"github.com/manjada/com/db_adapter/gorm"
)

type DBAdapter interface {
	AutoMigrate(data interface{}) error
	Create(data interface{}) error
}

func NewDbAdapter() DBAdapter {
	cfg := config.GetConfig()
	var err error
	var db DBAdapter
	if cfg.DbConfig.Orm == "orm" && cfg.DbConfig.Type == "postgres" {
		db, err = gorm.NewPostgresGormAdapter(cfg.DbConfig.Host)
		if err != nil {
			panic(err)
		}
	}
	return db
}
