package database

import (
	"database/sql"
	"errors"
)

func (db *DBConnectionStruct) Ping() error {
	err := db.conn.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (db *DBConnectionStruct) Exec(query string) (sql.Result, error) {
	if db.conn == nil {
		return nil, errors.New("Database connection is not initialized")
	}
	res, err := db.conn.Exec(query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (db *DBConnectionStruct) Query(query string) (*sql.Rows, error) {
	res, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (db *DBConnectionStruct) BatchCreate(queries []string) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}
	for _, query := range queries {
		_, err = tx.Exec(query)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
