package database

import (
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"github.com/pressly/goose/v3"
)

// UpgradeMigrations applies migrations to database
func (db *DBConnectionStruct) UpgradeMigrations() error {
	err := goose.Up(db.Conn, config.DefaultConfigValues.MigrationsDir)
	if err != nil {
		return err
	}
	return nil
}
