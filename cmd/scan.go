package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lgboyce/leakyrepo/config"
	"github.com/lgboyce/leakyrepo/git"
	"github.com/lgboyce/leakyrepo/ignore"
	"github.com/lgboyce/leakyrepo/scanner"
	"github.com/spf13/cobra"
)

var (
	jsonOutput string
	explain    bool
	interactive bool
	scanAll   bool
)

var scanCmd = &cobra.Command{
	Use:   "scan [files...]",
	Short: "Scan files for secrets",
	Long: `Scan files for secrets using regex rules and entropy detection.
If no files are specified, scans staged files in the git repository.`,
	RunE: runScan,
}

func init() {
	scanCmd.Flags().StringVar(&jsonOutput, "json", "", "Output results to JSON file")
	scanCmd.Flags().BoolVar(&explain, "explain", false, "Show explanation for each detected secret")
	scanCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive mode: prompt to ignore false positives")
	scanCmd.Flags().BoolVar(&scanAll, "all", false, "Scan all tracked files in the repository (default: scan staged files)")
}

func runScan(cmd *cobra.Command, args []string) error {
	workDir := getWorkingDir()

	// Find and load config
	configPath, err := findConfigPath(workDir)
	if err != nil {
		return fmt.Errorf("failed to find config: %w\nRun 'leakyrepo init' to create a default configuration", err)
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Load ignore patterns
	ignorePath := filepath.Join(workDir, ".leakyrepoignore")
	ignorePatterns, err := ignore.LoadIgnorePatterns(ignorePath)
	if err != nil {
		return fmt.Errorf("failed to load ignore patterns: %w", err)
	}

	// Create scanner
	scnr, err := scanner.NewScanner(cfg, ignorePatterns)
	if err != nil {
		return fmt.Errorf("failed to create scanner: %w", err)
	}

	// Determine files to scan
	var filesToScan []string
	if len(args) > 0 {
		// Use files provided as arguments
		for _, arg := range args {
			absPath, err := filepath.Abs(arg)
			if err != nil {
				return fmt.Errorf("failed to get absolute path for %s: %w", arg, err)
			}
			if _, err := os.Stat(absPath); err != nil {
				return fmt.Errorf("file not found: %s", arg)
			}
			filesToScan = append(filesToScan, absPath)
		}
	} else {
		// Get files from git
		repoRoot, err := git.GetRepoRoot(workDir)
		if err != nil {
			return fmt.Errorf("failed to find git repository: %w\nSpecify files to scan or run from within a git repository", err)
		}

		if scanAll {
			// Get all tracked files
			trackedFiles, err := git.GetAllTrackedFiles(repoRoot)
			if err != nil {
				return fmt.Errorf("failed to get tracked files: %w", err)
			}

			if len(trackedFiles) == 0 {
				fmt.Println("No tracked files in repository.")
				return nil
			}

			filesToScan = trackedFiles
		} else {
			// Get staged files from git
			stagedFiles, err := git.GetStagedFiles(repoRoot)
			if err != nil {
				return fmt.Errorf("failed to get staged files: %w", err)
			}

			if len(stagedFiles) == 0 {
				fmt.Println("No files staged for commit.")
				return nil
			}

			filesToScan = stagedFiles
		}
	}

	// Scan files
	var allResults []scanner.Result
	for _, file := range filesToScan {
		results, err := scnr.ScanFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to scan %s: %v\n", file, err)
			continue
		}
		allResults = append(allResults, results...)
	}

	// Output results
	if jsonOutput != "" {
		if err := outputJSON(allResults, jsonOutput); err != nil {
			return fmt.Errorf("failed to write JSON output: %w", err)
		}
		fmt.Printf("Results written to %s\n", jsonOutput)
	} else {
		outputHumanReadable(allResults, explain)
	}

	// Handle interactive mode
	if interactive && len(allResults) > 0 {
		shouldRescan, err := runInteractive(allResults, workDir)
		if err != nil {
			return fmt.Errorf("interactive mode error: %w", err)
		}

		if shouldRescan {
			// Reload ignore patterns after updating .leakyrepoignore
			ignorePatterns, err = ignore.LoadIgnorePatterns(filepath.Join(workDir, ".leakyrepoignore"))
			if err != nil {
				return fmt.Errorf("failed to reload ignore patterns: %w", err)
			}

			// Re-create scanner with updated ignore patterns
			scnr, err = scanner.NewScanner(cfg, ignorePatterns)
			if err != nil {
				return fmt.Errorf("failed to create scanner: %w", err)
			}

			// Re-scan files
			var newResults []scanner.Result
			for _, file := range filesToScan {
				results, err := scnr.ScanFile(file)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to scan %s: %v\n", file, err)
					continue
				}
				newResults = append(newResults, results...)
			}

			if len(newResults) > 0 {
				fmt.Printf("\n⚠️  Still found %d potential secret(s) after ignoring:\n\n", len(newResults))
				outputHumanReadable(newResults, explain)
				fmt.Println("\nPlease review remaining findings or add more ignore patterns.")
				os.Exit(1)
			} else {
				fmt.Println("✓ No secrets found! All ignored patterns applied successfully.")
				return nil
			}
		}
	}

	// Exit with error if secrets were found (and not in interactive mode)
	if len(allResults) > 0 {
		if !interactive {
			fmt.Println("\n💡 Tip: Run with --interactive (-i) to ignore false positives interactively")
			fmt.Println("   Or use: leakyrepo ignore <file>")
		}
		os.Exit(1)
	}

	return nil
}

func outputJSON(results []scanner.Result, outputPath string) error {
	// Convert to JSON format as specified
	type JSONResult struct {
		File    string `json:"file"`
		Line    int    `json:"line"`
		RuleID  string `json:"rule_id,omitempty"`
		Severity string `json:"severity"`
		Match   string `json:"match"`
	}

	jsonResults := make([]JSONResult, len(results))
	for i, r := range results {
		jsonResults[i] = JSONResult{
			File:     r.File,
			Line:     r.Line,
			RuleID:   r.RuleID,
			Severity: r.Severity,
			Match:    r.Match,
		}
	}

	data, err := json.MarshalIndent(jsonResults, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, data, 0644)
}

func outputHumanReadable(results []scanner.Result, explain bool) {
	if len(results) == 0 {
		fmt.Println("✓ No secrets found!")
		return
	}

	fmt.Printf("\n⚠️  Found %d potential secret(s):\n\n", len(results))

	severityColors := map[string]string{
		"low":      "\033[33m", // Yellow
		"medium":   "\033[33m", // Yellow
		"high":     "\033[31m", // Red
		"critical": "\033[31m", // Red
	}
	resetColor := "\033[0m"
	lockEmoji := "🔒"

	for _, result := range results {
		color := severityColors[result.Severity]
		if color == "" {
			color = "\033[33m" // Default to yellow
		}

		// Format: 🔒 [High] AWS Access Key found in config.env: AKIA...
		fmt.Printf("%s %s[%s] %s found in %s:%d\n",
			lockEmoji,
			color,
			result.Severity,
			result.Description,
			result.File,
			result.Line,
		)

		// Show masked match
		fmt.Printf("   Match: %s%s%s\n", color, result.Match, resetColor)

		// Show explanation if requested
		if explain {
			if result.DetectionType == "regex" {
				fmt.Printf("   Reason: Matched regex rule '%s' (pattern: %s)\n",
					result.RuleID,
					getRulePattern(result.RuleID),
				)
			} else if result.DetectionType == "entropy" {
				threshold := getEntropyThreshold()
				fmt.Printf("   Reason: High entropy detected (entropy above threshold: %.2f)\n",
					threshold,
				)
			}
		}

		fmt.Println()
	}
}

// Helper functions for explain output
func getRulePattern(ruleID string) string {
	// This would ideally load from config, but for simplicity we'll just return a placeholder
	workDir := getWorkingDir()
	configPath, err := findConfigPath(workDir)
	if err != nil {
		return "N/A"
	}
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return "N/A"
	}
	for _, rule := range cfg.Rules {
		if rule.ID == ruleID {
			return rule.Pattern
		}
	}
	return "N/A"
}

func getEntropyThreshold() float64 {
	workDir := getWorkingDir()
	configPath, err := findConfigPath(workDir)
	if err != nil {
		return 4.5
	}
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return 4.5
	}
	return cfg.EntropyThreshold
}

