package pkg_config

import (
	"GoRestify/pkg/dictionary"
	"GoRestify/pkg/models"
	"GoRestify/pkg/pkg_consts"
	"GoRestify/pkg/pkg_log"
	"GoRestify/pkg/pkg_redis"
	"GoRestify/pkg/pkg_types"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config of Application
var Config *ConfigInstance

// ConfigInstance its instance of an pkg_config
type ConfigInstance struct {

	// Debug mode
	IsDebug bool

	DB         *gorm.DB
	ActivityDB *gorm.DB

	Redis pkg_redis.RedisCon

	// Record Activity
	ActivityActive bool
	ActivityCh     chan models.Activity

	EnumLists map[string]interface{}

	// basic auth keys
	BasicAuthUsername string
	BasicAuthPassword string

	// JWT Secret Key
	JWTSecretKey string

	HTTPClient *http.Client
}

// Cnf its config struct to pkg_config
type Cnf struct {

	// Debug mode
	IsDebug string

	DbDSN         string
	ActivityDbDSN string
	RedisAddress  string

	SettingActive    bool
	SettingDomainApp pkg_types.Enum

	// MustBeInTypes its a list in core actions, to validate enum fields in database
	MustBeInTypes map[string]interface{}

	// basic auth keys, if project use basic auth in header
	BasicAuthUsername string
	BasicAuthPassword string

	// JWT Secret Key, if project use token in header
	JWTSecretKey string
}

// InitConfig its fill a pkg_config value
func InitConfig(cnf Cnf) (config *ConfigInstance, err error) {

	config = new(ConfigInstance) // Initialize Config

	config.IsDebug = cnf.IsDebug == "debug"

	// init a instance for http client
	config.HTTPClient = new(http.Client)

	if config.IsDebug {
		config.HTTPClient.Transport = &pkg_log.LoggingTransport{}
	}
	pkg_consts.HTTPClient = config.HTTPClient

	// init a connection for activity and setting
	activityConnection(cnf, config)
	settingConnection(cnf, config)

	// redis connection to fetch phrases in DBTranslate
	if config.Redis, err = pkg_redis.ConnectRedis(cnf.RedisAddress, config.IsDebug); err != nil {
		return
	}

	// init dictionary to translate backend terms
	if err = dictionary.Init(); err != nil {
		return
	}

	// Init server log
	pkg_log.Init(pkg_consts.LogServerOutput, config.IsDebug)

	// This action require to validation
	config.EnumLists = cnf.MustBeInTypes

	// this Env for basic Auth Middleware
	config.BasicAuthUsername = cnf.BasicAuthUsername
	config.BasicAuthPassword = cnf.BasicAuthPassword

	// this env for jwt secret key in middleware
	config.JWTSecretKey = cnf.JWTSecretKey

	Config = config

	return
}

func activityConnection(cnf Cnf, config *ConfigInstance) (err error) {

	// if database not active ignore ActivityDB connection
	if cnf.ActivityDbDSN != "" {
		config.ActivityActive = true
		var logLevel logger.LogLevel

		switch config.IsDebug {
		case true:
			logLevel = logger.Info
		case false:
			logLevel = logger.Silent
		}

		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second, // Slow SQL threshold
				LogLevel:      logLevel,    // Log level
				Colorful:      true,        // Disable color
			},
		)

		config.ActivityDB, err = gorm.Open(mysql.Open(cnf.ActivityDbDSN),
			&gorm.Config{
				Logger: newLogger,
			})
		if err != nil {
			return
		}

		config.ActivityDB.Table(models.ActivityTable).AutoMigrate(&models.Activity{})

		config.ActivityCh = make(chan models.Activity, 1)
	}
	return
}

func settingConnection(cnf Cnf, config *ConfigInstance) (err error) {

	// if database not active ignore SettingDB connection
	if cnf.SettingActive {
		var logLevel logger.LogLevel

		switch config.IsDebug {
		case true:
			logLevel = logger.Info
		case false:
			logLevel = logger.Silent
		}

		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second, // Slow SQL threshold
				LogLevel:      logLevel,    // Log level
				Colorful:      true,        // Disable color
			},
		)

		config.DB, err = gorm.Open(mysql.Open(cnf.DbDSN),
			&gorm.Config{
				Logger: newLogger,
			})
		if err != nil {
			return
		}

		config.DB.Table(models.SettingTable).AutoMigrate(&models.Setting{})

		if cnf.SettingDomainApp == "" {
			log.Fatal("SettingDomainApp is require for init setting")
		}
	}
	return
}
