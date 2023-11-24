package pkg_sql

import (
	"GoRestify/pkg/pkg_config"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MySQLConnectDB initiate the db connection in gorm
func MySQLConnectDB(DSN string) (DB *gorm.DB) {
	var err error

	var logLevel logger.LogLevel
	switch pkg_config.Config.IsDebug {
	case false:
		logLevel = logger.Silent
	case true:
		logLevel = logger.Info
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: 10 * time.Second, // Slow SQL threshold
			LogLevel:      logLevel,         // Log level
			Colorful:      true,             // Disable color
		},
	)

	DB, err = gorm.Open(mysql.Open(DSN),
		&gorm.Config{
			Logger: newLogger,
		})

	if err != nil {
		log.Fatalf("%v connection failed: %v\n", DSN, err.Error())
	}

	writeDB, err := DB.DB()
	if err != nil {
		log.Fatalln(err.Error())
	}

	writeDB.SetMaxIdleConns(250)
	writeDB.SetMaxOpenConns(250)
	writeDB.SetConnMaxLifetime(30 * time.Minute)

	return
}
