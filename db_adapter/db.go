package db_adapter

import (
	"fmt"
	"github.com/manjada/com/config"
	"github.com/manjada/com/db_adapter/gorm"
	_interface "github.com/manjada/com/db_adapter/interface"
	"github.com/manjada/com/db_adapter/native"
)

func NewDbAdapter() _interface.DBAdapter {
	cfg := config.GetConfig()
	var err error
	var db _interface.DBAdapter
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s", cfg.DbConfig.Host, cfg.DbConfig.User, cfg.DbConfig.Pass, cfg.DbConfig.DbName, cfg.DbConfig.Port, cfg.DbConfig.Timezone)
	switch cfg.DbConfig.Orm {
	case "gorm":
		if cfg.DbConfig.Type == "postgres" {
			db, err = gorm.NewPostgresGormAdapter(dsn)
			if err != nil {
				config.Panic(err)
			}
		}
	case "sqlx":
		if cfg.DbConfig.Type == "postgres" {
			db, err = native.NewPostgresAdapter(dsn)
			if err != nil {
				config.Panic(err)
			}
		}
	}
	return db
}
