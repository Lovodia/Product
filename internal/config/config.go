package config

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type DBConfig struct {
	Host     string `mapstructure:"dbhost"`
	Port     int    `mapstructure:"dbport"`
	User     string `mapstructure:"dbuser"`
	Password string `mapstructure:"dbpassword"`
	Name     string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

func (c *DBConfig) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("DBHost is required")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("DBPort must be between 1 and 65535")
	}
	if c.User == "" {
		return fmt.Errorf("DBUser is required")
	}
	if c.Password == "" {
		return fmt.Errorf("DBPassword is required")
	}
	if c.Name == "" {
		return fmt.Errorf("DBName is required")
	}
	validSSLModes := map[string]bool{
		"disable": true, "require": true, "verify-ca": true, "verify-full": true,
	}
	if !validSSLModes[c.SSLMode] {
		return fmt.Errorf("invalid SSLMode: %s", c.SSLMode)
	}
	return nil
}

func (c *DBConfig) DSN() string {
	hostPort := net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		c.User, c.Password, hostPort, c.Name, c.SSLMode)
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

func (c *ServerConfig) Validate() error {
	if c.Port == "" {
		return fmt.Errorf("Server.Port is required")
	}
	if _, err := strconv.Atoi(c.Port); err != nil {
		return fmt.Errorf("Server.Port must be a number")
	}
	return nil
}

type LoggerConfig struct {
	Level string `mapstructure:"log_level"`
}

func (c *LoggerConfig) Validate() error {
	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLogLevels[strings.ToLower(c.Level)] {
		return fmt.Errorf("invalid LogLevel: %s", c.Level)
	}
	return nil
}

type MigrationConfig struct {
	Path string `mapstructure:"migrations_path"`
}

func (c *MigrationConfig) Validate() error {
	if c.Path == "" {
		return fmt.Errorf("migrations_path is required")
	}
	return nil
}

type Config struct {
	DB        DBConfig        `mapstructure:",squash"`
	Server    ServerConfig    `mapstructure:"server"`
	Logger    LoggerConfig    `mapstructure:",squash"`
	Migration MigrationConfig `mapstructure:",squash"`
}

func (c *Config) Validate() error {
	if err := c.DB.Validate(); err != nil {
		return err
	}
	if err := c.Server.Validate(); err != nil {
		return err
	}
	if err := c.Logger.Validate(); err != nil {
		return err
	}
	if err := c.Migration.Validate(); err != nil {
		return err
	}
	return nil
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

	if cfg.Migration.Path == "" {
		cfg.Migration.Path = "./migrations"
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}
