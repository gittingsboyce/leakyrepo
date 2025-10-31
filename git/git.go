package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetStagedFiles returns a list of files that are staged for commit
func GetStagedFiles(repoRoot string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only", "--diff-filter=ACMR")
	cmd.Dir = repoRoot
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get staged files: %w", err)
	}

	if len(output) == 0 {
		return []string{}, nil
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	var stagedFiles []string
	for _, file := range files {
		if file != "" {
			// Convert to absolute path
			absPath := filepath.Join(repoRoot, file)
			stagedFiles = append(stagedFiles, absPath)
		}
	}

	return stagedFiles, nil
}

// GetRepoRoot returns the root directory of the git repository
func GetRepoRoot(startDir string) (string, error) {
	dir := startDir
	for {
		// Try to get git root using git command
		cmd := exec.Command("git", "rev-parse", "--show-toplevel")
		cmd.Dir = dir
		output, err := cmd.Output()
		if err == nil {
			return strings.TrimSpace(string(output)), nil
		}

		// Fallback: Check if .git directory exists
		gitDir := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root directory
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("not a git repository")
}

