package words

import (
	"strings"
	"testing"
)

func TestAddEmbeddedKoreanSources(t *testing.T) {
	sources := make(map[string]WordSource)
	result := addEmbeddedKoreanSources(sources)

	// Check that Korean sources were added
	expectedSources := []string{
		"Korean common words",
		"Korean tech terms", 
		"Korean sentences",
	}

	for _, sourceName := range expectedSources {
		if _, exists := result[sourceName]; !exists {
			t.Errorf("Expected Korean source '%s' not found", sourceName)
		}
	}

	// Verify Korean common words source has content
	koreanCommonSource := result["Korean common words"]
	if len(koreanCommonSource.Words) == 0 {
		t.Error("Korean common words source should have words")
	}

	// Verify Korean tech terms source has content
	koreanTechSource := result["Korean tech terms"]
	if len(koreanTechSource.Words) == 0 {
		t.Error("Korean tech terms source should have words")
	}

	// Verify Korean sentences source has content
	koreanSentenceSource := result["Korean sentences"]
	if len(koreanSentenceSource.Words) == 0 {
		t.Error("Korean sentences source should have sentences")
	}

	// Verify metadata
	if koreanCommonSource.Metadata.Name == "" {
		t.Error("Korean common words should have metadata name")
	}
}

func TestNewKoreanGenerator(t *testing.T) {
	paths := []string{} // Empty paths to test embedded sources only
	generator := NewKoreanGenerator(paths)

	// Check that generator has Korean sources
	expectedSources := []string{
		"Korean common words",
		"Korean tech terms",
		"Korean sentences",
	}

	for _, sourceName := range expectedSources {
		if _, exists := generator.poolsJson[sourceName]; !exists {
			t.Errorf("Korean generator should have source '%s'", sourceName)
		}
	}

	// Check default count
	if generator.Count != 300 {
		t.Errorf("Korean generator default count should be 300, got %d", generator.Count)
	}
}

func TestGenerateKoreanWords(t *testing.T) {
	generator := NewKoreanGenerator([]string{})
	
	// Test generating Korean common words
	words := generator.GenerateKoreanWords("Korean common words")
	
	if len(words) == 0 {
		t.Error("GenerateKoreanWords should return non-empty result")
	}

	// Convert to string and check for Korean content
	wordsStr := string(words)
	hasKorean := false
	for _, r := range words {
		if r >= 0xAC00 && r <= 0xD7AF { // Korean syllables
			hasKorean = true
			break
		}
	}

	if !hasKorean {
		t.Errorf("Generated words should contain Korean characters, got: %s", wordsStr)
	}
}

func TestGenerateKoreanSentences(t *testing.T) {
	generator := NewKoreanGenerator([]string{})
	
	// Test generating Korean sentences
	sentences := generator.GenerateKoreanSentences("Korean sentences")
	
	if len(sentences) == 0 {
		t.Error("GenerateKoreanSentences should return non-empty result")
	}

	// Convert to string and check for Korean content
	sentencesStr := string(sentences)
	hasKorean := false
	for _, r := range sentences {
		if r >= 0xAC00 && r <= 0xD7AF { // Korean syllables
			hasKorean = true
			break
		}
	}

	if !hasKorean {
		t.Errorf("Generated sentences should contain Korean characters, got: %s", sentencesStr)
	}

	// Check that sentences are properly spaced (should contain spaces between sentences)
	if !strings.Contains(sentencesStr, " ") {
		t.Error("Generated sentences should be properly spaced")
	}
}

func TestShuffleKoreanWords(t *testing.T) {
	generator := NewKoreanGenerator([]string{})
	
	words := []string{"안녕", "하세요", "프로그래밍", "데이터베이스", "알고리즘"}
	
	// Test multiple shuffles to ensure randomness
	shuffled1 := generator.shuffleKoreanWords(words)
	shuffled2 := generator.shuffleKoreanWords(words)
	
	// Check that original slice is not modified
	originalWords := []string{"안녕", "하세요", "프로그래밍", "데이터베이스", "알고리즘"}
	for i, word := range words {
		if word != originalWords[i] {
			t.Error("Original words slice should not be modified")
		}
	}
	
	// Check that shuffled slices have same length
	if len(shuffled1) != len(words) || len(shuffled2) != len(words) {
		t.Error("Shuffled slices should have same length as original")
	}
	
	// Check that all original words are present in shuffled slices
	for _, originalWord := range words {
		found1 := false
		found2 := false
		for _, shuffledWord := range shuffled1 {
			if shuffledWord == originalWord {
				found1 = true
				break
			}
		}
		for _, shuffledWord := range shuffled2 {
			if shuffledWord == originalWord {
				found2 = true
				break
			}
		}
		if !found1 || !found2 {
			t.Errorf("Word '%s' should be present in shuffled slices", originalWord)
		}
	}
}

func TestGenerateKoreanWithCount(t *testing.T) {
	generator := NewKoreanGenerator([]string{})
	
	// Test with specific count
	count := 5
	words := generator.GenerateKoreanWithCount("Korean common words", count)
	
	if len(words) == 0 {
		t.Error("GenerateKoreanWithCount should return non-empty result")
	}

	// Convert to string and count words (split by spaces)
	wordsStr := string(words)
	wordList := strings.Fields(wordsStr)
	
	// Should respect the count limit (allowing some flexibility for sentence boundaries)
	if len(wordList) > count*2 { // Allow some flexibility
		t.Errorf("Generated word count should be around %d, got %d words", count, len(wordList))
	}

	// Check for Korean content
	hasKorean := false
	for _, r := range words {
		if r >= 0xAC00 && r <= 0xD7AF { // Korean syllables
			hasKorean = true
			break
		}
	}

	if !hasKorean {
		t.Errorf("Generated words should contain Korean characters, got: %s", wordsStr)
	}
}

// Helper function to check if string contains Korean characters
func containsKoreanChars(s string) bool {
	for _, r := range s {
		if r >= 0xAC00 && r <= 0xD7AF { // Hangul syllables range
			return true
		}
	}
	return false
}