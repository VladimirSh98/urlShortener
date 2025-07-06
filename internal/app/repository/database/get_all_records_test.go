package database

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetAllRecords(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &ShortenRepository{Conn: db}

	rows := sqlmock.NewRows([]string{"id", "original_url", "created_at", "user_id", "archived"}).
		AddRow("abc123", "http://example.com", time.Now(), 1, false).
		AddRow("def456", "http://another.com", time.Now(), 2, true)

	mock.ExpectQuery("SELECT \\* FROM urls").
		WillReturnRows(rows)

	records, err := repo.GetAllRecords()
	require.NoError(t, err)
	require.Len(t, records, 2)
	assert.Equal(t, "abc123", records[0].ID)
	assert.Equal(t, "http://example.com", records[0].OriginalURL)
	assert.Equal(t, 1, records[0].UserID)
	assert.Equal(t, false, records[0].Archived)

	assert.Equal(t, "def456", records[1].ID)
	assert.Equal(t, "http://another.com", records[1].OriginalURL)
	assert.Equal(t, 2, records[1].UserID)
	assert.Equal(t, true, records[1].Archived)

	assert.NoError(t, mock.ExpectationsWereMet())
}
