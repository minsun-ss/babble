package config

import "github.com/spf13/viper"

type Config struct {
	*viper.Viper
}

func NewConfig() *Config {
	v := viper.New()
	v.SetEnvPrefix("BABEL")
	v.AutomaticEnv()
	return &Config{v}
}
