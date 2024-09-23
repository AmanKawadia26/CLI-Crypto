package utils

import (
	"cryptotracker/pkg/utils"
	"time"

	//"bytes"
	"os"
	"testing"
)

func TestGetHiddenInput(t *testing.T) {
	// Backup the original stdin
	originalStdin := os.Stdin

	// Create a pipe for testing stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	defer r.Close()
	defer w.Close()

	// Replace os.Stdin with our pipe
	os.Stdin = r

	// Write test input to the pipe in a separate goroutine
	testInput := ""
	go func() {
		defer w.Close() // Ensure the pipe is closed
		_, _ = w.Write([]byte(testInput + "\n"))
	}()

	// Wait a bit to ensure input is written
	time.Sleep(100 * time.Millisecond)

	// Test the GetHiddenInput function
	result := utils.GetHiddenInput("Enter password: ")

	// Restore the original stdin
	os.Stdin = originalStdin

	// Check if the result matches the test input
	if result != testInput {
		t.Errorf("GetHiddenInput() = %v, want %v", result, testInput)
	}
}

func TestGetHiddenInput_ErrorHandling(t *testing.T) {
	// Backup the original stdin
	originalStdin := os.Stdin

	// Create a pipe for testing stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	defer r.Close()
	defer w.Close()

	// Replace os.Stdin with our pipe
	os.Stdin = r

	// Simulate an error by closing the write end of the pipe
	if err := w.Close(); err != nil {
		t.Fatalf("Failed to close pipe writer: %v", err)
	}

	// Test the GetHiddenInput function
	result := utils.GetHiddenInput("Enter password: ")

	// Restore the original stdin
	os.Stdin = originalStdin

	// The result should be an empty string because of the error
	if result != "" {
		t.Errorf("GetHiddenInput() = %v, want %v", result, "")
	}
}
