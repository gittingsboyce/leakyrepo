package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/lgboyce/leakyrepo/config"
)

// Scanner scans files for secrets using regex rules and entropy detection
type Scanner struct {
	config        *config.Config
	compiledRules []compiledRule
	ignorePatterns []string
	workDir       string
}

type compiledRule struct {
	rule    config.Rule
	pattern *regexp.Regexp
}

// NewScanner creates a new scanner with the given configuration
func NewScanner(cfg *config.Config, ignorePatterns []string) (*Scanner, error) {
	workDir, _ := os.Getwd() // Get current working directory for relative path matching
	scanner := &Scanner{
		config:        cfg,
		compiledRules: make([]compiledRule, 0, len(cfg.Rules)),
		ignorePatterns: ignorePatterns,
		workDir:       workDir,
	}

	// Compile regex patterns
	for _, rule := range cfg.Rules {
		pattern, err := regexp.Compile(rule.Pattern)
		if err != nil {
			return nil, fmt.Errorf("failed to compile pattern for rule %s: %w", rule.ID, err)
		}
		scanner.compiledRules = append(scanner.compiledRules, compiledRule{
			rule:    rule,
			pattern: pattern,
		})
	}

	return scanner, nil
}

// ScanFile scans a single file for secrets
func (s *Scanner) ScanFile(filePath string) ([]Result, error) {
	// Check if file should be ignored
	if s.shouldIgnoreFile(filePath) {
		return nil, nil
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	var results []Result
	lines := strings.Split(string(content), "\n")
	fileExt := strings.ToLower(filepath.Ext(filePath))

	// Scan each line
	for lineNum, line := range lines {
		lineResults := s.scanLine(line, lineNum+1, filePath, fileExt)
		results = append(results, lineResults...)
	}

	return results, nil
}

// scanLine scans a single line for secrets
func (s *Scanner) scanLine(line string, lineNum int, filePath, fileExt string) []Result {
	var results []Result

	// Check against allowlist
	if s.config.Allowlist.Strings != nil {
		for _, allowed := range s.config.Allowlist.Strings {
			if strings.Contains(line, allowed) {
				return nil // This line is allowlisted
			}
		}
	}

	// Apply regex rules
	for _, compiled := range s.compiledRules {
		// Check if rule applies to this file type
		if len(compiled.rule.FileTypes) > 0 {
			applies := false
			for _, ft := range compiled.rule.FileTypes {
				if fileExt == strings.ToLower(ft) {
					applies = true
					break
				}
			}
			if !applies {
				continue
			}
		}

		// Find all matches
		matches := compiled.pattern.FindAllString(line, -1)
		for _, match := range matches {
			// Check if this match is allowlisted
			allowlisted := false
			for _, allowed := range s.config.Allowlist.Strings {
				if strings.Contains(match, allowed) {
					allowlisted = true
					break
				}
			}
			if allowlisted {
				continue
			}

			maskedMatch := MaskMatch(match, 4)
			results = append(results, Result{
				File:          filePath,
				Line:          lineNum,
				RuleID:        compiled.rule.ID,
				Severity:      compiled.rule.Severity,
				Match:         maskedMatch,
				Description:   compiled.rule.Description,
				DetectionType: "regex",
				ScannedAt:     time.Now(),
			})
		}
	}

	// Check for high-entropy strings
	// Split line by common delimiters and check each token
	tokens := tokenizeLine(line)
	for _, token := range tokens {
		if len(token) >= 8 && IsHighEntropy(token, s.config.EntropyThreshold) {
			// Check if token is allowlisted
			allowlisted := false
			for _, allowed := range s.config.Allowlist.Strings {
				if strings.Contains(token, allowed) {
					allowlisted = true
					break
				}
			}
			if allowlisted {
				continue
			}

			// Check if this high-entropy string was already matched by a regex rule
			alreadyMatched := false
			for _, result := range results {
				if result.Line == lineNum && strings.Contains(result.Match, token) {
					alreadyMatched = true
					break
				}
			}
			if !alreadyMatched {
				maskedMatch := MaskMatch(token, 4)
				results = append(results, Result{
					File:          filePath,
					Line:          lineNum,
					Severity:      "medium",
					Match:         maskedMatch,
					Description:   "High-entropy string detected (possible secret)",
					DetectionType: "entropy",
					ScannedAt:     time.Now(),
				})
			}
		}
	}

	return results
}

// tokenizeLine splits a line into potential secret tokens
func tokenizeLine(line string) []string {
	// Common delimiters for secrets in config files
	delimiters := []string{" ", "\t", "=", ":", ",", ";", "|", "\"", "'", "`"}
	tokens := []string{line}

	for _, delim := range delimiters {
		var newTokens []string
		for _, token := range tokens {
			parts := strings.Split(token, delim)
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if len(part) > 0 {
					newTokens = append(newTokens, part)
				}
			}
		}
		tokens = newTokens
	}

	return tokens
}

// shouldIgnoreFile checks if a file should be ignored based on patterns
func (s *Scanner) shouldIgnoreFile(filePath string) bool {
	// Check config allowlist
	for _, pattern := range s.config.Allowlist.Files {
		if matched, _ := filepath.Match(pattern, filePath); matched {
			return true
		}
		if matched, _ := filepath.Match(pattern, filepath.Base(filePath)); matched {
			return true
		}
	}

	// Check .leakyrepoignore patterns
	for _, pattern := range s.ignorePatterns {
		if matched, _ := filepath.Match(pattern, filePath); matched {
			return true
		}
		if matched, _ := filepath.Match(pattern, filepath.Base(filePath)); matched {
			return true
		}
		// Support directory patterns
		if strings.HasSuffix(pattern, "/**") {
			dirPattern := strings.TrimSuffix(pattern, "/**")
			// Try absolute path first
			if strings.HasPrefix(filePath, dirPattern+"/") {
				return true
			}
			// Also check relative path from workDir
			if relPath, err := filepath.Rel(s.workDir, filePath); err == nil {
				if strings.HasPrefix(relPath, dirPattern+"/") {
					return true
				}
			}
		}
	}

	return false
}

