package ui

import (
	"bytes"
	"cryptotracker/pkg/ui"

	//"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/fatih/color"
)

// Helper function to capture output of a function that prints to stdout
func captureOutput(f func()) string {
	// Create a temporary file to capture output
	tempFile, err := os.CreateTemp("", "test_output_*.txt")
	if err != nil {
		panic("Failed to create temp file: " + err.Error())
	}
	defer os.Remove(tempFile.Name())

	original := os.Stdout
	defer func() { os.Stdout = original }()

	os.Stdout = tempFile

	f()

	// Seek to the beginning of the file to read from it
	tempFile.Seek(0, 0)

	var buf bytes.Buffer
	_, err = buf.ReadFrom(tempFile)
	if err != nil {
		panic("Failed to read from temp file: " + err.Error())
	}

	return buf.String()
}

func TestClearScreen(t *testing.T) {
	// Skip on platforms where `clear` is not available or fails
	cmd := exec.Command("clear")
	if err := cmd.Run(); err != nil {
		t.Skip("Skipping ClearScreen test due to error: ", err)
	}
	// No specific output to check, just ensure the command runs without errors.
}

func TestDisplayWelcomeBanner(t *testing.T) {
	expectedOutput := color.CyanString("==============================\n" +
		"  Welcome to CryptoTracker! 🚀\n" +
		"==============================\n\n")
	output := captureOutput(ui.DisplayWelcomeBanner)
	if output != expectedOutput {
		t.Errorf("DisplayWelcomeBanner() = %q; expected %q", output, expectedOutput)
	}
}

func TestDisplayAuthMenu(t *testing.T) {
	expectedOutput := color.BlueString("\n\n==================================\n") +
		color.YellowString("        Authentication Menu       \n") +
		color.BlueString("==================================\n\n\n") +
		color.BlueString("1. Login\n") +
		color.BlueString("2. SignUp\n") +
		color.BlueString("3. Exit\n") +
		"\n"
	output := captureOutput(ui.DisplayAuthMenu)
	if output != expectedOutput {
		t.Errorf("DisplayAuthMenu() = %q; expected %q", output, expectedOutput)
	}
}

func TestDisplayMainMenu(t *testing.T) {
	expectedOutput := color.BlueString("\n==================================\n") +
		color.YellowString("             Main Menu            \n") +
		color.BlueString("==================================\n\n") +
		color.GreenString("1. View Top 10 Cryptocurrencies\n") +
		color.GreenString("2. Search for a Cryptocurrency\n") +
		color.GreenString("3. Set Price Alert\n") +
		//color.GreenString("4. Check if user is Admin\n") +
		color.GreenString("4. User Profile\n") +
		color.RedString("5. Logout\n") +
		"\n"
	output := captureOutput(ui.DisplayMainMenu)
	if output != expectedOutput {
		t.Errorf("DisplayMainMenu() = %q; expected %q", output, expectedOutput)
	}
}

func TestPrintError(t *testing.T) {
	expectedOutput := color.RedString("Sample error message\n")
	output := captureOutput(func() {
		ui.PrintError("Sample error message")
	})
	if output != expectedOutput {
		t.Errorf("PrintError() = %q; expected %q", output, expectedOutput)
	}
}
