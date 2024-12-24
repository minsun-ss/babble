package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Cfg    *viper.Viper
	DBpool *DB
}

type DB struct {
	*gorm.DB
}

// NewConfig generates a new configuration that imports environment values,
// instantiates a new database pool connection and serves as the centralized
// configuration struct for the application.
func NewConfig() *Config {
	v := viper.New()
	v.SetEnvPrefix("BABEL")
	v.AutomaticEnv()

	dbpool := NewDB(v)
	return &Config{
		Cfg:    v,
		DBpool: dbpool,
	}
}

// NewDB generates a new database connection pool to the Babel database.
func NewDB(config *viper.Viper) *DB {
	host := config.GetString("DB_HOST")
	user := config.GetString("DB_USER")
	password := config.GetString("DB_PASSWORD")
	dbname := config.GetString("DB_DBNAME")
	port := config.GetString("DB_PORT")
	connection_string := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(connection_string), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("no connection established babel database: %v", err))
	}

	connpool, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("no connection established to babel database: %v", err))
	}

	connpool.SetMaxOpenConns(20)
	connpool.SetConnMaxLifetime(time.Hour)

	return &DB{db}
}
