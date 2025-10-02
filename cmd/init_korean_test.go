package cmd

import (
	"testing"
)

func TestIsKoreanWordList(t *testing.T) {
	tests := []struct {
		name         string
		generatorKey string
		expected     bool
	}{
		// Embedded Korean wordlists
		{"Korean common words", "Korean common words", true},
		{"Korean tech terms", "Korean tech terms", true},
		{"Korean sentences", "Korean sentences", true},
		
		// Case variations
		{"Lowercase korean", "korean common words", true},
		{"Mixed case Korean", "Korean Tech Terms", true},
		{"UPPERCASE KOREAN", "KOREAN SENTENCES", true},
		
		// Korean language indicators
		{"Korean with 한국", "한국어 단어", true},
		{"Korean with 한글", "한글 연습", true},
		
		// File paths with Korean indicators
		{"Korean file path", "/path/to/korean_words.json", true},
		{"Korea in path", "/files/korea-tech.txt", true},
		{"Korean with underscore", "korean_words.json", true},
		{"Korean with dash", "korean-tech.txt", true},
		{"Korean with dot", "korean.json", true},
		{"Korea with underscore", "south_korea_cities.json", true},
		{"Korea with dash", "north-korea-data.txt", true},
		
		// Non-Korean wordlists
		{"English common", "Common words", false},
		{"Frankenstein", "Frankenstein sentences", false},
		{"Custom English", "/path/to/english_words.json", false},
		{"Random text", "some random text", false},
		{"Empty string", "", false},
		
		// Edge cases - should NOT match when korean/korea is part of another word
		{"Korean substring in word", "nonkoreanword", false},
		{"Korea substring in word", "nonkoreaword", false},
		
		// Should match when korean/korea is standalone or properly separated
		{"Just Korean", "Korean", true},
		{"Just korean", "korean", true},
		{"Just korea", "korea", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isKoreanWordList(tt.generatorKey)
			if result != tt.expected {
				t.Errorf("isKoreanWordList(%q) = %v, expected %v", tt.generatorKey, result, tt.expected)
			}
		})
	}
}

func TestIsKoreanWordListWithRealEmbeddedNames(t *testing.T) {
	// Test with the actual embedded wordlist names from the codebase
	realKoreanWordlists := []string{
		"Korean common words",
		"Korean tech terms", 
		"Korean sentences",
	}
	
	for _, wordlist := range realKoreanWordlists {
		t.Run("Real embedded: "+wordlist, func(t *testing.T) {
			if !isKoreanWordList(wordlist) {
				t.Errorf("isKoreanWordList(%q) should return true for real Korean wordlist", wordlist)
			}
		})
	}
	
	// Test with real English wordlists
	realEnglishWordlists := []string{
		"Common words",
		"Frankenstein sentences",
	}
	
	for _, wordlist := range realEnglishWordlists {
		t.Run("Real embedded: "+wordlist, func(t *testing.T) {
			if isKoreanWordList(wordlist) {
				t.Errorf("isKoreanWordList(%q) should return false for English wordlist", wordlist)
			}
		})
	}
}

func TestIsKoreanWordListWithFilePaths(t *testing.T) {
	// Test with realistic file paths
	tests := []struct {
		name     string
		filePath string
		expected bool
	}{
		{"Korean JSON file", "/home/user/.config/typioca/korean_words.json", true},
		{"Korean text file", "/tmp/korean-tech.txt", true},
		{"Korea region file", "/data/south-korea-cities.json", true},
		{"English JSON file", "/home/user/.config/typioca/english_words.json", false},
		{"Tech terms file", "/tmp/programming-terms.txt", false},
		{"Random file", "/data/random-data.json", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isKoreanWordList(tt.filePath)
			if result != tt.expected {
				t.Errorf("isKoreanWordList(%q) = %v, expected %v", tt.filePath, result, tt.expected)
			}
		})
	}
}