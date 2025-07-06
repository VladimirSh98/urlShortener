package shorten

import (
	dbrepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
	repoMock "github.com/VladimirSh98/urlShortener/mocks/repository"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"testing"
	"time"
)

func TestBatchCreate(t *testing.T) {
	t.Run("success batch create", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repoMock.NewMockShortenRepoInterface(ctrl)
		svc := ShortenService{
			Repo: mockRepo,
		}

		data := []dbrepo.ShortenBatchRequest{
			{Mask: "mask1", URL: "http://example.com/1"},
			{Mask: "mask2", URL: "http://example.com/2"},
		}

		mockRepo.EXPECT().BatchCreate(data).Return(nil)

		svc.BatchCreate(data)
	})

	t.Run("repo returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repoMock.NewMockShortenRepoInterface(ctrl)
		svc := ShortenService{
			Repo: mockRepo,
		}

		data := []dbrepo.ShortenBatchRequest{
			{Mask: "mask1", URL: "http://example.com/1"},
			{Mask: "mask2", URL: "http://example.com/2"},
		}

		mockRepo.EXPECT().BatchCreate(data).Return(errors.New("db error"))

		svc.BatchCreate(data)

		time.Sleep(100 * time.Millisecond)
	})
}
