package middleware

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthorization(t *testing.T) {
	type expect struct {
		token  string
		err    error
		userID int
	}
	type testRequest struct {
		addAuth       bool
		authorization string
	}
	tests := []struct {
		description string
		expect      expect
		testRequest testRequest
	}{
		{
			description: "Test #1. No cookie",
			expect: expect{
				err:    nil,
				token:  "",
				userID: 0,
			},
			testRequest: testRequest{
				addAuth: false,
			},
		},
		{
			description: "Test #2. Add cookie",
			expect: expect{
				err:    nil,
				token:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjE2fQ.GTOVjHe7U0vE4KJAQXf-ilZJj-nnYTdRM4FopYizpLw",
				userID: 16,
			},
			testRequest: testRequest{
				addAuth:       true,
				authorization: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjE2fQ.GTOVjHe7U0vE4KJAQXf-ilZJj-nnYTdRM4FopYizpLw",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", nil)
			if test.testRequest.addAuth {
				request.AddCookie(
					&http.Cookie{
						Name:  "Authorization",
						Value: test.testRequest.authorization,
					},
				)
			}
			token, userID, err := authorize(request)
			assert.Equal(t, test.expect.err, err)
			if test.testRequest.addAuth {
				assert.Equal(t, test.expect.token, token)
				assert.Equal(t, test.expect.userID, userID)
			}
		})
	}
}
