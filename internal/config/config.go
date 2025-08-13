package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

const (
	SparkDir   = "/etc/spark"
	ConfigFile = "spark.yml"
)

type TLSConfig struct {
	Enabled  bool   `mapstructure:"enabled" yaml:"enabled"`
	CertFile string `mapstructure:"certFile" yaml:"certFile"`
	KeyFile  string `mapstructure:"keyFile" yaml:"keyFile"`
}

type XrayConfig struct {
	Address string `mapstructure:"address" yaml:"address"`
	Port    uint16 `mapstructure:"port" yaml:"port"`
}

type InfernoConfig struct {
	Enabled bool   `mapstructure:"enabled" yaml:"enabled"`
	Server  string `mapstructure:"server" yaml:"server"`
	Token   string `mapstructure:"token" yaml:"token"`
}

type Config struct {
	ID                        string        `mapstructure:"nodeID" yaml:"id"`
	Address                   string        `mapstructure:"address" yaml:"address"`
	Port                      uint16        `mapstructure:"port" yaml:"port"`
	Prefork                   bool          `mapstructure:"prefork" yaml:"prefork"`
	Token                     string        `mapstructure:"token" yaml:"token"`
	TLS                       TLSConfig     `mapstructure:"tls" yaml:"tls"`
	Xray                      XrayConfig    `mapstructure:"xray" yaml:"xray"`
	Inferno                   InfernoConfig `mapstructure:"inferno" yaml:"inferno"`
	NodeStatusUpdateFrequency time.Duration `mapstructure:"nodeStatusUpdateFrequency" yaml:"nodeStatusUpdateFrequency"`
	NodeLeaseDurationSeconds  time.Duration `mapstructure:"nodeLeaseDurationSeconds" yaml:"nodeLeaseDurationSeconds"`
	NodeStatusReportFrequency time.Duration `mapstructure:"nodeStatusReportFrequency" yaml:"nodeStatusReportFrequency"`
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
