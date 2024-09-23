package config

import (
	"cryptotracker/pkg/config"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoadConfig_Success tests if LoadConfig loads a valid configuration successfully.
func TestLoadConfig_Success(t *testing.T) {
	// Create a temporary configuration file
	tempFile, err := os.CreateTemp("", "config_test_*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name()) // Clean up after the test

	// Write sample configuration to the temporary file
	configData := config.Config{
		APIKey: "dummy_api_key",
	}
	jsonData, err := json.Marshal(configData)
	require.NoError(t, err)

	_, err = tempFile.Write(jsonData)
	require.NoError(t, err)

	// Close the file to flush writes
	err = tempFile.Close()
	require.NoError(t, err)

	// Backup original config file name and restore it after test
	originalConfigPath := "config.json"
	defer func() { _ = os.Rename(tempFile.Name(), originalConfigPath) }()

	// Rename the temp file to "config.json" to simulate loading from the expected path
	err = os.Rename(tempFile.Name(), originalConfigPath)
	require.NoError(t, err)

	// Call LoadConfig to load the configuration
	err = config.LoadConfig()
	assert.NoError(t, err)

	// Validate the loaded configuration
	assert.Equal(t, "dummy_api_key", config.AppConfig.APIKey, "Expected API key to be loaded from config file")
}

// TestLoadConfig_FileNotFound tests if LoadConfig handles missing config files correctly.
func TestLoadConfig_FileNotFound(t *testing.T) {
	// Ensure the configuration file does not exist
	_ = os.Remove("config.json")

	// Call LoadConfig, which should return an error due to missing config file
	err := config.LoadConfig()

	// Check for the expected error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "configuration file not found")
}

// TestLoadConfig_InvalidJSON tests if LoadConfig handles invalid JSON correctly.
func TestLoadConfig_InvalidJSON(t *testing.T) {
	// Create a temporary invalid JSON file
	tempFile, err := os.CreateTemp("", "invalid_config_*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name()) // Clean up after the test

	_, err = tempFile.Write([]byte("{invalid json"))
	require.NoError(t, err)

	// Close the file to flush writes
	err = tempFile.Close()
	require.NoError(t, err)

	// Backup original config file name and restore it after test
	originalConfigPath := "config.json"
	defer func() { _ = os.Rename(tempFile.Name(), originalConfigPath) }()

	// Rename the temp file to "config.json" to simulate loading from the expected path
	err = os.Rename(tempFile.Name(), originalConfigPath)
	require.NoError(t, err)

	// Call LoadConfig, which should return an error due to invalid JSON
	err = config.LoadConfig()

	// Check for the expected error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error decoding configuration file")
}
