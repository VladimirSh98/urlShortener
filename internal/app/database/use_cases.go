package database

import "database/sql"

func (db *DBConnectionStruct) Ping() error {
	err := db.conn.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (db *DBConnectionStruct) Exec(query string) (sql.Result, error) {
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
