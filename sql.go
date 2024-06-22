package mjd

import (
	config2 "github.com/manjada/com/config"
	"github.com/manjada/com/log"
	"gorm.io/gorm"
)

var Db *gorm.DB

type DbConfig interface {
	Connect()
}

func init() {
	log.Info("Start Connecting to DB")
	var dbCon DbConfig
	config := config2.GetConfig()
	switch config.DbConfig.Type {
	case "postgresql":
		dbCon = Postgres{}
	default:
		log.log.Panic("Database type not found, please define in properties")
	}
	dbCon.Connect()
}
