package db_adapter

import (
	"github.com/manjada/com/config"
	"github.com/manjada/com/db_adapter/gorm"
	_interface "github.com/manjada/com/db_adapter/interface"
)

func NewDbAdapter() _interface.DBAdapter {
	cfg := config.GetConfig()
	var err error
	var db _interface.DBAdapter
	switch cfg.DbConfig.Orm {
	case "gorm":
		if cfg.DbConfig.Type == "postgres" {
			db, err = gorm.NewPostgresGormAdapter(cfg.DbConfig.Host)
			if err != nil {
				config.Panic(err)
			}
		}

	}
	return db
}
