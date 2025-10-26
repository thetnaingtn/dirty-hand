package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	// Server configuration
	Server ServerConfig `mapstructure:"server" yaml:"server"`

	// Database configuration
	Database DatabaseConfig `mapstructure:"database" yaml:"database"`

	// Environment
	Environment string `mapstructure:"environment" yaml:"environment"`
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Addr string `mapstructure:"addr" yaml:"addr"`
	Port string `mapstructure:"port" yaml:"port"`
}

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	DSN string `mapstructure:"dsn" yaml:"dsn"`
}

// NewConfig creates a new configuration instance
// It supports multiple configuration sources with the following precedence:
// 1. Command line flags (--config) and environment variables
// 2. Configuration file specified by CONFIG_FILE env var or --config flag
// 3. config.yaml in current directory
// 4. Default values
func NewConfig() (*Config, error) {
	configPath := parseConfigFlag()
	return NewConfigWithPath(configPath)
}

// parseConfigFlag parses the --config command-line flag
func parseConfigFlag() string {
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to configuration file")

	// Only parse if flags haven't been parsed yet
	if !flag.Parsed() {
		flag.Parse()
	}

	return configPath
}

// NewConfigWithPath creates a new configuration instance with a specific config file path
func NewConfigWithPath(configPath string) (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Configure viper
	setupViper(v, configPath)

	// Read configuration
	if err := readConfig(v); err != nil {
		return nil, fmt.Errorf("failed to read configuration: %w", err)
	}

	// Unmarshal configuration
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	return &config, nil
}

// NewConfigWithFlags creates a new configuration instance and allows external flag management
// This is useful when you want to define flags in your main application or use a different flag library
func NewConfigWithFlags(configPath string) (*Config, error) {
	return NewConfigWithPath(configPath)
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.addr", "localhost")
	v.SetDefault("server.port", "8080")

	// Database defaults
	v.SetDefault("database.dsn", "")

	// Environment default
	v.SetDefault("environment", "development")
}

// setupViper configures viper with file paths, environment variables, etc.
func setupViper(v *viper.Viper, configPath string) {
	// Configure file reading
	v.SetConfigType("yaml")
	v.SetConfigName("config")

	// Priority order for config file:
	// 1. --config flag (configPath parameter)
	// 2. CONFIG_FILE environment variable
	// 3. Default search paths

	// Handle --config flag first (highest priority)
	if configPath != "" {
		if filepath.Ext(configPath) != "" {
			// If it's a file path with extension
			v.SetConfigFile(configPath)
		} else {
			// If it's a directory path
			v.AddConfigPath(configPath)
		}
	} else {
		// Check for CONFIG_FILE environment variable (second priority)
		if configFile := os.Getenv("CONFIG_FILE"); configFile != "" {
			v.SetConfigFile(configFile)
		}
	}

	// Add default search paths (lowest priority)
	v.AddConfigPath(".")
	v.AddConfigPath("./internal/config")
	v.AddConfigPath("/etc/github.com/thetnaingtn/dirty-hand/")
	v.AddConfigPath("$HOME/.github.com/thetnaingtn/dirty-hand")

	// Configure environment variable reading
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Bind specific environment variables for backward compatibility
	bindEnvironmentVariables(v)
}

// bindEnvironmentVariables binds specific environment variables for backward compatibility
func bindEnvironmentVariables(v *viper.Viper) {
	// Server configuration
	v.BindEnv("server.addr", "SERVER_ADDR", "ADDR")
	v.BindEnv("server.port", "SERVER_PORT", "PORT")

	// Database configuration
	v.BindEnv("database.dsn", "DATABASE_DSN")

	// Environment
	v.BindEnv("environment", "ENVIRONMENT")
}

// readConfig reads the configuration file
func readConfig(v *viper.Viper) error {
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error since we have defaults
			return nil
		}
		// Config file was found but another error was produced
		return err
	}
	return nil
}

// IsDevelopment returns true if running in development environment
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
