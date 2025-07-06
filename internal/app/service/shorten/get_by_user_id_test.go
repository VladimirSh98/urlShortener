package shorten

import (
	dbrepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
	repoMock "github.com/VladimirSh98/urlShortener/mocks/repository"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetByUserID(t *testing.T) {
	t.Run("success get by user ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repoMock.NewMockShortenRepoInterface(ctrl)
		svc := ShortenService{
			Repo: mockRepo,
		}

		userID := 42
		expected := []dbrepo.Shorter{
			{ID: "mask1", OriginalURL: "http://example.com/1"},
			{ID: "mask2", OriginalURL: "http://example.com/2"},
		}

		mockRepo.EXPECT().GetByUserID(userID).Return(expected, nil)

		result, err := svc.GetByUserID(userID)
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

		userID := 42

		mockRepo.EXPECT().GetByUserID(userID).Return(nil, errors.New("db error"))

		result, err := svc.GetByUserID(userID)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "db error")
	})
}
