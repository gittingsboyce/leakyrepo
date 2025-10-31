package scanner

import (
	"math"
	"testing"
)

func TestCalculateShannonEntropy(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		tolerance float64
	}{
		{
			name:      "empty string",
			input:     "",
			expected:  0.0,
			tolerance: 0.01,
		},
		{
			name:      "single character",
			input:     "a",
			expected:  0.0,
			tolerance: 0.01,
		},
		{
			name:      "low entropy (repeated character)",
			input:     "aaaaaaaa",
			expected:  0.0,
			tolerance: 0.1,
		},
		{
			name:      "medium entropy (simple pattern)",
			input:     "abcabcabc",
			expected:  1.0,
			tolerance: 0.5,
		},
		{
			name:      "high entropy (random-like string)",
			input:     "AKIAIOSFODNN7EXAMPLE",
			expected:  4.0,
			tolerance: 0.5,
		},
		{
			name:      "very high entropy (random alphanumeric)",
			input:     "a1b2c3d4e5f6g7h8i9j0",
			expected:  4.0,
			tolerance: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateShannonEntropy(tt.input)
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("CalculateShannonEntropy(%q) = %f, expected ~%f (tolerance: %f)",
					tt.input, result, tt.expected, tt.tolerance)
			}
		})
	}
}

func TestIsHighEntropy(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		threshold float64
		expected  bool
	}{
		{
			name:      "low entropy string",
			input:     "aaaaaaaa",
			threshold: 4.5,
			expected:  false,
		},
		{
			name:      "high entropy string",
			input:     "AKIAIOSFODNN7EXAMPLE",
			threshold: 4.0,
			expected:  true,
		},
		{
			name:      "short string (should be false)",
			input:     "abc123",
			threshold: 2.0,
			expected:  false, // Too short
		},
		{
			name:      "mostly whitespace (should be false)",
			input:     "a       b       c       d",
			threshold: 2.0,
			expected:  false, // Too much whitespace
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsHighEntropy(tt.input, tt.threshold)
			if result != tt.expected {
				t.Errorf("IsHighEntropy(%q, %f) = %v, expected %v",
					tt.input, tt.threshold, result, tt.expected)
			}
		})
	}
}

