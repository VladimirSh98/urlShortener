package database

import (
	"database/sql"
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBConnectionStruct struct {
	Conn *sql.DB
}

var DBConnection = DBConnectionStruct{}

func (db *DBConnectionStruct) OpenConnection() error {
	var err error
	db.Conn, err = sql.Open("pgx", config.DatabaseDSN)
	if err != nil {
		return err
	}
	return nil
}

func (db *DBConnectionStruct) CloseConnection() {
	db.Conn.Close()
}
