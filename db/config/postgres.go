package config

import confProp "github.com/manjada/com/config"

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	TimeZone string
}

func NewPostgresConfig() *PostgresConfig {
	dbConfig := confProp.GetConfig().DbConfig
	return &PostgresConfig{
		Host:     dbConfig.Host,
		Port:     dbConfig.Port,
		User:     dbConfig.User,
		Password: dbConfig.Pass,
		DBName:   dbConfig.DbName,
		TimeZone: dbConfig.Timezone,
	}
}
