package cmd

import (
	"fmt"
	"strings"
	"testing"
)

func TestCreateUserFriendlyErrorMessageImplementation(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		generatorKey string
		expectEmoji  bool
		expectTips   bool
	}{
		{
			name: "Empty Korean content error",
			err: ContentValidationError{
				GeneratorKey: "korean-common",
				ListName:     "Korean common words",
				ErrorType:    "EmptyContent",
				Message:      "Generated Korean content is empty",
			},
			generatorKey: "korean-common",
			expectEmoji:  true,
			expectTips:   true,
		},
		{
			name: "No Korean content error",
			err: ContentValidationError{
				GeneratorKey: "korean-tech",
				ListName:     "Korean tech terms",
				ErrorType:    "NoKoreanContent",
				Message:      "Generated content does not contain Korean characters",
			},
			generatorKey: "korean-tech",
			expectEmoji:  true,
			expectTips:   true,
		},
		{
			name: "Content integrity error",
			err: ContentValidationError{
				GeneratorKey: "korean-sentences",
				ListName:     "Korean sentences",
				ErrorType:    "ContentIntegrityError",
				Message:      "Content integrity validation failed",
			},
			generatorKey: "korean-sentences",
			expectEmoji:  true,
			expectTips:   true,
		},
		{
			name: "Generic Korean error",
			err:  fmt.Errorf("Korean generator failed"),
			generatorKey: "korean-common",
			expectEmoji:  true,
			expectTips:   true,
		},
		{
			name: "Non-Korean error",
			err:  fmt.Errorf("Generator failed"),
			generatorKey: "common-words",
			expectEmoji:  true,
			expectTips:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := createUserFriendlyErrorMessage(tt.err, tt.generatorKey)
			
			// Check that the message is not empty
			if len(result) == 0 {
				t.Errorf("Expected non-empty error message, got empty string")
			}
			
			// Check for emoji if expected
			if tt.expectEmoji {
				if !strings.Contains(result, "❌") && !strings.Contains(result, "💡") {
					t.Errorf("Expected error message to contain emoji, got: %s", result)
				}
			}
			
			// Check for helpful tips if expected
			if tt.expectTips {
				if !strings.Contains(result, "💡") {
					t.Errorf("Expected error message to contain helpful tips (💡), got: %s", result)
				}
			}
			
			// Check that generator key is mentioned
			if !strings.Contains(result, tt.generatorKey) {
				t.Errorf("Expected error message to mention generator key '%s', got: %s", tt.generatorKey, result)
			}
		})
	}
}

func TestCreateTerminalEncodingWarning(t *testing.T) {
	generatorKey := "korean-common"
	warning := createTerminalEncodingWarning(generatorKey)
	
	// Check that warning is not empty
	if len(warning) == 0 {
		t.Errorf("Expected non-empty warning message, got empty string")
	}
	
	// Check for key components
	expectedComponents := []string{
		"⚠️",  // Warning emoji
		"Terminal Font Issues",
		"Terminal Encoding Issues", 
		"UTF-8",
		"Korean characters",
		"Recommended Terminal Fonts",
		"Quick Fixes",
		generatorKey,
	}
	
	for _, component := range expectedComponents {
		if !strings.Contains(warning, component) {
			t.Errorf("Expected warning to contain '%s', got: %s", component, warning)
		}
	}
	
	// Check for specific font recommendations
	fontRecommendations := []string{
		"Noto Sans CJK",
		"Source Han Sans",
		"Malgun Gothic",
		"AppleGothic",
	}
	
	foundFontRecommendation := false
	for _, font := range fontRecommendations {
		if strings.Contains(warning, font) {
			foundFontRecommendation = true
			break
		}
	}
	
	if !foundFontRecommendation {
		t.Errorf("Expected warning to contain at least one font recommendation, got: %s", warning)
	}
}

func TestCreateKoreanWordlistUnavailableGuidance(t *testing.T) {
	generatorKey := "korean-tech"
	guidance := createKoreanWordlistUnavailableGuidance(generatorKey)
	
	// Check that guidance is not empty
	if len(guidance) == 0 {
		t.Errorf("Expected non-empty guidance message, got empty string")
	}
	
	// Check for key components
	expectedComponents := []string{
		"📋", // List emoji
		"Possible Causes",
		"Solutions",
		"Alternative Options",
		"Expected Korean Wordlists",
		generatorKey,
		"korean-common",
		"korean-tech", 
		"korean-sentences",
	}
	
	for _, component := range expectedComponents {
		if !strings.Contains(guidance, component) {
			t.Errorf("Expected guidance to contain '%s', got: %s", component, guidance)
		}
	}
}

func TestDisplayUserFriendlyError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		generatorKey string
		expectLong   bool // Expect comprehensive error with additional guidance
	}{
		{
			name: "Empty Korean content with guidance",
			err: ContentValidationError{
				GeneratorKey: "korean-common",
				ListName:     "Korean common words",
				ErrorType:    "EmptyContent",
				Message:      "Generated Korean content is empty",
			},
			generatorKey: "korean-common",
			expectLong:   true,
		},
		{
			name: "Content integrity error with terminal warning",
			err: ContentValidationError{
				GeneratorKey: "korean-sentences",
				ListName:     "Korean sentences",
				ErrorType:    "ContentIntegrityError",
				Message:      "Content integrity validation failed",
			},
			generatorKey: "korean-sentences",
			expectLong:   true,
		},
		{
			name: "Generic Korean error with terminal warning",
			err:  fmt.Errorf("Korean generator initialization failed"),
			generatorKey: "korean-tech",
			expectLong:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := displayUserFriendlyError(tt.err, tt.generatorKey)
			
			// Check that the message is not empty
			if len(result) == 0 {
				t.Errorf("Expected non-empty error display, got empty string")
			}
			
			// Check for comprehensive content if expected
			if tt.expectLong {
				// Should contain both error message and additional guidance
				if len(result) < 200 { // Arbitrary threshold for "comprehensive"
					t.Errorf("Expected comprehensive error message (>200 chars), got %d chars: %s", len(result), result)
				}
				
				// Should contain multiple sections
				sectionCount := 0
				if strings.Contains(result, "❌") {
					sectionCount++
				}
				if strings.Contains(result, "⚠️") || strings.Contains(result, "📋") {
					sectionCount++
				}
				
				if sectionCount < 2 {
					t.Errorf("Expected multiple guidance sections, found %d sections in: %s", sectionCount, result)
				}
			}
		})
	}
}

func TestValidateTerminalKoreanSupport(t *testing.T) {
	supported, warnings := validateTerminalKoreanSupport()
	
	// This test depends on the actual terminal environment
	// We'll just verify the function returns reasonable values
	
	if !supported && len(warnings) == 0 {
		t.Errorf("If Korean support is not detected, warnings should be provided")
	}
	
	if supported && len(warnings) > 0 {
		// This is okay - warnings can be provided even if basic support is detected
		t.Logf("Korean support detected but with warnings: %v", warnings)
	}
	
	// Test that the function doesn't panic and returns valid data
	if len(warnings) > 10 {
		t.Errorf("Too many warnings returned (%d), might indicate an issue", len(warnings))
	}
}

func TestCreateKoreanSetupGuidance(t *testing.T) {
	guidance := createKoreanSetupGuidance()
	
	// Check that guidance is not empty
	if len(guidance) == 0 {
		t.Errorf("Expected non-empty setup guidance, got empty string")
	}
	
	// Check for key setup components
	expectedComponents := []string{
		"🇰🇷", // Korean flag emoji
		"Required Components",
		"Setup Steps",
		"Install Korean Fonts",
		"Configure Terminal",
		"Enable Korean Wordlists",
		"Troubleshooting",
		"UTF-8",
		"한글", // Korean text test
	}
	
	for _, component := range expectedComponents {
		if !strings.Contains(guidance, component) {
			t.Errorf("Expected setup guidance to contain '%s', got: %s", component, guidance)
		}
	}
	
	// Check for platform-specific instructions
	platforms := []string{"Linux", "macOS", "Windows"}
	foundPlatform := false
	for _, platform := range platforms {
		if strings.Contains(guidance, platform) {
			foundPlatform = true
			break
		}
	}
	
	if !foundPlatform {
		t.Errorf("Expected setup guidance to contain platform-specific instructions, got: %s", guidance)
	}
}

func TestLogUserFriendlyError(t *testing.T) {
	// This test verifies that the logging function doesn't panic
	// and handles different error types appropriately
	
	tests := []struct {
		name         string
		err          error
		generatorKey string
	}{
		{
			name: "Korean content validation error",
			err: ContentValidationError{
				GeneratorKey: "korean-common",
				ListName:     "Korean common words",
				ErrorType:    "EmptyContent",
				Message:      "Generated Korean content is empty",
			},
			generatorKey: "korean-common",
		},
		{
			name:         "Generic Korean error",
			err:          fmt.Errorf("Korean generator failed"),
			generatorKey: "korean-tech",
		},
		{
			name:         "Non-Korean error",
			err:          fmt.Errorf("Generator failed"),
			generatorKey: "common-words",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should not panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("logUserFriendlyError panicked: %v", r)
				}
			}()
			
			logUserFriendlyError(tt.err, tt.generatorKey)
		})
	}
}