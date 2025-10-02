package cmd

import (
	"testing"
)

func TestGetKoreanWordListMapping(t *testing.T) {
	mapping := getKoreanWordListMapping()
	
	// Test that all expected mappings exist
	expectedMappings := []struct {
		key            string
		expectedType   KoreanGenerationType
		expectedListName string
	}{
		{"korean-common", KoreanWords, "Korean common words"},
		{"korean-tech", KoreanWords, "Korean tech terms"},
		{"korean-sentences", KoreanSentences, "Korean sentences"},
		{"Korean common words", KoreanWords, "Korean common words"},
		{"Korean tech terms", KoreanWords, "Korean tech terms"},
		{"Korean sentences", KoreanSentences, "Korean sentences"},
	}
	
	for _, expected := range expectedMappings {
		if koreanMapping, exists := mapping[expected.key]; exists {
			if koreanMapping.GenerationType != expected.expectedType {
				t.Errorf("Expected generation type %v for key %s, got %v", 
					expected.expectedType, expected.key, koreanMapping.GenerationType)
			}
			if koreanMapping.ListName != expected.expectedListName {
				t.Errorf("Expected list name %s for key %s, got %s", 
					expected.expectedListName, expected.key, koreanMapping.ListName)
			}
		} else {
			t.Errorf("Expected mapping for key %s not found", expected.key)
		}
	}
}

func TestGetKoreanGenerationType(t *testing.T) {
	testCases := []struct {
		generatorKey   string
		shouldExist    bool
		expectedType   KoreanGenerationType
	}{
		// Direct mappings
		{"korean-common", true, KoreanWords},
		{"korean-tech", true, KoreanWords},
		{"korean-sentences", true, KoreanSentences},
		{"Korean common words", true, KoreanWords},
		{"Korean tech terms", true, KoreanWords},
		{"Korean sentences", true, KoreanSentences},
		
		// Fallback detection
		{"some-korean-words", true, KoreanWords},
		{"korean-sentence-list", true, KoreanSentences},
		{"custom-korean-문장", true, KoreanSentences},
		{"custom-korean-단어", true, KoreanWords},
		
		// Non-Korean keys
		{"english-words", false, ""},
		{"common-words", false, ""},
		{"frankenstein", false, ""},
	}
	
	for _, tc := range testCases {
		mapping, exists := getKoreanGenerationType(tc.generatorKey)
		
		if exists != tc.shouldExist {
			t.Errorf("For key %s, expected exists=%v, got %v", 
				tc.generatorKey, tc.shouldExist, exists)
			continue
		}
		
		if exists && mapping.GenerationType != tc.expectedType {
			t.Errorf("For key %s, expected type %v, got %v", 
				tc.generatorKey, tc.expectedType, mapping.GenerationType)
		}
	}
}

func TestIsKoreanWordListMapping(t *testing.T) {
	testCases := []struct {
		generatorKey string
		expected     bool
	}{
		// Korean indicators
		{"korean-common", true},
		{"korean-tech", true},
		{"Korean sentences", true},
		{"some-korean-words", true},
		{"한국어-단어", true},
		{"한글-문장", true},
		{"korea-words", true},
		
		// Non-Korean
		{"english-words", false},
		{"common-words", false},
		{"frankenstein", false},
		{"spanish-words", false},
	}
	
	for _, tc := range testCases {
		result := isKoreanWordList(tc.generatorKey)
		if result != tc.expected {
			t.Errorf("For key %s, expected %v, got %v", 
				tc.generatorKey, tc.expected, result)
		}
	}
}