package mjd

import (
	"gorm.io/gorm"
)

var Db *gorm.DB

type DbConfig interface {
	Connect()
}

func init() {
	Info("Start Connecting to DB")
	var dbCon DbConfig
	config := GetConfig()
	switch config.DbConfig.Type {
	case "postgresql":
		dbCon = Postgres{}
	default:
		log.Panic("Database type not found, please define in properties")
	}
	dbCon.Connect()
}
