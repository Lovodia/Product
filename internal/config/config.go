package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string `mapstructure:"dbhost"`
	DBPort     int    `mapstructure:"dbport"`
	DBUser     string `mapstructure:"dbuser"`
	DBPassword string `mapstructure:"dbpassword"`
	DBName     string `mapstructure:"dbname"`
	SSLMode    string `mapstructure:"sslmode"`

	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	err := viper.Unmarshal(&cfg)
	return &cfg, err
}
