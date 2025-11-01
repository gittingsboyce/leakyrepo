package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var ignoreCmd = &cobra.Command{
	Use:   "ignore",
	Short: "Add file or pattern to .leakyrepoignore",
	Long: `Add files or patterns to .leakyrepoignore to ignore false positives.

Examples:
  leakyrepo ignore scripts/test-homebrew.sh
  leakyrepo ignore --pattern "*test*.sh"
  leakyrepo ignore --file scripts/
`,
	RunE: runIgnore,
}

var (
	ignorePattern string
	ignoreFile    string
)

func init() {
	ignoreCmd.Flags().StringVar(&ignorePattern, "pattern", "", "Pattern to ignore (e.g., '*test*.sh')")
	ignoreCmd.Flags().StringVar(&ignoreFile, "file", "", "File or directory to ignore")
	// Note: Command is added in root.go to avoid duplicate registration
}

func runIgnore(cmd *cobra.Command, args []string) error {
	workDir := getWorkingDir()
	ignorePath := filepath.Join(workDir, ".leakyrepoignore")

	// Determine what to ignore
	var pattern string
	if ignorePattern != "" {
		pattern = ignorePattern
	} else if ignoreFile != "" {
		pattern = ignoreFile
	} else if len(args) > 0 {
		pattern = args[0]
	} else {
		return fmt.Errorf("specify a file, pattern, or use --file/--pattern flags")
	}

	// Read existing ignore file
	var existingPatterns []string
	if _, err := os.Stat(ignorePath); err == nil {
		data, err := os.ReadFile(ignorePath)
		if err != nil {
			return fmt.Errorf("failed to read .leakyrepoignore: %w", err)
		}
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				existingPatterns = append(existingPatterns, line)
			}
		}
	}

	// Check if pattern already exists
	for _, existing := range existingPatterns {
		if existing == pattern {
			fmt.Printf("Pattern '%s' already in .leakyrepoignore\n", pattern)
			return nil
		}
	}

	// Append new pattern
	file, err := os.OpenFile(ignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open .leakyrepoignore: %w", err)
	}
	defer file.Close()

	// Add comment if file is empty
	info, _ := file.Stat()
	if info.Size() == 0 {
		fmt.Fprintf(file, "# LeakyRepo Ignore File\n")
		fmt.Fprintf(file, "# Files and patterns listed here will be ignored during scanning\n\n")
	}

	// Add the pattern
	if _, err := fmt.Fprintf(file, "%s\n", pattern); err != nil {
		return fmt.Errorf("failed to write to .leakyrepoignore: %w", err)
	}

	fmt.Printf("✓ Added '%s' to .leakyrepoignore\n", pattern)
	return nil
}

