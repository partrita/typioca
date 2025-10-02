package cmd

import (
	"testing"
	"time"
)

func TestIsKoreanChar(t *testing.T) {
	tests := []struct {
		name     string
		input    rune
		expected bool
	}{
		{"Korean syllable 안", '안', true},
		{"Korean syllable 녕", '녕', true},
		{"Korean syllable 하", '하', true},
		{"Korean syllable 세", '세', true},
		{"Korean syllable 요", '요', true},
		{"English letter a", 'a', false},
		{"English letter Z", 'Z', false},
		{"Number 1", '1', false},
		{"Space", ' ', false},
		{"Punctuation .", '.', false},
		{"Korean Jamo ㄱ", 'ㄱ', true},
		{"Korean Jamo ㅏ", 'ㅏ', true},
		{"Korean Jamo ㅇ", 'ㅇ', true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isKoreanChar(tt.input)
			if result != tt.expected {
				t.Errorf("isKoreanChar(%c) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCalculateKoreanWPM(t *testing.T) {
	tests := []struct {
		name     string
		chars    int
		duration time.Duration
		expected int
	}{
		{"10 chars in 1 minute", 10, time.Minute, 5}, // (10 * 2.5) / 5 / 1 = 5
		{"20 chars in 2 minutes", 20, 2 * time.Minute, 5}, // (20 * 2.5) / 5 / 2 = 5
		{"50 chars in 1 minute", 50, time.Minute, 25}, // (50 * 2.5) / 5 / 1 = 25
		{"0 chars", 0, time.Minute, 0},
		{"Zero duration", 10, 0, 0},
		{"30 seconds", 10, 30 * time.Second, 10}, // (10 * 2.5) / 5 / 0.5 = 10
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateKoreanWPM(tt.chars, tt.duration)
			if result != tt.expected {
				t.Errorf("calculateKoreanWPM(%d, %v) = %d, expected %d", tt.chars, tt.duration, result, tt.expected)
			}
		})
	}
}

func TestCalculateKoreanAccuracy(t *testing.T) {
	tests := []struct {
		name     string
		input    []rune
		target   []rune
		expected float64
	}{
		{
			"Perfect match",
			[]rune("안녕하세요"),
			[]rune("안녕하세요"),
			100.0,
		},
		{
			"No match",
			[]rune("가나다라마"),
			[]rune("안녕하세요"),
			0.0,
		},
		{
			"Partial match",
			[]rune("안녕다라마"),
			[]rune("안녕하세요"),
			40.0, // 2 out of 5 correct
		},
		{
			"Input shorter than target",
			[]rune("안녕"),
			[]rune("안녕하세요"),
			40.0, // 2 out of 5 correct
		},
		{
			"Input longer than target",
			[]rune("안녕하세요반갑"),
			[]rune("안녕하세요"),
			100.0, // All target characters match
		},
		{
			"Empty target",
			[]rune("안녕"),
			[]rune(""),
			100.0,
		},
		{
			"Empty input",
			[]rune(""),
			[]rune("안녕하세요"),
			0.0,
		},
		{
			"Single character correct",
			[]rune("안"),
			[]rune("안녕하세요"),
			20.0, // 1 out of 5 correct
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateKoreanAccuracy(tt.input, tt.target)
			if result != tt.expected {
				t.Errorf("calculateKoreanAccuracy(%v, %v) = %.1f, expected %.1f", 
					string(tt.input), string(tt.target), result, tt.expected)
			}
		})
	}
}

func TestIsKoreanSyllable(t *testing.T) {
	tests := []struct {
		name     string
		input    rune
		expected bool
	}{
		{"Korean syllable 가", '가', true},
		{"Korean syllable 힣", '힣', true}, // Last Korean syllable
		{"Korean syllable 안", '안', true},
		{"Korean Jamo ㄱ", 'ㄱ', false},
		{"Korean Jamo ㅏ", 'ㅏ', false},
		{"English letter a", 'a', false},
		{"Number 1", '1', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isKoreanSyllable(tt.input)
			if result != tt.expected {
				t.Errorf("isKoreanSyllable(%c) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsKoreanJamo(t *testing.T) {
	tests := []struct {
		name     string
		input    rune
		expected bool
	}{
		{"Korean Jamo ㄱ", 'ㄱ', true},
		{"Korean Jamo ㅏ", 'ㅏ', true},
		{"Korean Jamo ㅇ", 'ㅇ', true},
		{"Korean syllable 안", '안', false},
		{"English letter a", 'a', false},
		{"Number 1", '1', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isKoreanJamo(tt.input)
			if result != tt.expected {
				t.Errorf("isKoreanJamo(%c) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCountKoreanSyllables(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Simple Korean text", "안녕하세요", 5},
		{"Mixed Korean and English", "안녕 hello 세요", 4},
		{"Korean with numbers", "안녕123하세요", 5},
		{"Only English", "hello world", 0},
		{"Empty string", "", 0},
		{"Korean with punctuation", "안녕하세요!", 5},
		{"Korean with spaces", "안녕 하세 요", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := countKoreanSyllables(tt.input)
			if result != tt.expected {
				t.Errorf("countKoreanSyllables(%s) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

// Benchmark tests for performance
func BenchmarkIsKoreanChar(b *testing.B) {
	testRune := '안'
	for i := 0; i < b.N; i++ {
		isKoreanChar(testRune)
	}
}

func BenchmarkCalculateKoreanWPM(b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculateKoreanWPM(100, time.Minute)
	}
}

func BenchmarkCalculateKoreanAccuracy(b *testing.B) {
	input := []rune("안녕하세요반갑습니다")
	target := []rune("안녕하세요반갑습니다")
	
	for i := 0; i < b.N; i++ {
		calculateKoreanAccuracy(input, target)
	}
}

func BenchmarkCountKoreanSyllables(b *testing.B) {
	text := "안녕하세요 반갑습니다 프로그래밍은 재미있습니다"
	
	for i := 0; i < b.N; i++ {
		countKoreanSyllables(text)
	}
}

func TestValidateKoreanWordList(t *testing.T) {
	tests := []struct {
		name     string
		words    []string
		hasError bool
	}{
		{
			"Valid Korean word list",
			[]string{"안녕", "하세요", "프로그래밍", "데이터베이스"},
			false,
		},
		{
			"Mixed Korean and English (majority Korean)",
			[]string{"안녕", "hello", "하세요", "프로그래밍", "데이터베이스"},
			false,
		},
		{
			"Mostly English words",
			[]string{"hello", "world", "programming", "안녕"},
			false, // Not an error, just not Korean
		},
		{
			"Empty word list",
			[]string{},
			false,
		},
		{
			"Single Korean word",
			[]string{"안녕하세요"},
			false,
		},
		{
			"Single English word",
			[]string{"hello"},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateKoreanWordList(tt.words)
			if (err != nil) != tt.hasError {
				t.Errorf("validateKoreanWordList(%v) error = %v, expected error = %v", 
					tt.words, err, tt.hasError)
			}
		})
	}
}

func TestDetectKoreanContent(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected bool
	}{
		{
			"Pure Korean text",
			"안녕하세요",
			true,
		},
		{
			"Mixed Korean and English (majority Korean)",
			"안녕 hello 하세요",
			true,
		},
		{
			"Mixed Korean and English (minority Korean)",
			"hello 안녕 world programming",
			false,
		},
		{
			"Pure English text",
			"hello world",
			false,
		},
		{
			"Korean with numbers and punctuation",
			"안녕하세요! 123번입니다.",
			true,
		},
		{
			"Empty string",
			"",
			false,
		},
		{
			"Only numbers and punctuation",
			"123!@#",
			false,
		},
		{
			"Korean Jamo characters",
			"ㄱㄴㄷㄹㅁㅂㅅㅇ",
			true,
		},
		{
			"Single Korean character",
			"안",
			true,
		},
		{
			"Korean with spaces",
			"안녕 하세 요",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectKoreanContent(tt.text)
			if result != tt.expected {
				t.Errorf("detectKoreanContent(%s) = %v, expected %v", tt.text, result, tt.expected)
			}
		})
	}
}