package mjd

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
}

func (p Postgres) Connect() {
	if Db == nil {
		config := GetConfig().DbConfig
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s", config.Host, config.User, config.Pass, config.DbName, config.Port, config.Timezone)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			SkipDefaultTransaction: true,
		})
		if err != nil {
			Panic(err)
		}

		Info("Success connect to DB")
		Db = db
	}

}
