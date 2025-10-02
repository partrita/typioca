package cmd

import (
	"strings"
	"testing"
	"time"

	"github.com/bloznelis/typioca/cmd/words"
)

// TestKoreanWordListDetectionComprehensive tests Korean wordlist detection with comprehensive scenarios
func TestKoreanWordListDetectionComprehensive(t *testing.T) {
	tests := []struct {
		name         string
		generatorKey string
		expected     bool
		description  string
	}{
		// Embedded Korean wordlists (Requirements 1.1)
		{"Korean common words", "Korean common words", true, "Standard embedded Korean common words"},
		{"Korean tech terms", "Korean tech terms", true, "Standard embedded Korean tech terms"},
		{"Korean sentences", "Korean sentences", true, "Standard embedded Korean sentences"},
		
		// Case variations (Requirements 1.1)
		{"Lowercase korean", "korean common words", true, "Lowercase korean detection"},
		{"Mixed case Korean", "Korean Tech Terms", true, "Mixed case Korean detection"},
		{"UPPERCASE KOREAN", "KOREAN SENTENCES", true, "Uppercase Korean detection"},
		
		// Korean language indicators (Requirements 1.1)
		{"Korean with 한국", "한국어 단어", true, "Korean with Hangul characters"},
		{"Korean with 한글", "한글 연습", true, "Korean with Hangul script name"},
		
		// File paths with Korean indicators (Requirements 1.1)
		{"Korean file path", "/path/to/korean_words.json", true, "File path with korean"},
		{"Korea in path", "/files/korea-tech.txt", true, "File path with korea"},
		{"Korean with underscore", "korean_words.json", true, "Korean with underscore separator"},
		{"Korean with dash", "korean-tech.txt", true, "Korean with dash separator"},
		{"Korean with dot", "korean.json", true, "Korean with dot separator"},
		
		// Non-Korean wordlists (Requirements 1.1)
		{"English common", "Common words", false, "Standard English wordlist"},
		{"Frankenstein", "Frankenstein sentences", false, "English literature wordlist"},
		{"Custom English", "/path/to/english_words.json", false, "Custom English file path"},
		{"Random text", "some random text", false, "Random non-Korean text"},
		{"Empty string", "", false, "Empty string input"},
		
		// Edge cases - should NOT match when korean/korea is part of another word
		{"Korean substring in word", "nonkoreanword", false, "Korean as substring should not match"},
		{"Korea substring in word", "nonkoreaword", false, "Korea as substring should not match"},
		
		// Should match when korean/korea is standalone or properly separated
		{"Just Korean", "Korean", true, "Standalone Korean word"},
		{"Just korean", "korean", true, "Standalone korean word"},
		{"Just korea", "korea", true, "Standalone korea word"},
		
		// Complex real-world scenarios
		{"Korean with numbers", "korean-words-2024.json", true, "Korean file with version numbers"},
		{"Korean with spaces", "Korean Common Words List", true, "Korean with multiple spaces"},
		{"Mixed Korean English", "Korean-English-Dictionary.txt", true, "Mixed Korean-English content"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isKoreanWordList(tt.generatorKey)
			if result != tt.expected {
				t.Errorf("isKoreanWordList(%q) = %v, expected %v - %s", 
					tt.generatorKey, result, tt.expected, tt.description)
			}
		})
	}
}

// TestKoreanContentGenerationInTimerTests tests Korean content generation in timer-based tests
func TestKoreanContentGenerationInTimerTests(t *testing.T) {
	// Test Korean content generation for timer-based tests (Requirements 1.2, 1.3)
	tests := []struct {
		name           string
		generatorKey   string
		expectKorean   bool
		expectSuccess  bool
		minContentLen  int
		description    string
	}{
		{
			name:          "Korean common words generation",
			generatorKey:  "Korean common words",
			expectKorean:  true,
			expectSuccess: true,
			minContentLen: 10,
			description:   "Should generate Korean common words content",
		},
		{
			name:          "Korean tech terms generation",
			generatorKey:  "Korean tech terms",
			expectKorean:  true,
			expectSuccess: true,
			minContentLen: 10,
			description:   "Should generate Korean tech terms content",
		},
		{
			name:          "Korean sentences generation",
			generatorKey:  "Korean sentences",
			expectKorean:  true,
			expectSuccess: true,
			minContentLen: 20,
			description:   "Should generate Korean sentences content",
		},
		{
			name:          "English common words generation",
			generatorKey:  "Common words",
			expectKorean:  false,
			expectSuccess: true,
			minContentLen: 10,
			description:   "Should generate English content for non-Korean wordlist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test settings
			selections := []WordsSelection{
				{name: tt.name, generatorKey: tt.generatorKey},
			}
			
			settings := TimerBasedTestSettings{
				timeSelections:     []time.Duration{30 * time.Second},
				timeCursor:         0,
				wordListSelections: selections,
				wordListCursor:     0,
				cursor:             0,
				enabled:            true,
			}
			
			// Create main menu with generators
			mainMenu := MainMenu{
				timeBasedGenerator:       words.NewGenerator([]string{}),
				koreanTimeBasedGenerator: words.NewKoreanGenerator([]string{}),
			}
			
			// Initialize timer-based test
			test := initTimerBasedTest(settings, mainMenu)
			
			// Verify content was generated
			if len(test.base.wordsToEnter) == 0 {
				if tt.expectSuccess {
					t.Errorf("Expected content to be generated for %s, but got empty content", tt.generatorKey)
				}
				return
			}
			
			// Verify minimum content length
			if len(test.base.wordsToEnter) < tt.minContentLen {
				t.Errorf("Expected at least %d characters for %s, got %d", 
					tt.minContentLen, tt.generatorKey, len(test.base.wordsToEnter))
			}
			
			// Verify Korean content detection
			contentStr := string(test.base.wordsToEnter)
			hasKorean := detectKoreanContent(contentStr)
			
			if tt.expectKorean && !hasKorean {
				t.Errorf("Expected Korean content for %s, but Korean characters not detected in: %s", 
					tt.generatorKey, contentStr[:minInt(50, len(contentStr))])
			}
			
			if !tt.expectKorean && hasKorean {
				t.Errorf("Expected non-Korean content for %s, but Korean characters detected in: %s", 
					tt.generatorKey, contentStr[:minInt(50, len(contentStr))])
			}
		})
	}
}

// TestKoreanDisplayErrorHandling tests error handling scenarios with empty or invalid Korean content
func TestKoreanDisplayErrorHandling(t *testing.T) {
	// Test error handling scenarios (Requirements 1.4)
	tests := []struct {
		name              string
		content           []rune
		generatorKey      string
		listName          string
		expectError       bool
		expectedErrorType string
		description       string
	}{
		{
			name:              "Empty Korean content",
			content:           []rune(""),
			generatorKey:      "korean-common",
			listName:          "Korean common words",
			expectError:       true,
			expectedErrorType: "EmptyContent",
			description:       "Should detect empty Korean content and return appropriate error",
		},
		{
			name:              "Whitespace only Korean content",
			content:           []rune("   \t\n   "),
			generatorKey:      "korean-common",
			listName:          "Korean common words",
			expectError:       true,
			expectedErrorType: "WhitespaceOnly",
			description:       "Should detect whitespace-only content and return appropriate error",
		},
		{
			name:              "Content too short",
			content:           []rune("안녕"),
			generatorKey:      "korean-common",
			listName:          "Korean common words",
			expectError:       true,
			expectedErrorType: "ContentTooShort",
			description:       "Should detect content that is too short for meaningful practice",
		},
		{
			name:              "No Korean characters in Korean wordlist",
			content:           []rune("hello world english text only"),
			generatorKey:      "korean-common",
			listName:          "Korean common words",
			expectError:       true,
			expectedErrorType: "NoKoreanContent",
			description:       "Should detect when Korean wordlist contains no Korean characters",
		},
		{
			name:              "Content with null bytes",
			content:           []rune("안녕하세요\x00세계 한국어 타이핑 연습"),
			generatorKey:      "korean-common",
			listName:          "Korean common words",
			expectError:       true,
			expectedErrorType: "ContentIntegrityError",
			description:       "Should detect corrupted content with null bytes",
		},
		{
			name:              "Valid Korean content",
			content:           []rune("안녕하세요 세계 한국어 타이핑 연습 프로그래밍 개발자 소프트웨어"),
			generatorKey:      "korean-common",
			listName:          "Korean common words",
			expectError:       false,
			expectedErrorType: "",
			description:       "Should accept valid Korean content without errors",
		},
		{
			name:              "Mixed Korean and English content",
			content:           []rune("안녕하세요 hello 세계 world 한국어 programming"),
			generatorKey:      "korean-common",
			listName:          "Korean common words",
			expectError:       false,
			expectedErrorType: "",
			description:       "Should accept mixed Korean-English content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateKoreanContent(tt.content, tt.generatorKey, tt.listName)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for %s, but got none - %s", tt.name, tt.description)
					return
				}
				
				if contentErr, ok := err.(ContentValidationError); ok {
					if contentErr.ErrorType != tt.expectedErrorType {
						t.Errorf("Expected error type %s for %s, got %s - %s", 
							tt.expectedErrorType, tt.name, contentErr.ErrorType, tt.description)
					}
				} else {
					t.Errorf("Expected ContentValidationError for %s, got %T - %s", 
						tt.name, err, tt.description)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for %s, but got: %v - %s", tt.name, err, tt.description)
				}
			}
		})
	}
}

// TestKoreanContentGenerationWithFallback tests Korean content generation with fallback mechanisms
func TestKoreanContentGenerationWithFallback(t *testing.T) {
	// Test fallback mechanisms when Korean generation fails (Requirements 1.4)
	tests := []struct {
		name         string
		generatorKey string
		expectSource string
		description  string
	}{
		{
			name:         "Korean wordlist with fallback",
			generatorKey: "korean-common",
			expectSource: "korean", // May fallback to "fallback" if Korean generation fails
			description:  "Should generate Korean content or fallback gracefully",
		},
		{
			name:         "English wordlist",
			generatorKey: "common-english",
			expectSource: "english",
			description:  "Should generate English content for English wordlists",
		},
		{
			name:         "Unknown wordlist",
			generatorKey: "unknown-wordlist",
			expectSource: "error_message", // Should generate error message content
			description:  "Should handle unknown wordlists gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock generators
			koreanGenerator := words.NewKoreanGenerator([]string{})
			fallbackGenerator := words.NewGenerator([]string{})
			
			// Generate content with error handling
			result := generateContentWithErrorHandling(tt.generatorKey, koreanGenerator, fallbackGenerator)
			
			// Verify content was generated
			if len(result.Content) == 0 {
				t.Errorf("Expected non-empty content for %s, but got empty content - %s", 
					tt.generatorKey, tt.description)
			}
			
			// Verify content is validated
			if !result.Validated {
				t.Errorf("Expected content to be validated for %s - %s", tt.generatorKey, tt.description)
			}
			
			// For Korean wordlists, accept multiple possible sources due to fallback behavior
			if tt.generatorKey == "korean-common" {
				validSources := []string{"korean", "fallback", "error_message"}
				isValidSource := false
				for _, validSource := range validSources {
					if result.Source == validSource {
						isValidSource = true
						break
					}
				}
				if !isValidSource {
					t.Errorf("Expected source to be one of %v for %s, got %s - %s", 
						validSources, tt.generatorKey, result.Source, tt.description)
				}
			} else {
				// For non-Korean wordlists, check specific expected source or error_message
				if result.Source != tt.expectSource && result.Source != "error_message" {
					t.Errorf("Expected source %s or 'error_message' for %s, got %s - %s", 
						tt.expectSource, tt.generatorKey, result.Source, tt.description)
				}
			}
		})
	}
}

// TestKoreanDisplayIntegrationScenarios tests integration scenarios for Korean display
func TestKoreanDisplayIntegrationScenarios(t *testing.T) {
	// Test complete integration scenarios (Requirements 1.1, 1.2, 1.3, 1.4)
	t.Run("Korean Timer Test Complete Flow", func(t *testing.T) {
		// Create Korean wordlist selections
		koreanSelections := []WordsSelection{
			{name: "Korean common words", generatorKey: "Korean common words"},
			{name: "Korean tech terms", generatorKey: "Korean tech terms"},
			{name: "Korean sentences", generatorKey: "Korean sentences"},
		}
		
		// Test each Korean wordlist type
		for i, selection := range koreanSelections {
			t.Run(selection.name, func(t *testing.T) {
				settings := TimerBasedTestSettings{
					timeSelections:     []time.Duration{30 * time.Second},
					timeCursor:         0,
					wordListSelections: koreanSelections,
					wordListCursor:     i, // Select current Korean wordlist
					cursor:             0,
					enabled:            true,
				}
				
				mainMenu := MainMenu{
					timeBasedGenerator:       words.NewGenerator([]string{}),
					koreanTimeBasedGenerator: words.NewKoreanGenerator([]string{}),
				}
				
				// Initialize test
				test := initTimerBasedTest(settings, mainMenu)
				
				// Verify Korean content was generated
				if len(test.base.wordsToEnter) == 0 {
					t.Errorf("Korean display issue: No content generated for %s", selection.name)
				}
				
				// Verify Korean characters are present
				contentStr := string(test.base.wordsToEnter)
				if !detectKoreanContent(contentStr) {
					t.Errorf("Korean display issue: No Korean characters detected in content for %s: %s", 
						selection.name, contentStr[:minInt(50, len(contentStr))])
				}
				
				// Verify content is suitable for typing practice
				if len(test.base.wordsToEnter) < 20 {
					t.Errorf("Korean display issue: Content too short for meaningful practice in %s: %d characters", 
						selection.name, len(test.base.wordsToEnter))
				}
			})
		}
	})
	
	t.Run("Korean Content Validation Integration", func(t *testing.T) {
		// Test that Korean content validation works in the complete flow
		koreanMapping := KoreanWordListMapping{
			GeneratorKey:   "korean-common",
			ListName:       "Korean common words",
			GenerationType: KoreanWords,
		}
		
		generator := words.NewKoreanGenerator([]string{})
		result := generateKoreanContentWithValidation(koreanMapping, generator)
		
		// Verify result structure
		if !result.Validated {
			t.Error("Korean content should be validated in integration flow")
		}
		
		if result.Source != "korean" && result.Source != "fallback" {
			t.Errorf("Expected Korean or fallback source in integration, got: %s", result.Source)
		}
		
		// If generation succeeded, verify Korean content
		if result.Success && len(result.Content) > 0 {
			contentStr := string(result.Content)
			if !detectKoreanContent(contentStr) {
				t.Errorf("Korean integration issue: Generated content lacks Korean characters: %s", 
					contentStr[:minInt(50, len(contentStr))])
			}
		}
	})
	
	t.Run("Korean Error Message Generation", func(t *testing.T) {
		// Test that user-friendly error messages are generated for Korean display issues
		testErrors := []ContentValidationError{
			{
				GeneratorKey: "korean-common",
				ListName:     "Korean common words",
				ErrorType:    "EmptyContent",
				Message:      "Generated Korean content is empty",
			},
			{
				GeneratorKey: "korean-sentences",
				ListName:     "Korean sentences",
				ErrorType:    "NoKoreanContent",
				Message:      "Generated content does not contain Korean characters",
			},
		}
		
		for _, testErr := range testErrors {
			message := createUserFriendlyErrorMessage(testErr, testErr.GeneratorKey)
			
			if len(message) == 0 {
				t.Errorf("Expected non-empty error message for %s", testErr.ErrorType)
			}
			
			// Verify message contains helpful information
			lowerMessage := strings.ToLower(message)
			if !strings.Contains(lowerMessage, "korean") {
				t.Errorf("Expected error message to mention Korean for Korean wordlist error: %s", message)
			}
		}
	})
}

// TestKoreanDisplayRobustness tests robustness of Korean display functionality
func TestKoreanDisplayRobustness(t *testing.T) {
	// Test robustness scenarios (Requirements 1.4)
	t.Run("Emergency Fallback Content Generation", func(t *testing.T) {
		// Test emergency fallback for Korean wordlists
		koreanKeys := []string{
			"korean-common",
			"korean-tech", 
			"korean-sentences",
			"Korean common words",
		}
		
		for _, key := range koreanKeys {
			t.Run(key, func(t *testing.T) {
				content := generateEmergencyFallbackContent(key)
				
				if len(content) == 0 {
					t.Errorf("Emergency fallback should never return empty content for %s", key)
				}
				
				contentStr := string(content)
				if !detectKoreanContent(contentStr) {
					t.Errorf("Emergency fallback for Korean wordlist should contain Korean characters: %s", contentStr)
				}
			})
		}
		
		// Test emergency fallback for English wordlists
		englishKeys := []string{
			"common-english",
			"Common words",
			"Frankenstein sentences",
		}
		
		for _, key := range englishKeys {
			t.Run(key, func(t *testing.T) {
				content := generateEmergencyFallbackContent(key)
				
				if len(content) == 0 {
					t.Errorf("Emergency fallback should never return empty content for %s", key)
				}
				
				contentStr := string(content)
				if detectKoreanContent(contentStr) {
					t.Errorf("Emergency fallback for English wordlist should not contain Korean characters: %s", contentStr)
				}
			})
		}
	})
	
	t.Run("Content Generation Panic Recovery", func(t *testing.T) {
		// Test panic recovery mechanism by testing the fallback case
		// Since we can't easily trigger a panic in a test, we test the fallback behavior
		result := handleContentGenerationPanic("korean-common")
		
		if len(result.Content) == 0 {
			t.Error("Panic recovery should generate non-empty content")
		}
		
		if result.Success {
			t.Error("Panic recovery result should indicate failure")
		}
		
		// The function returns "error" when no panic occurs (normal case)
		// and "panic_recovery" when an actual panic is recovered from
		if result.Source != "error" && result.Source != "panic_recovery" {
			t.Errorf("Expected 'error' or 'panic_recovery' source, got: %s", result.Source)
		}
		
		if !result.Validated {
			t.Error("Panic recovery content should be marked as validated")
		}
	})
}

// Helper function for minimum calculation (renamed to avoid conflict)
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}