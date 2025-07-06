package shorten

import (
	dbrepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
	repoMock "github.com/VladimirSh98/urlShortener/mocks/repository"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetByOriginalURL(t *testing.T) {
	t.Run("success get by original URL", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repoMock.NewMockShortenRepoInterface(ctrl)
		svc := ShortenService{
			Repo: mockRepo,
		}

		originalURL := "http://example.com"
		expected := dbrepo.Shorter{ID: "mask123", OriginalURL: originalURL}

		mockRepo.EXPECT().GetByOriginalURL(originalURL).Return(expected, nil)

		result, err := svc.GetByOriginalURL(originalURL)
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

		originalURL := "http://example.com"

		mockRepo.EXPECT().GetByOriginalURL(originalURL).Return(dbrepo.Shorter{}, errors.New("not found"))

		result, err := svc.GetByOriginalURL(originalURL)
		assert.Error(t, err)
		assert.Equal(t, dbrepo.Shorter{}, result)
		assert.EqualError(t, err, "not found")
	})
}
