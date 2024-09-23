package validation

import (
	"cryptotracker/pkg/validation"
	"testing"
)

// TestIsValidEmail tests the IsValidEmail function
func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"test@example.com", true},
		{"user.name+tag+sorting@example.com", true},
		{"user@sub.example.com", true},
		{"user@ex.co", true},
		{"user@domain.co.in", true},
		{"plainaddress", false},
		{"@missingusername.com", false},
		{"username@.com", false},
		{"username@com", false},
		{"username@com.", false},
		{"username@.com.", false},
	}

	for _, test := range tests {
		result := validation.IsValidEmail(test.email)
		if result != test.expected {
			t.Errorf("IsValidEmail(%q) = %v; expected %v", test.email, result, test.expected)
		}
	}
}
