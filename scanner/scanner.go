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

// isBinaryFile checks if a file is binary by examining its content
func isBinaryFile(content []byte) bool {
	// Check for common binary file signatures
	if len(content) < 4 {
		return false
	}
	
	// PNG: 89 50 4E 47
	if len(content) >= 4 && content[0] == 0x89 && content[1] == 0x50 && content[2] == 0x4E && content[3] == 0x47 {
		return true
	}
	
	// JPEG: FF D8 FF
	if len(content) >= 3 && content[0] == 0xFF && content[1] == 0xD8 && content[2] == 0xFF {
		return true
	}
	
	// GIF: 47 49 46 38
	if len(content) >= 4 && content[0] == 0x47 && content[1] == 0x49 && content[2] == 0x46 && content[3] == 0x38 {
		return true
	}
	
	// Check if file contains too many null bytes or non-printable characters
	// If more than 30% of bytes are non-printable (excluding common whitespace), it's likely binary
	nonPrintableCount := 0
	for i := 0; i < len(content) && i < 512; i++ { // Check first 512 bytes
		b := content[i]
		if b < 32 && b != '\t' && b != '\n' && b != '\r' {
			nonPrintableCount++
		}
	}
	
	checkedBytes := len(content)
	if checkedBytes > 512 {
		checkedBytes = 512
	}
	
	if checkedBytes > 0 && float64(nonPrintableCount)/float64(checkedBytes) > 0.3 {
		return true
	}
	
	return false
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

	// Skip binary files
	if isBinaryFile(content) {
		return nil, nil
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
		// Minimum length must match IsHighEntropy requirement (16 chars)
		if len(token) >= 16 && IsHighEntropy(token, s.config.EntropyThreshold) {
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
// Focuses on common secret patterns (key=value, key:value) rather than aggressive splitting
func tokenizeLine(line string) []string {
	var tokens []string
	
	// First, try to extract values from common key-value patterns
	// Pattern: key=value or key:value (common in .env, config files)
	kvPatterns := []string{"=", ":", " = ", " : "}
	for _, sep := range kvPatterns {
		if strings.Contains(line, sep) {
			parts := strings.SplitN(line, sep, 2)
			if len(parts) == 2 {
				// Extract the value part (right side)
				value := strings.TrimSpace(parts[1])
				// Remove quotes if present
				value = strings.Trim(value, "\"'`")
				if len(value) >= 16 {
					tokens = append(tokens, value)
				}
			}
		}
	}
	
	// If no key-value patterns found, split by whitespace and common delimiters
	// but only for longer tokens to reduce false positives
	if len(tokens) == 0 {
		delimiters := []string{" ", "\t", ",", ";", "|"}
		currentTokens := []string{line}
		
		for _, delim := range delimiters {
			var newTokens []string
			for _, token := range currentTokens {
				parts := strings.Split(token, delim)
				for _, part := range parts {
					part = strings.TrimSpace(part)
					// Only keep tokens that are long enough and don't look like code
					if len(part) >= 16 && !strings.Contains(part, "${") && 
					   !strings.HasPrefix(part, "<") && !strings.Contains(part, "[") {
						newTokens = append(newTokens, part)
					}
				}
			}
			currentTokens = newTokens
		}
		tokens = currentTokens
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
		// Try matching against absolute path
		if matched, _ := filepath.Match(pattern, filePath); matched {
			return true
		}
		// Try matching against base filename
		if matched, _ := filepath.Match(pattern, filepath.Base(filePath)); matched {
			return true
		}
		// Try matching against relative path from workDir (important for patterns added via interactive mode)
		if relPath, err := filepath.Rel(s.workDir, filePath); err == nil {
			if matched, _ := filepath.Match(pattern, relPath); matched {
				return true
			}
			// Also try normalized path separators (handle Windows vs Unix paths)
			normalizedPattern := strings.ReplaceAll(pattern, "\\", "/")
			normalizedRelPath := strings.ReplaceAll(relPath, "\\", "/")
			if matched, _ := filepath.Match(normalizedPattern, normalizedRelPath); matched {
				return true
			}
		}
		// Support directory patterns ending with /
		if strings.HasSuffix(pattern, "/") {
			dirPattern := strings.TrimSuffix(pattern, "/")
			// Try absolute path first
			if strings.HasPrefix(filePath, dirPattern+"/") {
				return true
			}
			// Also check relative path from workDir
			if relPath, err := filepath.Rel(s.workDir, filePath); err == nil {
				if strings.HasPrefix(relPath, dirPattern+"/") || relPath == dirPattern {
					return true
				}
				// Also check normalized paths
				normalizedPattern := strings.ReplaceAll(dirPattern, "\\", "/")
				normalizedRelPath := strings.ReplaceAll(relPath, "\\", "/")
				if strings.HasPrefix(normalizedRelPath, normalizedPattern+"/") || normalizedRelPath == normalizedPattern {
					return true
				}
			}
		}
		// Support directory patterns ending with /**
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
				// Also check normalized paths
				normalizedPattern := strings.ReplaceAll(dirPattern, "\\", "/")
				normalizedRelPath := strings.ReplaceAll(relPath, "\\", "/")
				if strings.HasPrefix(normalizedRelPath, normalizedPattern+"/") {
					return true
				}
			}
		}
	}

	return false
}

