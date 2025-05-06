package logger

import (
	"fmt"
	"os"
)

type Config struct {
	LogType   string `yaml:"type"`   // "console" or "file"
	LogDir    string `yaml:"dir"`    // Directory for log files when LogType is "file"
	LogLevel  string `yaml:"level"`  // "debug", "info", "warn", "error"
	LogFormat string `yaml:"format"` // "text" or "json"
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		LogType:   "console",
		LogDir:    "./logs",
		LogLevel:  "info",
		LogFormat: "text",
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate LogType
	if c.LogType != "console" && c.LogType != "file" {
		return fmt.Errorf("invalid log type: %s", c.LogType)
	}

	// Validate LogLevel
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLevels[c.LogLevel] {
		return fmt.Errorf("invalid log level: %s", c.LogLevel)
	}

	// Validate LogFormat
	if c.LogFormat != "text" && c.LogFormat != "json" {
		return fmt.Errorf("invalid log format: %s", c.LogFormat)
	}

	// Validate LogDir if using file logging
	if c.LogType == "file" {
		if c.LogDir == "" {
			return fmt.Errorf("log directory is required for file logging")
		}
		// Create log directory if it doesn't exist
		if err := os.MkdirAll(c.LogDir, 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %v", err)
		}
	}
	return nil
}
