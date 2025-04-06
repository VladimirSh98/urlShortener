package database

import (
	"database/sql"
)

type DBConnectionStruct struct {
	conn *sql.DB
}
