package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bloznelis/typioca/cmd/words"
)

func TestKoreanCustomFileLoading(t *testing.T) {
	// Create temporary directory for test files
	tempDir, err := os.MkdirTemp("", "korean_custom_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test Korean JSON file loading
	t.Run("Korean JSON file detection", func(t *testing.T) {
		// Create a Korean JSON word list file
		koreanWords := words.WordSource{
			Metadata: words.Metadata{
				Name:       "Test Korean Words",
				Size:       5,
				PackagedAt: time.Now().Format(time.RFC3339),
				Version:    1,
			},
			Words: []string{"안녕하세요", "반갑습니다", "프로그래밍", "컴퓨터", "한국어"},
		}

		jsonFile := filepath.Join(tempDir, "korean_words.json")
		file, err := os.Create(jsonFile)
		if err != nil {
			t.Fatalf("Failed to create JSON file: %v", err)
		}

		encoder := json.NewEncoder(file)
		err = encoder.Encode(koreanWords)
		file.Close()
		if err != nil {
			t.Fatalf("Failed to write JSON file: %v", err)
		}

		// Test language detection
		detectedLanguage := detectWordListLanguage(jsonFile)
		if detectedLanguage != "korean" {
			t.Errorf("Expected language 'korean', got '%s'", detectedLanguage)
		}
	})

	// Test Korean text file loading
	t.Run("Korean text file detection", func(t *testing.T) {
		// Create a Korean text word list file
		textFile := filepath.Join(tempDir, "korean_words.txt")
		file, err := os.Create(textFile)
		if err != nil {
			t.Fatalf("Failed to create text file: %v", err)
		}

		koreanWords := []string{
			"데이터베이스",
			"소프트웨어",
			"알고리즘",
			"인공지능",
			"머신러닝",
		}

		for _, word := range koreanWords {
			_, err = file.WriteString(word + "\n")
			if err != nil {
				t.Fatalf("Failed to write to text file: %v", err)
			}
		}
		file.Close()

		// Test language detection
		detectedLanguage := detectWordListLanguage(textFile)
		if detectedLanguage != "korean" {
			t.Errorf("Expected language 'korean', got '%s'", detectedLanguage)
		}
	})

	// Test English file detection (should not be detected as Korean)
	t.Run("English file detection", func(t *testing.T) {
		// Create an English JSON word list file
		englishWords := words.WordSource{
			Metadata: words.Metadata{
				Name:       "Test English Words",
				Size:       5,
				PackagedAt: time.Now().Format(time.RFC3339),
				Version:    1,
			},
			Words: []string{"hello", "world", "programming", "computer", "english"},
		}

		jsonFile := filepath.Join(tempDir, "english_words.json")
		file, err := os.Create(jsonFile)
		if err != nil {
			t.Fatalf("Failed to create JSON file: %v", err)
		}

		encoder := json.NewEncoder(file)
		err = encoder.Encode(englishWords)
		file.Close()
		if err != nil {
			t.Fatalf("Failed to write JSON file: %v", err)
		}

		// Test language detection
		detectedLanguage := detectWordListLanguage(jsonFile)
		if detectedLanguage != "english" {
			t.Errorf("Expected language 'english', got '%s'", detectedLanguage)
		}
	})

	// Test mixed content file detection
	t.Run("Mixed content file detection", func(t *testing.T) {
		// Create a file with mixed Korean and English content
		mixedWords := words.WordSource{
			Metadata: words.Metadata{
				Name:       "Test Mixed Words",
				Size:       6,
				PackagedAt: time.Now().Format(time.RFC3339),
				Version:    1,
			},
			Words: []string{"안녕하세요", "hello", "프로그래밍", "programming", "한국어", "english"},
		}

		jsonFile := filepath.Join(tempDir, "mixed_words.json")
		file, err := os.Create(jsonFile)
		if err != nil {
			t.Fatalf("Failed to create JSON file: %v", err)
		}

		encoder := json.NewEncoder(file)
		err = encoder.Encode(mixedWords)
		file.Close()
		if err != nil {
			t.Fatalf("Failed to write JSON file: %v", err)
		}

		// Test language detection - should detect as Korean since Korean ratio > 30%
		detectedLanguage := detectWordListLanguage(jsonFile)
		if detectedLanguage != "korean" {
			t.Errorf("Expected language 'korean' for mixed content, got '%s'", detectedLanguage)
		}
	})
}

func TestKoreanContentDetection(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected bool
	}{
		{"Pure Korean text", "안녕하세요 반갑습니다", true},
		{"Pure English text", "hello world programming", false},
		{"Mixed with Korean majority", "안녕하세요 hello 프로그래밍", true},
		{"Mixed with English majority", "hello 안녕 world programming", false},
		{"Empty text", "", false},
		{"Numbers only", "123 456 789", false},
		{"Korean with punctuation", "안녕하세요! 반갑습니다.", true},
		{"Single Korean character", "안", true},
		{"Korean Jamo characters", "ㄱㄴㄷㄹㅁㅂㅅㅇㅈㅊㅋㅌㅍㅎ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectKoreanContent(tt.text)
			if result != tt.expected {
				t.Errorf("detectKoreanContent(%q) = %v, expected %v", tt.text, result, tt.expected)
			}
		})
	}
}

func TestKoreanCustomFileIntegration(t *testing.T) {
	// Create temporary directory for test files
	tempDir, err := os.MkdirTemp("", "korean_integration_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a Korean custom file
	koreanWords := words.WordSource{
		Metadata: words.Metadata{
			Name:       "Custom Korean Words",
			Size:       3,
			PackagedAt: time.Now().Format(time.RFC3339),
			Version:    1,
		},
		Words: []string{"사용자정의", "한글단어", "테스트"},
	}

	jsonFile := filepath.Join(tempDir, "custom_korean.json")
	file, err := os.Create(jsonFile)
	if err != nil {
		t.Fatalf("Failed to create JSON file: %v", err)
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(koreanWords)
	file.Close()
	if err != nil {
		t.Fatalf("Failed to write JSON file: %v", err)
	}

	// Test that the Korean generator can load the custom file
	generator := words.NewKoreanGenerator([]string{jsonFile})
	
	// Generate content using the custom Korean file (use file path as key)
	content := generator.GenerateKoreanWords(jsonFile)
	
	if len(content) == 0 {
		t.Error("Expected generated content to be non-empty")
	}

	// Verify that the generated content contains Korean characters
	contentStr := string(content)
	if !detectKoreanContent(contentStr) {
		t.Errorf("Expected generated content to contain Korean characters, got: %s", contentStr)
	}
}

func TestKoreanFileValidation(t *testing.T) {
	tests := []struct {
		name     string
		words    []string
		hasError bool
	}{
		{"Valid Korean words", []string{"안녕하세요", "반갑습니다", "프로그래밍"}, false},
		{"Valid English words", []string{"hello", "world", "programming"}, false},
		{"Mixed content", []string{"안녕하세요", "hello", "프로그래밍"}, false},
		{"Empty list", []string{}, false},
		{"Single Korean word", []string{"안녕하세요"}, false},
		{"Single English word", []string{"hello"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateKoreanWordList(tt.words)
			hasError := err != nil
			if hasError != tt.hasError {
				t.Errorf("validateKoreanWordList(%v) error = %v, expected error = %v", tt.words, hasError, tt.hasError)
			}
		})
	}
}