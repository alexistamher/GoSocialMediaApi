package db

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func StartDB() *gorm.DB {
	var DB *gorm.DB
	var err error

	DB, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	return DB
}
