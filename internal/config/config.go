package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type Config struct {
	DBHost     string `mapstructure:"dbhost"`
	DBPort     int    `mapstructure:"dbport"`
	DBUser     string `mapstructure:"dbuser"`
	DBPassword string `mapstructure:"dbpassword"`
	DBName     string `mapstructure:"dbname"`
	SSLMode    string `mapstructure:"sslmode"`

	Server ServerConfig `mapstructure:"server"`

	LogLevel string `mapstructure:"log_level"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")

	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBName == "" {
		return nil, fmt.Errorf("missing database configuration fields")
	}
	if cfg.Server.Port == "" {
		return nil, fmt.Errorf("missing server port")
	}
	return &cfg, nil
}
