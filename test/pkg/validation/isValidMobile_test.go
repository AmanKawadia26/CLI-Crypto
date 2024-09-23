package validation

import (
	"cryptotracker/pkg/validation"
	"testing"
)

// TestIsValidMobile tests the IsValidMobile function
func TestIsValidMobile(t *testing.T) {
	tests := []struct {
		mobile   int
		expected bool
	}{
		{1234567890, true},
		{9876543210, true},
		{123456789, false},
		{12345678901, false},
		{12345678, false},
		{12345678901234, false},
		{0, false},
		{-1234567890, false},
	}

	for _, test := range tests {
		result := validation.IsValidMobile(test.mobile)
		if result != test.expected {
			t.Errorf("IsValidMobile(%d) = %v; expected %v", test.mobile, result, test.expected)
		}
	}
}
