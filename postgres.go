package mjd

import (
	"fmt"
	config2 "github.com/manjada/com/config"
	"github.com/manjada/com/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
}

func (p Postgres) Connect() {
	if Db == nil {
		config := config2.GetConfig().DbConfig
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s", config.Host, config.User, config.Pass, config.DbName, config.Port, config.Timezone)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			SkipDefaultTransaction: true,
		})
		if err != nil {
			log.Panic(err)
		}

		log.Info("Success connect to DB")
		if config.Debug {
			db = db.Debug()
		}
		Db = db
	}

}
