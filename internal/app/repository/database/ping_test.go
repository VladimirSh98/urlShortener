package database

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPing(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPing()

	repo := &ShortenRepository{Conn: db}

	err = repo.Ping()
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
