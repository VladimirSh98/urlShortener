package database

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetByOriginalURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &ShortenRepository{Conn: db}
	expected := Shorter{
		ID:          "abc123",
		OriginalURL: "http://example.com",
		CreatedAt:   time.Now(),
		UserID:      1,
		Archived:    false,
	}

	mock.ExpectQuery("SELECT \\* FROM urls WHERE original_url = 'http://example.com' limit 1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "original_url", "created_at", "user_id", "archived"}).
			AddRow(expected.ID, expected.OriginalURL, expected.CreatedAt, expected.UserID, expected.Archived))

	result, err := repo.GetByOriginalURL("http://example.com")
	require.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
	assert.Equal(t, expected.OriginalURL, result.OriginalURL)
	assert.Equal(t, expected.UserID, result.UserID)
	assert.Equal(t, expected.Archived, result.Archived)
	assert.NoError(t, mock.ExpectationsWereMet())
}
