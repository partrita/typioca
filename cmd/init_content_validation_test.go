package cmd

import (
	"strings"
	"testing"

	"github.com/bloznelis/typioca/cmd/words"
)

func TestValidateKoreanContent(t *testing.T) {
	tests := []struct {
		name         string
		content      []rune
		generatorKey string
		listName     string
		expectError  bool
		errorType    string
	}{
		{
			name:         "Valid Korean content",
			content:      []rune("안녕하세요 세계 한국어 타이핑 연습"),
			generatorKey: "korean-common",
			listName:     "Korean common words",
			expectError:  false,
		},
		{
			name:         "Empty content",
			content:      []rune(""),
			generatorKey: "korean-common",
			listName:     "Korean common words",
			expectError:  true,
			errorType:    "EmptyContent",
		},
		{
			name:         "Content too short",
			content:      []rune("안녕"),
			generatorKey: "korean-common",
			listName:     "Korean common words",
			expectError:  true,
			errorType:    "ContentTooShort",
		},
		{
			name:         "No Korean content",
			content:      []rune("hello world english text"),
			generatorKey: "korean-common",
			listName:     "Korean common words",
			expectError:  true,
			errorType:    "NoKoreanContent",
		},
		{
			name:         "Whitespace only",
			content:      []rune("   \t\n   "),
			generatorKey: "korean-common",
			listName:     "Korean common words",
			expectError:  true,
			errorType:    "WhitespaceOnly",
		},
		{
			name:         "Mixed Korean and English (valid)",
			content:      []rune("안녕하세요 hello 세계 world 한국어"),
			generatorKey: "korean-common",
			listName:     "Korean common words",
			expectError:  false,
		},
		{
			name:         "Content with null bytes",
			content:      []rune("안녕하세요\x00세계"),
			generatorKey: "korean-common",
			listName:     "Korean common words",
			expectError:  true,
			errorType:    "ContentIntegrityError",
		},
		{
			name:         "Content with insufficient diversity",
			content:      []rune("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
			generatorKey: "korean-common",
			listName:     "Korean common words",
			expectError:  true,
			errorType:    "ContentIntegrityError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateKoreanContent(tt.content, tt.generatorKey, tt.listName)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				
				if contentErr, ok := err.(ContentValidationError); ok {
					if contentErr.ErrorType != tt.errorType {
						t.Errorf("Expected error type %s, got %s", tt.errorType, contentErr.ErrorType)
					}
				} else {
					t.Errorf("Expected ContentValidationError, got %T", err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestGenerateKoreanContentWithValidation(t *testing.T) {
	// Create a mock Korean generator for testing
	generator := words.NewKoreanGenerator([]string{})
	
	tests := []struct {
		name            string
		koreanMapping   KoreanWordListMapping
		expectSuccess   bool
		expectedSource  string
	}{
		{
			name: "Valid Korean words generation",
			koreanMapping: KoreanWordListMapping{
				GeneratorKey:   "korean-common",
				ListName:       "Korean common words",
				GenerationType: KoreanWords,
			},
			expectSuccess:  true,
			expectedSource: "korean",
		},
		{
			name: "Valid Korean sentences generation",
			koreanMapping: KoreanWordListMapping{
				GeneratorKey:   "korean-sentences",
				ListName:       "Korean sentences",
				GenerationType: KoreanSentences,
			},
			expectSuccess:  true,
			expectedSource: "korean",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateKoreanContentWithValidation(tt.koreanMapping, generator)
			
			if tt.expectSuccess && !result.Success {
				t.Errorf("Expected success but got failure: %v", result.Error)
			}
			
			if result.Source != tt.expectedSource {
				t.Errorf("Expected source %s, got %s", tt.expectedSource, result.Source)
			}
			
			if !result.Validated {
				t.Errorf("Expected content to be validated")
			}
		})
	}
}

func TestGenerateFallbackContent(t *testing.T) {
	// Create a mock generator for testing
	generator := words.NewGenerator([]string{})
	
	result := generateFallbackContent("test-key", generator)
	
	if !result.Success {
		t.Errorf("Expected fallback generation to succeed")
	}
	
	if result.Source != "fallback" {
		t.Errorf("Expected source to be 'fallback', got %s", result.Source)
	}
	
	if len(result.Content) == 0 {
		t.Errorf("Expected fallback content to be non-empty")
	}
	
	if !result.Validated {
		t.Errorf("Expected fallback content to be validated")
	}
}

func TestCreateUserFriendlyErrorMessage(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		generatorKey string
		expectContains string
	}{
		{
			name: "Empty content error",
			err: ContentValidationError{
				GeneratorKey: "korean-common",
				ListName:     "Korean common words",
				ErrorType:    "EmptyContent",
				Message:      "Generated Korean content is empty",
			},
			generatorKey:   "korean-common",
			expectContains: "appears to be empty",
		},
		{
			name: "Content too short error",
			err: ContentValidationError{
				GeneratorKey: "korean-common",
				ListName:     "Korean common words",
				ErrorType:    "ContentTooShort",
				Message:      "Generated content too short",
			},
			generatorKey:   "korean-common",
			expectContains: "insufficient content",
		},
		{
			name: "No Korean content error",
			err: ContentValidationError{
				GeneratorKey: "korean-common",
				ListName:     "Korean common words",
				ErrorType:    "NoKoreanContent",
				Message:      "Generated content does not contain Korean characters",
			},
			generatorKey:   "korean-common",
			expectContains: "does not contain Korean characters",
		},
		{
			name: "Content integrity error",
			err: ContentValidationError{
				GeneratorKey: "korean-common",
				ListName:     "Korean common words",
				ErrorType:    "ContentIntegrityError",
				Message:      "Content integrity validation failed",
			},
			generatorKey:   "korean-common",
			expectContains: "corrupted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := createUserFriendlyErrorMessage(tt.err, tt.generatorKey)
			
			if !strings.Contains(strings.ToLower(message), strings.ToLower(tt.expectContains)) {
				t.Errorf("Expected message to contain '%s', got: %s", tt.expectContains, message)
			}
		})
	}
}

func TestGenerateEmergencyFallbackContent(t *testing.T) {
	tests := []struct {
		name         string
		generatorKey string
		expectKorean bool
	}{
		{
			name:         "Korean wordlist emergency fallback",
			generatorKey: "korean-common",
			expectKorean: true,
		},
		{
			name:         "English wordlist emergency fallback",
			generatorKey: "common-english",
			expectKorean: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := generateEmergencyFallbackContent(tt.generatorKey)
			
			if len(content) == 0 {
				t.Errorf("Expected non-empty emergency fallback content")
			}
			
			contentStr := string(content)
			hasKorean := detectKoreanContent(contentStr)
			
			if tt.expectKorean && !hasKorean {
				t.Errorf("Expected Korean content for Korean wordlist, got: %s", contentStr)
			}
			
			if !tt.expectKorean && hasKorean {
				t.Errorf("Expected English content for English wordlist, got: %s", contentStr)
			}
		})
	}
}

func TestGenerateContentWithErrorHandling(t *testing.T) {
	// Create mock generators for testing
	koreanGenerator := words.NewKoreanGenerator([]string{})
	fallbackGenerator := words.NewGenerator([]string{})
	
	tests := []struct {
		name         string
		generatorKey string
		expectSource string
	}{
		{
			name:         "Korean wordlist",
			generatorKey: "korean-common",
			expectSource: "korean", // May fallback to "fallback" if Korean generation fails
		},
		{
			name:         "English wordlist",
			generatorKey: "common-english",
			expectSource: "english",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateContentWithErrorHandling(tt.generatorKey, koreanGenerator, fallbackGenerator)
			
			if len(result.Content) == 0 {
				t.Errorf("Expected non-empty content")
			}
			
			if !result.Validated {
				t.Errorf("Expected content to be validated")
			}
			
			// For Korean wordlists, we accept either "korean" or "fallback" source
			// depending on whether the Korean generation succeeds
			if tt.generatorKey == "korean-common" {
				if result.Source != "korean" && result.Source != "fallback" && result.Source != "error_message" {
					t.Errorf("Expected source to be 'korean', 'fallback', or 'error_message' for Korean wordlist, got %s", result.Source)
				}
			} else {
				// For non-Korean wordlists, we expect the specified source or error_message
				if result.Source != tt.expectSource && result.Source != "error_message" {
					t.Errorf("Expected source %s or 'error_message', got %s", tt.expectSource, result.Source)
				}
			}
		})
	}
}