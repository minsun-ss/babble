package utils

import (
	"time"

	"babel/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func DBPool(config *config.Config) *DB {
	host := config.GetString("DB_HOST")
	user := config.GetString("DB_USER")
	password := config.GetString("DB_PASSWORD")
	dbname := config.GetString("DB_DBNAME")
	port := config.GetString("DB_PORT")
	connection_string := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(connection_string), &gorm.Config{})

	if err != nil {
		return nil
	}

	connpool, err := db.DB()

	if err != nil {
		return nil
	}

	connpool.SetMaxOpenConns(20)
	connpool.SetConnMaxLifetime(time.Hour)

	return &DB{db}
}
