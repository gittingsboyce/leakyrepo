package scanner

import (
	"math"
	"strings"
)

// CalculateShannonEntropy calculates the Shannon entropy of a string
// Higher entropy values indicate more randomness, which is often a sign of secrets
// Formula: H(X) = -Σ P(x) * log2(P(x))
// Returns entropy value (typically between 0 and 8 for ASCII strings)
func CalculateShannonEntropy(s string) float64 {
	if len(s) == 0 {
		return 0
	}

	// Count frequency of each character
	freq := make(map[rune]int)
	for _, char := range s {
		freq[char]++
	}

	// Calculate entropy
	entropy := 0.0
	length := float64(len(s))

	for _, count := range freq {
		probability := float64(count) / length
		if probability > 0 {
			entropy -= probability * math.Log2(probability)
		}
	}

	return entropy
}

// IsHighEntropy checks if a string has entropy above the threshold
func IsHighEntropy(s string, threshold float64) bool {
	// Filter out very short strings - require at least 16 characters for entropy detection
	// This reduces false positives from code constructs
	if len(s) < 16 {
		return false
	}

	// Skip strings that are mostly whitespace or punctuation
	alphanumericCount := 0
	hasPrintableASCII := false
	for _, char := range s {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			alphanumericCount++
		}
		// Check for non-printable characters (binary data)
		if char < 32 && char != '\t' && char != '\n' && char != '\r' {
			return false // Contains binary data
		}
		if char >= 32 && char <= 126 {
			hasPrintableASCII = true
		}
	}
	
	// Require at least 60% alphanumeric characters
	if alphanumericCount < len(s)*3/5 {
		return false
	}
	
	// Must have printable ASCII characters
	if !hasPrintableASCII {
		return false
	}

	// Filter out common code patterns that have high entropy
	// Template strings like ${variable}, ${im***lips
	if strings.Contains(s, "${") || strings.Contains(s, "$(") {
		return false
	}
	
	// JSX/HTML tags like </Dr***tem>
	if strings.HasPrefix(s, "</") || strings.HasPrefix(s, "<") && strings.Contains(s, ">") {
		return false
	}
	
	// CSS class patterns like w-[c***4))]
	if strings.Contains(s, "[") && strings.Contains(s, "]") {
		return false
	}
	
	// Common code patterns with parentheses
	if strings.Count(s, "(") > 2 || strings.Count(s, ")") > 2 {
		return false
	}
	
	// Variable names with common prefixes/suffixes
	if strings.HasPrefix(s, "play") || strings.HasPrefix(s, "Drop") || 
	   strings.HasSuffix(s, "Name") || strings.HasSuffix(s, "DATE") || 
	   strings.HasSuffix(s, "NGED") {
		return false
	}

	entropy := CalculateShannonEntropy(s)
	return entropy >= threshold
}

