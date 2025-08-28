package config

import (
	"fmt"
	"net"
	"strconv"
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

	LogLevel       string `mapstructure:"log_level"`
	MigrationsPath string `mapstructure:"migrations_path"`
}

func (c *Config) DSN() string {
	hostPort := net.JoinHostPort(c.DBHost, strconv.Itoa(c.DBPort))

	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		c.DBUser, c.DBPassword, hostPort, c.DBName, c.SSLMode)
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")

	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if cfg.MigrationsPath == "" {
		cfg.MigrationsPath = "./migrations"
	}

	return &cfg, nil
}
