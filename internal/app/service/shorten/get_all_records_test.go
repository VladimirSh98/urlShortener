package shorten

import (
	dbrepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
	repoMock "github.com/VladimirSh98/urlShortener/mocks/repository"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAllRecords(t *testing.T) {
	t.Run("success get all records", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repoMock.NewMockShortenRepoInterface(ctrl)
		svc := ShortenService{
			Repo: mockRepo,
		}

		expectedRecords := []dbrepo.Shorter{
			{ID: "mask1", OriginalURL: "http://example.com/1"},
			{ID: "mask2", OriginalURL: "http://example.com/2"},
		}

		mockRepo.EXPECT().GetAllRecords().Return(expectedRecords, nil)

		records, err := svc.GetAllRecords()
		assert.NoError(t, err)
		assert.Equal(t, expectedRecords, records)
	})

	t.Run("error from repo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repoMock.NewMockShortenRepoInterface(ctrl)
		svc := ShortenService{
			Repo: mockRepo,
		}

		mockRepo.EXPECT().GetAllRecords().Return(nil, errors.New("db failure"))

		records, err := svc.GetAllRecords()
		assert.Nil(t, records)
		assert.EqualError(t, err, "db failure")
	})
}
