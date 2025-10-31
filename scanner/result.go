package scanner

import "time"

// Result represents a detected secret
type Result struct {
	// File is the path to the file where the secret was found
	File string `json:"file"`
	// Line is the line number where the secret was found (1-indexed)
	Line int `json:"line"`
	// RuleID is the identifier of the rule that matched (empty for entropy-based detection)
	RuleID string `json:"rule_id,omitempty"`
	// Severity indicates the severity level
	Severity string `json:"severity"`
	// Match is the matched string (may be masked)
	Match string `json:"match"`
	// Description explains what was detected
	Description string `json:"description,omitempty"`
	// DetectionType is either "regex" or "entropy"
	DetectionType string `json:"detection_type"`
	// ScannedAt is the timestamp when this result was generated
	ScannedAt time.Time `json:"scanned_at"`
}

// MaskMatch masks a sensitive string, showing only first and last few characters
func MaskMatch(s string, visibleChars int) string {
	if len(s) <= visibleChars*2 {
		// Too short to mask effectively
		return "***"
	}
	if visibleChars <= 0 {
		visibleChars = 4
	}
	return s[:visibleChars] + "***" + s[len(s)-visibleChars:]
}

