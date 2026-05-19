package database

import (
	"ImageCaptioningService/config"
	"ImageCaptioningService/services/DatabaseService/models"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

func Connect() error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.Mysql.USER, cfg.Mysql.PASSWORD, cfg.Mysql.HOST, cfg.Mysql.PORT, cfg.Mysql.DB, cfg.Mysql.CHARSET, cfg.Mysql.PARSE_TIME, cfg.Mysql.TIMEZONE)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	migrateError := db.AutoMigrate(models.Request{}, models.RequestMetadata{})
	if migrateError != nil {
		return migrateError
	}

	dbConn = db
	return nil
}

func GetConnection() (*gorm.DB, error) {
	if dbConn == nil {

		err := Connect()
		if err != nil {
			return nil, errors.New(fmt.Sprintf("database connection is not initialized : %v", err.Error()))
		}
	}
	return dbConn, nil
}

func CloseDatabase(db *gorm.DB) error {
	mySqlDB, err := db.DB()
	if err != nil {
		return err
	}
	err = mySqlDB.Close()
	if err != nil {
		return err
	}

	return nil
}
