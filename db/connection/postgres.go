package connection

import (
	"fmt"
	confProp "github.com/manjada/com/config"
	"github.com/manjada/com/db/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var GormDb PostgresGorm

type PostgresGorm struct {
	cfg *config.PostgresConfig
	DB  *gorm.DB
}

func (p *PostgresGorm) Connect() error {

	// If a connection already exists, return immediately
	if p.DB != nil {
		GormDb = *p
		return nil
	}

	p.cfg = config.NewPostgresConfig()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s", p.cfg.Host, p.cfg.User, p.cfg.Password, p.cfg.DBName, p.cfg.Port, p.cfg.TimeZone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		confProp.Panic(err)
	}

	confProp.Info("Success connect to DB")
	if confProp.GetConfig().DbConfig.Debug {
		db = db.Debug()
	}
	p.DB = db
	return nil
}

func (p *PostgresGorm) GetDB() *gorm.DB {
	return p.DB
}
