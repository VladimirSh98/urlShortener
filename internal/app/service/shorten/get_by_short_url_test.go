package shorten

import (
	dbrepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
	repoMock "github.com/VladimirSh98/urlShortener/mocks/repository"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetByShortURL(t *testing.T) {
	t.Run("success get by short URL", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repoMock.NewMockShortenRepoInterface(ctrl)
		svc := ShortenService{
			Repo: mockRepo,
		}

		shortURL := "mask123"
		expected := dbrepo.Shorter{ID: shortURL, OriginalURL: "http://example.com"}

		mockRepo.EXPECT().GetByShortURL(shortURL).Return(expected, nil)

		result, err := svc.GetByShortURL(shortURL)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("repo returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repoMock.NewMockShortenRepoInterface(ctrl)
		svc := ShortenService{
			Repo: mockRepo,
		}

		shortURL := "mask123"

		mockRepo.EXPECT().GetByShortURL(shortURL).Return(dbrepo.Shorter{}, errors.New("not found"))

		result, err := svc.GetByShortURL(shortURL)
		assert.Error(t, err)
		assert.Equal(t, dbrepo.Shorter{}, result)
		assert.EqualError(t, err, "not found")
	})
}
