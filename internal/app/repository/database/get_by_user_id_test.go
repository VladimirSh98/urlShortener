package database

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &ShortenRepository{Conn: db}
	expected := []Shorter{
		{
			ID:          "abc123",
			OriginalURL: "http://example.com",
			CreatedAt:   time.Now(),
			UserID:      1,
			Archived:    false,
		},
		{
			ID:          "def456",
			OriginalURL: "http://another.com",
			CreatedAt:   time.Now(),
			UserID:      1,
			Archived:    true,
		},
	}

	rows := sqlmock.NewRows([]string{"id", "original_url", "created_at", "user_id", "archived"}).
		AddRow(expected[0].ID, expected[0].OriginalURL, expected[0].CreatedAt, expected[0].UserID, expected[0].Archived).
		AddRow(expected[1].ID, expected[1].OriginalURL, expected[1].CreatedAt, expected[1].UserID, expected[1].Archived)

	mock.ExpectQuery("SELECT \\* FROM urls WHERE user_id = '1'").
		WillReturnRows(rows)

	result, err := repo.GetByUserID(1)
	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, expected[0].ID, result[0].ID)
	assert.Equal(t, expected[1].OriginalURL, result[1].OriginalURL)
	assert.NoError(t, mock.ExpectationsWereMet())
}
