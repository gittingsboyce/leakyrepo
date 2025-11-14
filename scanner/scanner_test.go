package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lgboyce/leakyrepo/config"
)

func TestScanner_ScanFile(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.env")
	testContent := `# Test file
AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
API_KEY=sk_live_1234567890abcdefghijklmnopqrstuvwxyz
PASSWORD=secretpassword123
`
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create config with AWS access key rule
	cfg := &config.Config{
		EntropyThreshold: 4.5,
		Rules: []config.Rule{
			{
				ID:          "aws_access_key",
				Description: "AWS Access Key",
				Severity:    "high",
				Pattern:     `AKIA[0-9A-Z]{16}`,
				FileTypes:   []string{".env"},
			},
		},
	}

	// Create scanner
	scnr, err := NewScanner(cfg, []string{})
	if err != nil {
		t.Fatalf("Failed to create scanner: %v", err)
	}

	// Scan file
	results, err := scnr.ScanFile(testFile)
	if err != nil {
		t.Fatalf("Failed to scan file: %v", err)
	}

	// Check results
	if len(results) == 0 {
		t.Error("Expected to find at least one secret, but found none")
	}

	// Verify AWS access key was detected
	foundAWSKey := false
	for _, result := range results {
		if result.RuleID == "aws_access_key" && result.Severity == "high" {
			foundAWSKey = true
			if result.Line != 2 {
				t.Errorf("Expected AWS key on line 2, got line %d", result.Line)
			}
			break
		}
	}

	if !foundAWSKey {
		t.Error("Expected to find AWS access key, but it was not detected")
	}
}

func TestScanner_shouldIgnoreFile(t *testing.T) {
	cfg := &config.Config{
		Allowlist: config.Allowlist{
			Files: []string{".git/**", "*.lock"},
		},
	}

	ignorePatterns := []string{"node_modules/**", "vendor/"}

	scnr, err := NewScanner(cfg, ignorePatterns)
	if err != nil {
		t.Fatalf("Failed to create scanner: %v", err)
	}

	tests := []struct {
		name     string
		filePath string
		expected bool
	}{
		{
			name:     "git file should be ignored",
			filePath: ".git/config",
			expected: true,
		},
		{
			name:     "lock file should be ignored",
			filePath: "package-lock.json",
			expected: true,
		},
		{
			name:     "node_modules should be ignored",
			filePath: "node_modules/package/index.js",
			expected: true,
		},
		{
			name:     "regular file should not be ignored",
			filePath: "src/main.go",
			expected: false,
		},
	}
	
	// Add relative path pattern to ignorePatterns for the new test
	// This tests that a relative path pattern (like those added via interactive mode)
	// can match an absolute file path
	ignorePatterns = append(ignorePatterns, "src/pages/Index.tsx")
	scnr.ignorePatterns = ignorePatterns
	
	// Get current working directory for the absolute path test
	workDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	
	tests = append(tests, struct {
		name     string
		filePath string
		expected bool
	}{
		name:     "relative path pattern should match absolute path",
		filePath: filepath.Join(workDir, "src/pages/Index.tsx"),
		expected: true,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scnr.shouldIgnoreFile(tt.filePath)
			if result != tt.expected {
				t.Errorf("shouldIgnoreFile(%q) = %v, expected %v",
					tt.filePath, result, tt.expected)
			}
		})
	}
}

