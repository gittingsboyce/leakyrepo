package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	if cfg.EntropyThreshold != 4.5 {
		t.Errorf("Expected EntropyThreshold to be 4.5, got %f", cfg.EntropyThreshold)
	}

	if len(cfg.Rules) == 0 {
		t.Error("Expected default config to have at least one rule")
	}

	// Check for AWS access key rule
	foundAWSRule := false
	for _, rule := range cfg.Rules {
		if rule.ID == "aws_access_key" {
			foundAWSRule = true
			if rule.Severity != "high" {
				t.Errorf("Expected AWS rule severity to be 'high', got %q", rule.Severity)
			}
			break
		}
	}

	if !foundAWSRule {
		t.Error("Expected default config to include AWS access key rule")
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".leakyrepo.yml")

	// Create default config
	cfg := DefaultConfig()

	// Save config
	if err := SaveConfig(cfg, configPath); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Check file exists
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("Config file was not created: %v", err)
	}

	// Load config
	loadedCfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify loaded config matches
	if loadedCfg.EntropyThreshold != cfg.EntropyThreshold {
		t.Errorf("EntropyThreshold mismatch: expected %f, got %f",
			cfg.EntropyThreshold, loadedCfg.EntropyThreshold)
	}

	if len(loadedCfg.Rules) != len(cfg.Rules) {
		t.Errorf("Rules count mismatch: expected %d, got %d",
			len(cfg.Rules), len(loadedCfg.Rules))
	}
}

