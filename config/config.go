package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	TLS struct {
		Enabled  bool   `mapstructure:"enabled"`
		CertFile string `mapstructure:"cert_file"`
		KeyFile  string `mapstructure:"key_file"`
	} `mapstructure:"tls"`

	Server struct {
		Addr    string `mapstructure:"addr"`
		Prefork bool   `mapstructure:"prefork"`
	} `mapstructure:"server"`

	Xray struct {
		Addr string `mapstructure:"addr"`
	} `mapstructure:"xray"`

	Auth struct {
		Secret string `mapstructure:"secret"`
	} `mapstructure:"auth"`
}

var AppConfig *Config

func LoadConfig(configPath string) error {
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	return nil
}
