package ui

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"os/exec"
)

// Define color variables using Fatih's color package
var (
	colorRed    = color.New(color.FgRed).SprintFunc()
	colorGreen  = color.New(color.FgGreen).SprintFunc()
	colorYellow = color.New(color.FgYellow).SprintFunc()
	colorBlue   = color.New(color.FgBlue).SprintFunc()
	colorPurple = color.New(color.FgMagenta).SprintFunc()
	colorCyan   = color.New(color.FgCyan).SprintFunc()
)

// ClearScreen clears the terminal screen
func ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// DisplayWelcomeBanner shows the welcome message
func DisplayWelcomeBanner() {
	ClearScreen()
	fmt.Println(colorCyan("=============================="))
	fmt.Println(colorCyan("  Welcome to CryptoTracker! 🚀"))
	fmt.Println(colorCyan("=============================="))
	fmt.Println()
}

// DisplayAuthMenu shows the authentication menu
func DisplayAuthMenu() {
	fmt.Println()
	fmt.Println()
	fmt.Println(colorBlue("=================================="))
	fmt.Println(colorYellow("        Authentication Menu       "))
	fmt.Println(colorBlue("=================================="))
	fmt.Println()
	fmt.Println()
	//fmt.Println(colorYellow("Authentication Menu:"))
	fmt.Println(colorBlue("1. Login"))
	fmt.Println(colorBlue("2. SignUp"))
	fmt.Println(colorBlue("3. Exit"))
	fmt.Println()
}

// DisplayMainMenu shows the main menu for regular users
func DisplayMainMenu() {
	fmt.Println()
	fmt.Println(colorBlue("=================================="))
	fmt.Println(colorYellow("             Main Menu            "))
	fmt.Println(colorBlue("=================================="))
	fmt.Println()
	fmt.Println(colorGreen("1. View Top 10 Cryptocurrencies"))
	fmt.Println(colorGreen("2. Search for a Cryptocurrency"))
	fmt.Println(colorGreen("3. Set Price Alert"))
	fmt.Println(colorGreen("4. User Profile"))
	fmt.Println(colorRed("5. Logout"))
	fmt.Println()
}

// PrintError prints an error message in red
func PrintError(message string) {
	fmt.Println(colorRed(message))
}
