package utils

import "github.com/spf13/viper"

type Config struct {
	*viper.Viper
}

func GetConfig() *Config {
	v := viper.New()
	v.SetEnvPrefix("BABEL")
	v.AutomaticEnv()
	return &Config{v}
}
