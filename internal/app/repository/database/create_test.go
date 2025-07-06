package database

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &ShortenRepository{Conn: db}

	mask := "abc123"
	originalURL := "http://example.com"
	userID := 42

	query := fmt.Sprintf(
		"INSERT INTO urls \\(id, original_url, user_id, archived\\) VALUES \\('%s', '%s', '%d', false\\);",
		mask, originalURL, userID,
	)
	mock.ExpectExec(query).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = repo.Create(mask, originalURL, userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
