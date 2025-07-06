package database

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOpenConnection(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	conn := &DBConnectionStruct{
		Conn: mockDB,
	}

	err = conn.OpenConnection()
	assert.NoError(t, err)
}

func TestCloseConnection(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	assert.NoError(t, err)

	conn := &DBConnectionStruct{
		Conn: mockDB,
	}

	conn.CloseConnection()

	assert.NoError(t, err)
}
