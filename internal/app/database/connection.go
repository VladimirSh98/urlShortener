package database

import (
	"database/sql"

	"github.com/VladimirSh98/urlShortener/internal/app/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type dBConnectionStruct struct {
	Conn *sql.DB
}

// DBConnection contains connection to database
var DBConnection = dBConnectionStruct{}

// OpenConnection open database connection
func (db *dBConnectionStruct) OpenConnection() error {
	var err error
	db.Conn, err = sql.Open("pgx", config.DatabaseDSN)
	if err != nil {
		return err
	}
	return nil
}

// CloseConnection close database connection
func (db *dBConnectionStruct) CloseConnection() {
	db.Conn.Close()
}
