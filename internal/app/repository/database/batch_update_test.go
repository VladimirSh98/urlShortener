package database

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBatchUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &ShortenRepository{Conn: db}

	data := []string{"abc123", "def456"}
	userID := 42

	mock.ExpectExec(`UPDATE urls SET archived = true WHERE id = ANY\(\$1::text\[\]\) and user_id = \$2`).
		WithArgs(pq.Array(data), userID).
		WillReturnResult(sqlmock.NewResult(0, int64(len(data))))

	err = repo.BatchUpdate(data, userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
