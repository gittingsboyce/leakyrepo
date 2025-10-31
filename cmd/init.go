package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lgboyce/leakyrepo/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize LeakyRepo configuration",
	Long: `Creates a default .leakyrepo.yml configuration file in the current directory.
This file contains regex rules for common secret patterns and entropy thresholds.`,
	RunE: runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	workDir := getWorkingDir()
	configPath := filepath.Join(workDir, ".leakyrepo.yml")

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("configuration file already exists: %s\nUse --force to overwrite", configPath)
	}

	// Create default config
	cfg := config.DefaultConfig()

	// Save config
	if err := config.SaveConfig(cfg, configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("✓ Created default configuration at %s\n", configPath)
	fmt.Println("\nConfiguration includes:")
	fmt.Println("  - AWS Access Key detection")
	fmt.Println("  - Generic API Key pattern detection")
	fmt.Println("  - High-entropy string detection (threshold: 4.5)")
	fmt.Println("\nYou can customize the configuration by editing .leakyrepo.yml")

	return nil
}

