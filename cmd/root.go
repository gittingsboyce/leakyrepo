package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "leakyrepo",
	Short: "LeakyRepo - A secrets detection tool",
	Long: `LeakyRepo is a secrets detection tool that combines regex-based 
pattern matching with entropy-based detection to catch API keys, tokens, 
and credentials before they're committed to version control.`,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(installHookCmd)
	rootCmd.AddCommand(ignoreCmd)
}

// getWorkingDir returns the current working directory
func getWorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
		os.Exit(1)
	}
	return dir
}

// findConfigPath searches for .leakyrepo.yml starting from the given directory
func findConfigPath(startDir string) (string, error) {
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

	return "", fmt.Errorf("config file .leakyrepo.yml not found")
}

