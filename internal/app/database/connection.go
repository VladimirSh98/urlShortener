package database

import (
	"database/sql"
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var DBConnection = DBConnectionStruct{}

func (db *DBConnectionStruct) OpenConnection() error {
	var err error
	db.conn, err = sql.Open("pgx", config.DatabaseDSN)
	if err != nil {
		return err
	}
	return nil
}

func (db *DBConnectionStruct) CloseConnection() {
	db.conn.Close()
}
