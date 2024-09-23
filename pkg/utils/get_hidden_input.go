package utils

import (
	"fmt"
	"golang.org/x/term"
	"os"
	"strings"
)

// GetHiddenInput securely gets user input for passwords
func GetHiddenInput(prompt string) string {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Errorf("Error reading password: " + err.Error())
		return ""
	}
	fmt.Println() // Print a newline after input
	return strings.TrimSpace(string(bytePassword))
}
