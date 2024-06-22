package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/manjada/com/log"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var (
	viperInit *viper.Viper
	conf      *Config
	ext       = "yaml"
)

type Config struct {
	AppJwt struct {
		AccessSecret  string `mapstructure:"access_secret"`
		RefreshSecret string `mapstructure:"refresh_secret"`
	} `mapstructure:"app_secret"`

	FatSecret struct {
		ClientId  string `mapstructure:"client_id"`
		ClientKey string `mapstructure:"client_key"`
		UrlToken  string `mapstructure:"url_token"`
	} `mapstructure:"fatsecret"`

	AppHost struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"app_host"`

	DbConfig struct {
		Host     string `mapstructure:"host"`
		User     string `mapstructure:"user"`
		Pass     string `mapstructure:"pass"`
		DbName   string `mapstructure:"dbname"`
		Port     int    `mapstructure:"port"`
		Timezone string `mapstructure:"timezone"`
		Debug    bool   `mapstructure:"debug"`
		Type     string `mapstructure:"type"`
	} `mapstructure:"db_config"`

	NoSqlConfig struct {
		Host  string `mapstructure:"host"`
		User  string `mapstructure:"user"`
		Pass  string `mapstructure:"pass"`
		Debug bool   `mapstructure:"debug"`
	} `mapstructure:"no_sql_config"`

	LogFile struct {
		PathWindows string `mapstructure:"path_windows"`
		PathUnix    string `mapstructure:"path_unix"`
		Level       string `mapstructure:"level"`
	} `mapstructure:"log_file_config"`

	Swagger struct {
		Enable bool `mapstructure:"enable"`
	} `mapstructure:"swagger"`

	Redis struct {
		Address string `mapstructure:"address"`
		Pass    string `mapstructure:"pass"`
		Index   int    `mapstructure:"index"`
	} `mapstructure:"redis"`

	Smtp struct {
		Host     string `mapstructure:"host"`
		User     string `mapstructure:"user"`
		Port     int    `mapstructure:"port"`
		Password string `mapstructure:"password"`
	} `mapstructure:"smtp_config"`

	Message struct {
		Host string `mapstructure:"host"`
		User string `mapstructure:"user"`
		Pass string `mapstructure:"pass"`
	} `mapstructure:"message"`
}

func init() {
	log.Info("Start Load Config")
	application := "application"
	if viperInit == nil {

		if os.Getenv("WS_ENV") != "" {
			application = fmt.Sprintf("%s-%s", application, os.Getenv("WS_ENV"))
		}

		log.Info("Config path => ", application)
		viperInit = viper.New()
		viperInit.SetConfigType(ext)
		viperInit.AddConfigPath("./resource")

		viperInit.SetConfigName(application)
		viperInit.AutomaticEnv()
		err := viperInit.ReadInConfig()

		if err != nil {
			log.Error(err)
		}

		viperInit.OnConfigChange(func(in fsnotify.Event) {
			log.Info(fmt.Sprintf("Config file changed: %s", in.Name))
		})
		viperInit.WatchConfig()

		log.Info("Config loaded successfully...")
		log.Info("Getting environment variables...")
		for _, k := range viperInit.AllKeys() {
			value := viperInit.GetString(k)
			if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
				viperInit.Set(k, getEnvOrPanic(strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")))
			}
		}

		err = viperInit.Unmarshal(&conf)
		if err != nil {
			log.Panic(err)
		}
		log.Info("Config data ", conf.DbConfig.Host)

	}
}

func getEnvOrPanic(env string) string {
	res := os.Getenv(env)
	if len(res) == 0 {
		log.Panic(fmt.Errorf("Mandatory env variable not found:" + env))
	}
	return res
}

func GetConfig() *Config {
	return conf
}
