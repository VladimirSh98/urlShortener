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

// ShortenRepoInterface interface for repository
type ShortenRepoInterface interface {
	BatchCreate(data []ShortenBatchRequest) error
	BatchUpdate(data []string, userID int) error
	Create(mask string, originalURL string, userID int) (sql.Result, error)
	GetAllRecords() ([]Shorter, error)
	GetByOriginalURL(originalURL string) (Shorter, error)
	GetByShortURL(shortURL string) (Shorter, error)
	GetByUserID(userID int) ([]Shorter, error)
	Ping() error
}

// NewShortenRepository create new repository
func NewShortenRepository(conn *sql.DB) ShortenRepoInterface {
	return &ShortenRepository{Conn: conn}
}
