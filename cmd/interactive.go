package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lgboyce/leakyrepo/scanner"
)

// runInteractive handles interactive mode for ignoring false positives
func runInteractive(results []scanner.Result, workDir string) (bool, error) {
	if len(results) == 0 {
		return true, nil
	}

	fmt.Println()
	fmt.Println("Would you like to ignore any of these findings?")
	fmt.Println()

	// Group results by file
	fileGroups := make(map[string][]scanner.Result)
	for _, result := range results {
		fileGroups[result.File] = append(fileGroups[result.File], result)
	}

	reader := bufio.NewReader(os.Stdin)

	// Show options for each file
	ignoreFiles := make(map[string]bool)
	ignorePatterns := make([]string, 0)

	for file, fileResults := range fileGroups {
		relFile, _ := filepath.Rel(workDir, file)
		if relFile == "" {
			relFile = file
		}

		fmt.Printf("File: %s (%d finding(s))\n", relFile, len(fileResults))
		fmt.Println("  Options:")
		fmt.Println("    [f] Ignore this entire file")
		fmt.Println("    [n] Next file (keep findings)")
		fmt.Println("    [a] Ignore all remaining files")
		fmt.Println("    [q] Quit (manual fix required)")
		fmt.Print("  Choice [f/n/a/q]: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "f":
			ignoreFiles[relFile] = true
			fmt.Printf("  ✓ Will ignore file: %s\n", relFile)
		case "a":
			// Ignore all remaining files
			for f := range fileGroups {
				rel, _ := filepath.Rel(workDir, f)
				if rel == "" {
					rel = f
				}
				ignoreFiles[rel] = true
			}
			fmt.Println("  ✓ Will ignore all remaining files")
			return applyIgnores(ignoreFiles, ignorePatterns, workDir)
		case "q":
			fmt.Println("\nQuitting. Manual fix required.")
			fmt.Println("To ignore files, run: leakyrepo ignore <file>")
			return false, nil
		case "n":
			fmt.Println("  → Keeping findings for this file")
		default:
			fmt.Println("  → Keeping findings for this file")
		}
		fmt.Println()
	}

	return applyIgnores(ignoreFiles, ignorePatterns, workDir)
}

// applyIgnores adds ignored files to .leakyrepoignore
func applyIgnores(ignoreFiles map[string]bool, ignorePatterns []string, workDir string) (bool, error) {
	if len(ignoreFiles) == 0 && len(ignorePatterns) == 0 {
		return false, nil
	}

	ignorePath := filepath.Join(workDir, ".leakyrepoignore")

	// Read existing patterns
	existingPatterns := make(map[string]bool)
	if _, err := os.Stat(ignorePath); err == nil {
		data, err := os.ReadFile(ignorePath)
		if err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" && !strings.HasPrefix(line, "#") {
					existingPatterns[line] = true
				}
			}
		}
	}

	// Open file for appending
	file, err := os.OpenFile(ignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return false, fmt.Errorf("failed to open .leakyrepoignore: %w", err)
	}
	defer file.Close()

	// Add comment if file is empty
	info, _ := file.Stat()
	if info.Size() == 0 {
		fmt.Fprintf(file, "# LeakyRepo Ignore File\n")
		fmt.Fprintf(file, "# Files and patterns listed here will be ignored during scanning\n\n")
	}

	// Add new patterns
	added := 0
	for filePattern := range ignoreFiles {
		if !existingPatterns[filePattern] {
			if _, err := fmt.Fprintf(file, "%s\n", filePattern); err != nil {
				return false, fmt.Errorf("failed to write to .leakyrepoignore: %w", err)
			}
			existingPatterns[filePattern] = true
			added++
		}
	}

	for _, pattern := range ignorePatterns {
		if !existingPatterns[pattern] {
			if _, err := fmt.Fprintf(file, "%s\n", pattern); err != nil {
				return false, fmt.Errorf("failed to write to .leakyrepoignore: %w", err)
			}
			existingPatterns[pattern] = true
			added++
		}
	}

	if added > 0 {
		fmt.Printf("\n✓ Added %d pattern(s) to .leakyrepoignore\n", added)
		fmt.Println("✓ Re-running scan...")
		fmt.Println()
		return true, nil
	}

	return false, nil
}

