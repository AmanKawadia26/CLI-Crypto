package validation

import "fmt"

// IsValidMobile validates the mobile number format
func IsValidMobile(mobile int) bool {
	return len(fmt.Sprintf("%d", mobile)) == 10
}
