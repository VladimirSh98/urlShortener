package shorten

import (
	customErr "github.com/VladimirSh98/urlShortener/internal/app/errors"
	"github.com/VladimirSh98/urlShortener/internal/app/repository/database"
	repoMock "github.com/VladimirSh98/urlShortener/mocks/repository"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreate(t *testing.T) {

	const (
		mask        = "mask123"
		originalURL = "http://example.com"
		userID      = 1
	)

	t.Run("success create", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repoMock.NewMockShortenRepoInterface(ctrl)
		svc := ShortenService{
			Repo: mockRepo,
		}
		mockRepo.EXPECT().Create(mask, originalURL, userID).Return(nil, nil).AnyTimes()

		maskResult, err := svc.Create(mask, originalURL, userID)
		assert.NoError(t, err)
		assert.Equal(t, mask, maskResult)
	})

	t.Run("duplicate URL constraint violation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repoMock.NewMockShortenRepoInterface(ctrl)
		svc := ShortenService{
			Repo: mockRepo,
		}
		pgErr := &pgconn.PgError{Code: pgerrcode.UniqueViolation}

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, pgErr).AnyTimes()
		mockRepo.EXPECT().GetByOriginalURL(gomock.Any()).Return(database.Shorter{ID: "oldmask"}, nil).AnyTimes()

		maskResult, err := svc.Create(mask, originalURL, userID)
		assert.Equal(t, "oldmask", maskResult)
		assert.Equal(t, customErr.ErrConstraintViolation, err)
	})

	t.Run("other error on create", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repoMock.NewMockShortenRepoInterface(ctrl)
		svc := ShortenService{
			Repo: mockRepo,
		}
		mockRepo.EXPECT().Create(mask, originalURL, userID).Return(nil, errors.New("db error")).AnyTimes()

		maskResult, err := svc.Create(mask, originalURL, userID)
		assert.NoError(t, err)
		assert.Equal(t, mask, maskResult)
	})
}
