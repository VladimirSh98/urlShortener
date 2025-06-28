package service

import (
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	customErr "github.com/VladimirSh98/urlShortener/internal/app/errors"
	dbrepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
	"github.com/VladimirSh98/urlShortener/internal/app/repository/memory"
	"github.com/golang/mock/gomock"
	"testing"

	shortenMock "github.com/VladimirSh98/urlShortener/mocks/shorten"
	"github.com/stretchr/testify/assert"
)

func TestPrefillDataFromDB(t *testing.T) {
	type expect struct {
		err error
	}
	type testRequest struct {
		results          []dbrepo.Shorter
		errGetAllRecords error
	}
	tests := []struct {
		description string
		expect      expect
		testRequest testRequest
	}{
		{
			description: "Test #1. Database error",
			expect: expect{
				err: customErr.ErrConstraintViolation,
			},
			testRequest: testRequest{
				errGetAllRecords: customErr.ErrConstraintViolation,
			},
		},
		{
			description: "Test #1. Success",
			expect: expect{
				err: nil,
			},
			testRequest: testRequest{
				results: []dbrepo.Shorter{
					{
						ID:          "ffsdafd",
						OriginalURL: "http://test.com",
						UserID:      1,
						Archived:    false,
					},
				},
				errGetAllRecords: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := shortenMock.NewMockShortenServiceInterface(ctrl)
			mockService.EXPECT().GetAllRecords().Return(
				test.testRequest.results, test.testRequest.errGetAllRecords).AnyTimes()
			err := prefillFromDB(mockService)
			assert.Equal(t, err, test.expect.err)
		})
	}
}

func TestPrefillDataFromFile(t *testing.T) {
	t.Run("upload data", func(t *testing.T) {
		config.DBFilePath = "test_data/test.json"
		err := prefillFromFile()
		assert.NoError(t, err)
		result, ok := memory.Get("ETe5ORyc")
		assert.True(t, ok)
		assert.Equal(t, result, "http://test.com")
	})
}
