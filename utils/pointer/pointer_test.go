package pointer_test

import (
	"github.com/project-ippl-dev/tanding-api/utils/pointer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUtils_GetValueFromPointer(t *testing.T) {
	var testValue *string
	var emptyVal string
	var existVal = "exist"

	testCases := []struct {
		description string
		req         *string
		expectedRes string
	}{
		{
			description: "if value with type pointer string is nil",
			req:         testValue,
			expectedRes: emptyVal,
		},
		{
			description: "if value with type pointer string is not nil",
			req:         &existVal,
			expectedRes: existVal,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			res := pointer.GetValueFromPointer(testCase.req)
			assert.Equal(t, testCase.expectedRes, res, "res must be %t", testCase.expectedRes)
		})
	}
}

func TestUtils_ConvertToPointer(t *testing.T) {
	var existVal = "exist"

	testCases := []struct {
		description string
		req         string
		expectedRes *string
	}{
		{
			description: "convert string to pointer",
			req:         existVal,
			expectedRes: &existVal,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			res := pointer.ConvertToPointer(testCase.req)
			assert.Equal(t, testCase.expectedRes, res, "res must be %t", testCase.expectedRes)
		})
	}
}
