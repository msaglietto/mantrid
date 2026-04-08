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
}

// defaultConfig provides default values for all configuration options
var defaultConfig = Config{
	StorageType: "json",
	LogLevel:    "info",
	LogFormat:   "json",
}

// Load reads the configuration from multiple sources in the following order:
// 1. Default values
// 2. Configuration file (config.yaml)
// 3. Environment variables (MANTRID_*)
// 4. Command line flags (--config)
func Load(configFilePath ...string) (*Config, error) {
	// Start with default configuration
	cfg := defaultConfig

	// Initialize Viper
	v := viper.New()

	// Set up Viper to read environment variables
	v.SetEnvPrefix("MANTRID")
	v.AutomaticEnv()

	// Set default values
	v.SetDefault("storage_type", defaultConfig.StorageType)
	v.SetDefault("log_level", defaultConfig.LogLevel)
	v.SetDefault("log_format", defaultConfig.LogFormat)

	// Set up the default alias file path
	defaultAliasFile := filepath.Join(getConfigDir(), "aliases.json")
	v.SetDefault("alias_file", defaultAliasFile)

	// Determine config file path: explicit > env var > default
	configFile := ""
	if len(configFilePath) > 0 && configFilePath[0] != "" {
		configFile = configFilePath[0]
	} else if envConfig := v.GetString("config"); envConfig != "" {
		configFile = envConfig
	}

	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(getConfigDir())
		v.AddConfigPath(".")
	}

	// Try to read the config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal the configuration into our Config struct
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate the configuration
	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
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

// ExampleConfig returns an example configuration as a YAML string
func ExampleConfig() string {
	return `# Storage configuration
alias_file: "~/.mantrid/aliases.json"
storage_type: "json"

# Logging configuration
log_level: "info"
log_format: "json"
`
}
