package validation

import (
	"cryptotracker/pkg/validation"
	"testing"
)

// TestIsValidPassword tests the IsValidPassword function
func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		password string
		expected bool
	}{
		{"P@ssw0rd", true},
		{"AmanKawadia@1", true},
		{"Pravin@1234", true},
		{"short", false},
		{"NoSpecialChar1", false},
		{"NOLOWER123!", false},
		{"noupper123!", false},
		{"NoNumber!", false},
		{"12345678", false},
		{"ABCD1234", false},
		{"abcd1234", false},
		{"ABCDabcd", false},
		{"ABCD1234@", false},
		{"1234!@#$", false},
	}

	for _, test := range tests {
		result := validation.IsValidPassword(test.password)
		if result != test.expected {
			t.Errorf("IsValidPassword(%q) = %v; expected %v", test.password, result, test.expected)
		}
	}
}
