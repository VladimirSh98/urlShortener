package database

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBatchCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &ShortenRepository{Conn: db}
	data := []ShortenBatchRequest{
		{Mask: "abc123", URL: "http://example.com", UserID: 1},
		{Mask: "def456", URL: "http://example.org", UserID: 2},
	}
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO urls .*`).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(`INSERT INTO urls .*`).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(2, 1))
	mock.ExpectCommit()
	err = repo.BatchCreate(data)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
