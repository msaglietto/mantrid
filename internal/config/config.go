// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds all configuration values for the application
type Config struct {
	// Storage configuration
	AliasFile   string `mapstructure:"alias_file"`
	StorageType string `mapstructure:"storage_type"`

	// Logging configuration
	LogLevel  string `mapstructure:"log_level"`
	LogFormat string `mapstructure:"log_format"`

	// Cloud sync configuration (for future use)
	CloudEnabled bool   `mapstructure:"cloud_enabled"`
	CloudURL     string `mapstructure:"cloud_url"`
}

// defaultConfig provides default values for all configuration options
var defaultConfig = Config{
	StorageType:  "json",
	LogLevel:     "info",
	LogFormat:    "json",
	CloudEnabled: false,
}

// Load reads the configuration from multiple sources in the following order:
// 1. Default values
// 2. Configuration file (config.yaml)
// 3. Environment variables (MANTRID_*)
// 4. Command line flags (not implemented in this example)
func Load() (*Config, error) {
	// Start with default configuration
	config := defaultConfig

	// Initialize Viper
	v := viper.New()

	// Set up Viper to read environment variables
	v.SetEnvPrefix("MANTRID") // Environment variables will be prefixed with MANTRID_
	v.AutomaticEnv()          // Automatically read environment variables

	// Set default values
	v.SetDefault("storage_type", defaultConfig.StorageType)
	v.SetDefault("log_level", defaultConfig.LogLevel)
	v.SetDefault("log_format", defaultConfig.LogFormat)
	v.SetDefault("cloud_enabled", defaultConfig.CloudEnabled)

	// Set up the default alias file path
	defaultAliasFile := filepath.Join(getConfigDir(), "aliases.json")
	v.SetDefault("alias_file", defaultAliasFile)

	// Check if a specific config file is specified via environment variable
	if configFile := v.GetString("config"); configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		// Configure Viper to read the config file
		v.SetConfigName("config") // Name of config file (without extension)
		v.SetConfigType("yaml")   // Config file type

		// Add paths where Viper should look for the config file
		v.AddConfigPath(getConfigDir()) // First check in .mantrid directory
		v.AddConfigPath(".")            // Then check current directory
	}

	// Try to read the config file
	if err := v.ReadInConfig(); err != nil {
		// It's okay if we can't find a config file, but other errors should be reported
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal the configuration into our Config struct
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate the configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// getConfigDir returns the path to the configuration directory
func getConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(homeDir, ".mantrid")
}

// validateConfig checks if the configuration is valid
func validateConfig(cfg *Config) error {
	// Validate log level
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[cfg.LogLevel] {
		return fmt.Errorf("invalid log level: %s", cfg.LogLevel)
	}

	// Validate storage type
	validStorageTypes := map[string]bool{
		"json":   true,
		"memory": true,
	}
	if !validStorageTypes[cfg.StorageType] {
		return fmt.Errorf("invalid storage type: %s", cfg.StorageType)
	}

	return nil
}

// Example config.yaml file
func ExampleConfig() string {
	return `
# Storage configuration
alias_file: "~/.mantrid/aliases.json"
storage_type: "json"

# Logging configuration
log_level: "info"
log_format: "json"

# Cloud sync configuration
cloud_enabled: false
cloud_url: ""
`
}
