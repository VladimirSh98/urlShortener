package shorten

import (
	repoMock "github.com/VladimirSh98/urlShortener/mocks/repository"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBatchUpdate(t *testing.T) {
	t.Run("success batch update", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repoMock.NewMockShortenRepoInterface(ctrl)
		svc := ShortenService{
			Repo: mockRepo,
		}

		data := []string{"mask1", "mask2"}
		userID := 1

		mockRepo.EXPECT().BatchUpdate(data, userID).Return(nil)

		err := svc.BatchUpdate(data, userID)
		assert.NoError(t, err)
	})

	t.Run("repo returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repoMock.NewMockShortenRepoInterface(ctrl)
		svc := ShortenService{
			Repo: mockRepo,
		}

		data := []string{"mask1", "mask2"}
		userID := 1

		mockRepo.EXPECT().BatchUpdate(data, userID).Return(errors.New("update error"))

		err := svc.BatchUpdate(data, userID)
		assert.EqualError(t, err, "update error")
	})
}
