// internal/config/config_test.go
package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/msaglietto/mantrid/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("default configuration", func(t *testing.T) {
		cfg, err := config.Load()
		require.NoError(t, err)

		// Check default values
		assert.Equal(t, "json", cfg.StorageType)
		assert.Equal(t, "info", cfg.LogLevel)
		assert.Equal(t, "json", cfg.LogFormat)
		assert.False(t, cfg.CloudEnabled)
	})

	t.Run("configuration from file", func(t *testing.T) {
		// Create a temporary directory for the test
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		// Create a test configuration file
		configContent := []byte(`
storage_type: "memory"
log_level: "debug"
log_format: "text"
cloud_enabled: true
cloud_url: "https://example.com"
`)
		err := os.WriteFile(configPath, configContent, 0644)
		require.NoError(t, err)

		// Set environment variable to point to our test config
		os.Setenv("MANTRID_CONFIG", configPath)
		defer os.Unsetenv("MANTRID_CONFIG")

		// Load the configuration
		cfg, err := config.Load()
		require.NoError(t, err)

		// Verify the loaded configuration
		assert.Equal(t, "memory", cfg.StorageType)
		assert.Equal(t, "debug", cfg.LogLevel)
		assert.Equal(t, "text", cfg.LogFormat)
		assert.True(t, cfg.CloudEnabled)
		assert.Equal(t, "https://example.com", cfg.CloudURL)
	})

	t.Run("configuration from environment variables", func(t *testing.T) {
		// Set environment variables
		os.Setenv("MANTRID_LOG_LEVEL", "debug")
		os.Setenv("MANTRID_STORAGE_TYPE", "memory")
		defer func() {
			os.Unsetenv("MANTRID_LOG_LEVEL")
			os.Unsetenv("MANTRID_STORAGE_TYPE")
		}()

		// Load the configuration
		cfg, err := config.Load()
		require.NoError(t, err)

		// Verify environment variables were applied
		assert.Equal(t, "debug", cfg.LogLevel)
		assert.Equal(t, "memory", cfg.StorageType)
	})

	t.Run("invalid configuration", func(t *testing.T) {
		os.Setenv("MANTRID_LOG_LEVEL", "invalid")
		defer os.Unsetenv("MANTRID_LOG_LEVEL")

		_, err := config.Load()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid log level")
	})
}
