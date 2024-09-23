package validation

import "unicode"

// IsValidPassword validates the password strength
func IsValidPassword(password string) bool {
	var (
		hasMinLen      = len(password) >= 8
		hasUpper       = false
		hasLower       = false
		hasNumber      = false
		hasSpecialChar = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecialChar = true
		}
	}

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecialChar
}
