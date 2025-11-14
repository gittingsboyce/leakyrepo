package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the LeakyRepo configuration structure
type Config struct {
	// Rules defines regex-based detection rules
	Rules []Rule `yaml:"rules"`
	// EntropyThreshold is the minimum entropy value for high-entropy string detection
	EntropyThreshold float64 `yaml:"entropy_threshold,omitempty"`
	// Allowlist contains patterns that should be ignored
	Allowlist Allowlist `yaml:"allowlist,omitempty"`
}

// Rule defines a regex pattern for secret detection
type Rule struct {
	// ID is a unique identifier for the rule
	ID string `yaml:"id"`
	// Description explains what the rule detects
	Description string `yaml:"description"`
	// Severity indicates the severity level (low, medium, high, critical)
	Severity string `yaml:"severity"`
	// Pattern is the regex pattern to match
	Pattern string `yaml:"pattern"`
	// FileTypes specifies which file extensions this rule applies to (empty means all files)
	FileTypes []string `yaml:"file_types,omitempty"`
}

// Allowlist contains patterns that should be ignored
type Allowlist struct {
	// Files specifies file patterns to ignore (supports glob patterns)
	Files []string `yaml:"files,omitempty"`
	// Strings specifies string patterns to ignore (exact matches)
	Strings []string `yaml:"strings,omitempty"`
}

// DefaultConfig returns a default configuration with common secret detection rules
func DefaultConfig() *Config {
	return &Config{
		EntropyThreshold: 5.5, // Increased from 4.5 to reduce false positives
		Rules: []Rule{
			{
				ID:          "aws_access_key",
				Description: "AWS Access Key",
				Severity:    "high",
				Pattern:     `AKIA[0-9A-Z]{16}`,
				FileTypes:   []string{".env", ".yaml", ".yml", ".json", ".py", ".js", ".ts", ".go"},
			},
			{
				ID:          "generic_api_key",
				Description: "Generic API Key pattern",
				Severity:    "medium",
				Pattern:     `(?i)(api[_-]?key|apikey)\s*[:=]\s*['"]?([a-zA-Z0-9_\-]{20,})['"]?`,
				FileTypes:   []string{".env", ".yaml", ".yml", ".json", ".py", ".js", ".ts", ".go"},
			},
		},
		Allowlist: Allowlist{
			Files: []string{
				".leakyrepoignore",
				".git/**",
				"*.png", "*.jpg", "*.jpeg", "*.gif", "*.ico", "*.svg", "*.webp", // Images
				"*.pdf", "*.zip", "*.tar", "*.gz", "*.bz2", // Archives
				"*.exe", "*.dll", "*.so", "*.dylib", // Binaries
				"*.woff", "*.woff2", "*.ttf", "*.eot", // Fonts
			},
		},
	}
}

// LoadConfig loads configuration from a file path
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// SaveConfig saves configuration to a file path
func SaveConfig(config *Config, path string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// FindConfig searches for a config file starting from the given directory and walking up
func FindConfig(startDir string) (string, error) {
	dir := startDir
	for {
		configPath := filepath.Join(dir, ".leakyrepo.yml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root directory
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("config file not found")
}

