package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	APIServerAddr string
	XrayAPIAddr   string
	Domain        string
	Secret        string
	Development   bool
}

var AppConfig *Config

func LoadConfig() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	viper.AutomaticEnv()
	AppConfig = &Config{
		APIServerAddr: viper.GetString("API_SERVER_ADDR"),
		XrayAPIAddr:   viper.GetString("XRAY_API_ADDR"),
		Domain:        viper.GetString("DOMAIN"),
		Secret:        viper.GetString("SECRET"),
		Development:   viper.GetBool("DEV"),
	}
}
