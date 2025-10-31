package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var installHookCmd = &cobra.Command{
	Use:   "install-hook",
	Short: "Install Git pre-commit hook",
	Long: `Installs a Git pre-commit hook that will automatically scan staged files
for secrets before allowing a commit. The commit will be blocked if secrets are detected.`,
	RunE: runInstallHook,
}

func runInstallHook(cmd *cobra.Command, args []string) error {
	workDir := getWorkingDir()

	// Find .git directory
	gitDir := filepath.Join(workDir, ".git")
	if _, err := os.Stat(gitDir); err != nil {
		return fmt.Errorf("not a git repository: .git directory not found")
	}

	// Check if git hooks directory exists
	hooksDir := filepath.Join(gitDir, "hooks")
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return fmt.Errorf("failed to create hooks directory: %w", err)
	}

	// Path to pre-commit hook
	hookPath := filepath.Join(hooksDir, "pre-commit")

	// Check if hook already exists
	if _, err := os.Stat(hookPath); err == nil {
		return fmt.Errorf("pre-commit hook already exists at %s\nRemove it first to reinstall", hookPath)
	}

	// Get absolute path to leakyrepo binary
	// We'll use the binary name directly and assume it's in PATH
	// In a real scenario, you might want to detect the binary location
	binaryName := "leakyrepo"

	// Check if binary exists in PATH
	if _, err := exec.LookPath(binaryName); err != nil {
		// Try to find the binary relative to current directory (for development)
		possiblePaths := []string{
			filepath.Join(workDir, "leakyrepo"),
			filepath.Join(workDir, "bin", "leakyrepo"),
			filepath.Join(workDir, "dist", "leakyrepo"),
		}

		found := false
		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				absPath, _ := filepath.Abs(path)
				binaryName = absPath
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("leakyrepo binary not found in PATH or local directories.\nPlease build the binary first with 'go build -o leakyrepo'")
		}
	}

	// Create pre-commit hook script
	hookScript := fmt.Sprintf(`#!/bin/sh
# LeakyRepo pre-commit hook
# This hook scans staged files for secrets

# Run leakyrepo scan (this will scan staged files automatically)
%s scan

# Exit with the same code as leakyrepo
exit_code=$?
if [ $exit_code -ne 0 ]; then
    echo ""
    echo "❌ Commit blocked: secrets detected in staged files!"
    echo "Review the findings above and remove secrets before committing."
    exit $exit_code
fi

exit 0
`, binaryName)

	// Write hook script
	if err := os.WriteFile(hookPath, []byte(hookScript), 0755); err != nil {
		return fmt.Errorf("failed to write pre-commit hook: %w", err)
	}

	// Make hook executable
	if err := os.Chmod(hookPath, 0755); err != nil {
		return fmt.Errorf("failed to make hook executable: %w", err)
	}

	fmt.Printf("✓ Pre-commit hook installed at %s\n", hookPath)
	fmt.Println("\nThe hook will now scan staged files before each commit.")
	fmt.Println("If secrets are detected, the commit will be blocked.")

	return nil
}

