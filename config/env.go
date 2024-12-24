package config

import (
	"babel/utils"

	"github.com/spf13/viper"
)

type Config struct {
	Cfg    *viper.Viper
	DBpool *utils.DB
}

func NewConfig() *Config {
	v := viper.New()
	v.SetEnvPrefix("BABEL")
	v.AutomaticEnv()

	dbpool := utils.NewDB(v)
	return &Config{
		Cfg:    v,
		DBpool: dbpool,
	}
}
