package database

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortenRepositoryGetByShortURL(t *testing.T) {
	tests := []struct {
		name          string
		shortURL      string
		mockClosure   func(mock sqlmock.Sqlmock)
		expected      Shorter
		expectedError error
	}{
		{
			name:     "успешное получение записи",
			shortURL: "abc123",
			mockClosure: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "original_url", "created_at", "user_id", "archived"}).
					AddRow("abc123", "http://example.com", time.Now(), 1, false)
				mock.ExpectQuery(`SELECT \* FROM urls WHERE id = 'abc123' limit 1`).
					WillReturnRows(rows)
			},
			expected: Shorter{
				ID:          "abc123",
				OriginalURL: "http://example.com",
				UserID:      1,
				Archived:    false,
			},
			expectedError: nil,
		},
		{
			name:     "запись не найдена",
			shortURL: "notfound",
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM urls WHERE id = 'notfound' limit 1`).
					WillReturnError(sql.ErrNoRows)
			},
			expected:      Shorter{},
			expectedError: sql.ErrNoRows,
		},
		{
			name:     "ошибка базы данных",
			shortURL: "dberror",
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM urls WHERE id = 'dberror' limit 1`).
					WillReturnError(errors.New("database error"))
			},
			expected:      Shorter{},
			expectedError: errors.New("database error"),
		},
		{
			name:          "нет соединения с БД",
			shortURL:      "any",
			mockClosure:   func(mock sqlmock.Sqlmock) {},
			expected:      Shorter{},
			expectedError: errors.New("DB connection is not open"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			repo := &ShortenRepository{Conn: db}
			if tt.name == "нет соединения с БД" {
				repo.Conn = nil
			}
			tt.mockClosure(mock)
			result, err := repo.GetByShortURL(tt.shortURL)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.OriginalURL, result.OriginalURL)
			assert.Equal(t, tt.expected.UserID, result.UserID)
			assert.Equal(t, tt.expected.Archived, result.Archived)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
