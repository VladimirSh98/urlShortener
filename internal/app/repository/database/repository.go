package database

import (
	"database/sql"
	"time"
)

// ShortenRepository contains database connection
type ShortenRepository struct {
	Conn *sql.DB
}

// Shorter contains info about short URL
type Shorter struct {
	ID          string
	OriginalURL string
	CreatedAt   time.Time
	UserID      int
	Archived    bool
}

// ShortenBatchRequest struct for batch update
type ShortenBatchRequest struct {
	URL    string
	Mask   string
	UserID int
}
