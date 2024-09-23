package validation

import (
	"cryptotracker/pkg/validation"
	"testing"
)

// TestIsValidUsername tests the IsValidUsername function
func TestIsValidUsername(t *testing.T) {
	tests := []struct {
		username string
		expected bool
	}{
		{"user123", true},
		{"user_name", true},
		{"UserName", true},
		{"user", true},
		{"username_", true},
		{"_username", true},
		{"u123", true},
		{"user!123", false},
		{"user name", false},
		{"user@name", false},
		{"user-name", false},
		{"user.name", false},
		{"", false},
		{"user name123", false},
	}

	for _, test := range tests {
		result := validation.IsValidUsername(test.username)
		if result != test.expected {
			t.Errorf("IsValidUsername(%q) = %v; expected %v", test.username, result, test.expected)
		}
	}
}
