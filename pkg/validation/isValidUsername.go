package validation

import "regexp"

// IsValidUsername validates the username format
func IsValidUsername(username string) bool {
	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return validUsername.MatchString(username)
}
