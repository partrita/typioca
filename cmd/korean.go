package cmd

import (
	"fmt"
	"log"
	"time"
	"unicode"
)

// isKoreanChar checks if a rune is a Korean character (Hangul)
func isKoreanChar(r rune) bool {
	return unicode.Is(unicode.Hangul, r)
}

// calculateKoreanWPM calculates Words Per Minute for Korean text
// Korean characters are typically counted differently than English
// Each Korean syllable is considered equivalent to 2-3 English characters
func calculateKoreanWPM(chars int, duration time.Duration) int {
	if duration.Minutes() == 0 {
		return 0
	}
	
	// Adjust Korean characters to equivalent English character count
	// Korean syllables are more complex than English letters
	adjustedChars := float64(chars) * 2.5
	
	// Standard WPM calculation: (characters / 5) / minutes
	// The division by 5 is the standard word length assumption
	wpm := adjustedChars / 5.0 / duration.Minutes()
	
	return int(wpm)
}

// calculateKoreanAccuracy calculates typing accuracy for Korean text
// Compares input runes with target runes and returns accuracy percentage
func calculateKoreanAccuracy(input, target []rune) float64 {
	if len(target) == 0 {
		return 100.0
	}
	
	correctChars := 0
	minLength := len(input)
	if len(target) < minLength {
		minLength = len(target)
	}
	
	// Count correct characters up to the minimum length
	for i := 0; i < minLength; i++ {
		if input[i] == target[i] {
			correctChars++
		}
	}
	
	// Calculate accuracy based on target length
	// This accounts for incomplete input (shorter than target)
	accuracy := float64(correctChars) / float64(len(target)) * 100.0
	
	// Ensure accuracy doesn't exceed 100%
	if accuracy > 100.0 {
		accuracy = 100.0
	}
	
	return accuracy
}

// isKoreanSyllable checks if a rune is a complete Korean syllable
// Korean syllables are in the range U+AC00 to U+D7AF
func isKoreanSyllable(r rune) bool {
	return r >= 0xAC00 && r <= 0xD7AF
}

// isKoreanJamo checks if a rune is a Korean Jamo (component of syllables)
// Jamo includes initial consonants, vowels, and final consonants
func isKoreanJamo(r rune) bool {
	// Hangul Jamo (U+1100–U+11FF)
	// Hangul Compatibility Jamo (U+3130–U+318F)
	// Hangul Jamo Extended-A (U+A960–U+A97F)
	// Hangul Jamo Extended-B (U+D7B0–U+D7FF)
	return (r >= 0x1100 && r <= 0x11FF) ||
		   (r >= 0x3130 && r <= 0x318F) ||
		   (r >= 0xA960 && r <= 0xA97F) ||
		   (r >= 0xD7B0 && r <= 0xD7FF)
}

// countKoreanSyllables counts the number of Korean syllables in a string
// This is useful for more accurate Korean text metrics
func countKoreanSyllables(text string) int {
	count := 0
	for _, r := range text {
		if isKoreanSyllable(r) {
			count++
		}
	}
	return count
}

// validateKoreanWordList validates that a word list contains proper Korean content
func validateKoreanWordList(words []string) error {
	if len(words) == 0 {
		return nil // Empty list is valid
	}
	
	koreanWordCount := 0
	for _, word := range words {
		if detectKoreanContent(word) {
			koreanWordCount++
		}
	}
	
	// At least 50% of words should contain Korean characters
	koreanRatio := float64(koreanWordCount) / float64(len(words))
	if koreanRatio < 0.5 {
		return nil // Not enough Korean content, but not an error - just not Korean
	}
	
	return nil
}

// detectKoreanContent analyzes text content to determine if it contains Korean
func detectKoreanContent(text string) bool {
	if len(text) == 0 {
		return false
	}
	
	koreanCharCount := 0
	totalCharCount := 0
	
	for _, r := range text {
		if unicode.IsLetter(r) {
			totalCharCount++
			if isKoreanChar(r) {
				koreanCharCount++
			}
		}
	}
	
	// Consider it Korean if more than 30% of letters are Korean characters
	if totalCharCount == 0 {
		return false
	}
	
	koreanRatio := float64(koreanCharCount) / float64(totalCharCount)
	isKorean := koreanRatio > 0.3
	
	// Log Korean content detection for debugging
	if isKorean {
		log.Printf("[KOREAN-DEBUG] Korean content detected: %d/%d chars are Korean (%.1f%%)", 
			koreanCharCount, totalCharCount, koreanRatio*100)
	} else if koreanCharCount > 0 {
		log.Printf("[KOREAN-DEBUG] Some Korean chars found but below threshold: %d/%d chars (%.1f%%)", 
			koreanCharCount, totalCharCount, koreanRatio*100)
	}
	
	return isKorean
}

// validateKoreanContentIntegrity performs additional integrity checks on Korean content
func validateKoreanContentIntegrity(content []rune) error {
	if len(content) == 0 {
		return nil // Empty content is handled elsewhere
	}
	
	contentStr := string(content)
	log.Printf("[KOREAN-DEBUG] Validating Korean content integrity: %d characters", len(content))
	
	// Check for null bytes or other control characters that might indicate corruption
	for i, r := range content {
		if r == 0 {
			log.Printf("[KOREAN-ERROR] Content integrity check failed: null byte at position %d", i)
			return fmt.Errorf("content contains null byte at position %d", i)
		}
		if r < 32 && r != 9 && r != 10 && r != 13 { // Allow tab, newline, carriage return
			log.Printf("[KOREAN-ERROR] Content integrity check failed: invalid control character (code %d) at position %d", r, i)
			return fmt.Errorf("content contains invalid control character (code %d) at position %d", r, i)
		}
	}
	
	// Check for reasonable content structure (not just repeated characters)
	if len(contentStr) > 20 {
		// Count unique characters
		uniqueChars := make(map[rune]bool)
		koreanChars := make(map[rune]bool)
		for _, r := range content {
			if unicode.IsLetter(r) {
				uniqueChars[r] = true
				if isKoreanChar(r) {
					koreanChars[r] = true
				}
			}
		}
		
		log.Printf("[KOREAN-DEBUG] Content diversity analysis: %d unique chars, %d unique Korean chars", 
			len(uniqueChars), len(koreanChars))
		
		// If we have very few unique characters relative to content length, it might be corrupted
		if len(uniqueChars) < 3 && len(content) > 50 {
			log.Printf("[KOREAN-ERROR] Content integrity check failed: insufficient character diversity (%d unique characters)", len(uniqueChars))
			return fmt.Errorf("content appears to have insufficient character diversity (only %d unique characters)", len(uniqueChars))
		}
	}
	
	log.Printf("[KOREAN-DEBUG] Korean content integrity validation passed")
	return nil
}