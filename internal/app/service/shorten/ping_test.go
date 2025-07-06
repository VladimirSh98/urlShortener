package shorten

import (
	repoMock "github.com/VladimirSh98/urlShortener/mocks/repository"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMock.NewMockShortenRepoInterface(ctrl)
	service := ShortenService{Repo: mockRepo}

	t.Run("successful ping", func(t *testing.T) {
		mockRepo.EXPECT().Ping().Return(nil).Times(1)

		err := service.Ping()
		assert.NoError(t, err)
	})

	t.Run("ping with error", func(t *testing.T) {
		expectedErr := errors.New("database connection error")
		mockRepo.EXPECT().Ping().Return(expectedErr).Times(1)
		err := service.Ping()
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}
