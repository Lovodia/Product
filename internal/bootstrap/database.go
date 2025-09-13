package bootstrap

import (
	"github.com/Lovodia/Product/internal/config"
	"github.com/Lovodia/Product/internal/db"
)

func InitDatabase(cfg *config.Config) (*db.Database, error) {
	return db.New(cfg)
}
