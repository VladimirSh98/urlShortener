package database

import (
	"database/sql"
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var DBConnection *sql.DB

func OpenConnectionDB() error {
	var err error
	DBConnection, err = sql.Open("pgx", config.DatabaseDSN)
	if err != nil {
		return err
	}
	return nil
}
