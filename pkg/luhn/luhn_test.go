package luhn_test

import (
	"testing"

	"github.com/sodiqit/gophermart/pkg/luhn"
)

func TestValidateString(t *testing.T) {
	tests := []struct {
		name           string
		value          string
		expectedResult bool
	}{
		{
			name:           "valid number",
			value:          "4561261212345467",
			expectedResult: true,
		},
		{
			name:           "invalid number",
			value:          "4561261212345464",
			expectedResult: false,
		},
		{
			name:           "valid number",
			value:          "12345678903",
			expectedResult: true,
		},
		{
			name:           "invalid number",
			value:          "12345678904",
			expectedResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := luhn.ValidateString(tt.value); got != tt.expectedResult {
				t.Errorf("ValidateString() = %v, want %v", got, tt.expectedResult)
			}
		})
	}
}
