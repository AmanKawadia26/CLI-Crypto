package validation

import "regexp"

// IsValidEmail validates the email format
func IsValidEmail(email string) bool {
	validEmail := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return validEmail.MatchString(email)
}
