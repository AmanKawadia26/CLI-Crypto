package utils

import (
	"cryptotracker/pkg/utils"
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		password string
		expected string
	}{
		{"admin_password", "6d4525c2a21f9be1cca9e41f3aa402e0765ee5fcc3e7fea34a169b1730ae386e"},
		{"Rishabh@123", "f6d8b1203495a937b766051690a69a3176ceb7f5937c885672c4b3779e1ac07a"},
	}

	for _, test := range tests {
		result := utils.HashPassword(test.password)
		if result != test.expected {
			t.Errorf("HashPassword(%q) = %q; expected %q", test.password, result, test.expected)
		}
	}
}
