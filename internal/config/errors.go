package config

import "errors"

var (
	ErrDBHostRequired         = errors.New("DBHost is required")
	ErrDBPortInvalid          = errors.New("DBPort must be betewen 1 and 65535")
	ErrDBUserRequired         = errors.New("DBUser is required")
	ErrDBPasswordRequired     = errors.New("DBPassword is required")
	ErrDBNameRequired         = errors.New("DBName is required")
	ErrInvalidSSLMode         = errors.New("invalide SSLMode")
	ErrServerPortRequired     = errors.New("Server.Port is required")
	ErrServerPortMustBeNumber = errors.New("Server.Port must be a number")
	ErrInvalidLogLevel        = errors.New("invalid LogLevel")
	ErrMigrationsPathRequired = errors.New("migrations_path is required")
)
