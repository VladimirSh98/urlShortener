package database

import (
	"database/sql"
	"time"
)

type ShortenRepository struct {
	Conn *sql.DB
}

type Shorter struct {
	ID          string
	OriginalURL string
	CreatedAt   time.Time
	UserID      int
	Archived    bool
}

type ShortenBatchRequest struct {
	URL    string
	Mask   string
	UserID int
}
