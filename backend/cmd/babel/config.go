/*
Package config holds all application configuration.

This includes:
- Environment variables
- Database configuration
- Logging configuration
*/
package babel

import (
	"embed"
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Config holds configuration information about the app
type Config struct {
	Cfg     *viper.Viper
	DB      *gorm.DB
	BabelFS embed.FS
	ApiCfg  *huma.Config
}

// NewConfig generates a new configuration that imports environment values,
// instantiates a new database pool connection and serves as the centralized
// configuration struct for the application.
func NewConfig(static embed.FS) *Config {
	v := viper.New()
	v.SetEnvPrefix("BABEL")
	v.AutomaticEnv()

	// set up gorm config
	db := NewDB(v)

	// set up api config
	apicfg := huma.DefaultConfig("Babel API", "1.0.0")
	apicfg.DocsPath = "/api/docs/"
	apicfg.OpenAPIPath = "/api/openapi/"

	return &Config{
		Cfg:     v,
		DB:      db,
		BabelFS: static,
		ApiCfg:  &apicfg,
	}
}

// NewDB generates a new database connection pool to the Babel database
// from environment variables.
func NewDB(config *viper.Viper) *gorm.DB {
	host := config.GetString("DB_HOST")
	user := config.GetString("DB_USER")
	password := config.GetString("DB_PASSWORD")
	dbname := config.GetString("DB_DBNAME")
	port := config.GetString("DB_PORT")
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
		user, password, host, port, dbname)

	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("no connection established babel database: %v", err))
	}

	// set up connection pool settings
	connPool, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("no connection established to babel database: %v", err))
	}

	connPool.SetMaxOpenConns(20)
	connPool.SetConnMaxLifetime(time.Hour)

	return db
}
