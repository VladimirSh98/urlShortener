package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthorization(t *testing.T) {
	type expect struct {
		length int
	}
	tests := []struct {
		description string
		expect      expect
	}{
		{
			description: "Test #1. Success",
			expect: expect{
				length: 8,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			mask := CreateRandomMask()
			assert.Equal(t, len(mask), test.expect.length)
		})
	}
}
