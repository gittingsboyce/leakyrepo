package scanner

import (
	"math"
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
	// Filter out very short strings and common patterns
	if len(s) < 8 {
		return false
	}

	// Skip strings that are mostly whitespace or punctuation
	alphanumericCount := 0
	for _, char := range s {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			alphanumericCount++
		}
	}
	if alphanumericCount < len(s)/2 {
		return false
	}

	entropy := CalculateShannonEntropy(s)
	return entropy >= threshold
}

