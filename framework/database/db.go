package database

import (
	"github.com/video-encoder/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitGorm(databaseConfig *Config) *gorm.DB {

	dsn := getDSN(databaseConfig)

	logLevel := logger.Error
	if config.IsLocal() {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logLevel)})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func getDSN(config *Config) string {
	return "host=" + config.Hostname + " user=" + config.Username + " password=" + config.Password +
		" dbname=" + config.Database + " port=" + config.Port + " sslmode=disable TimeZone=UTC"
}
