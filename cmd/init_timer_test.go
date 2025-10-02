package cmd

import (
	"testing"
	"time"

	"github.com/bloznelis/typioca/cmd/words"
)

func TestInitTimerBasedTestWithKoreanContent(t *testing.T) {
	// Create test settings with Korean wordlist
	koreanSelections := []WordsSelection{
		{name: "Korean common words", generatorKey: "Korean common words"},
		{name: "Korean tech terms", generatorKey: "Korean tech terms"},
		{name: "Korean sentences", generatorKey: "Korean sentences"},
	}
	
	settings := TimerBasedTestSettings{
		timeSelections:     []time.Duration{30 * time.Second},
		timeCursor:         0,
		wordListSelections: koreanSelections,
		wordListCursor:     0, // Select Korean common words
		cursor:             0,
		enabled:            true,
	}
	
	// Create main menu with Korean generators
	mainMenu := MainMenu{
		timeBasedGenerator:       words.NewGenerator([]string{}),
		koreanTimeBasedGenerator: words.NewKoreanGenerator([]string{}),
	}
	
	// Test Korean common words
	test := initTimerBasedTest(settings, mainMenu)
	
	// Verify that content was generated
	if len(test.base.wordsToEnter) == 0 {
		t.Error("Expected Korean content to be generated, but got empty content")
	}
	
	// Test Korean tech terms
	settings.wordListCursor = 1
	test = initTimerBasedTest(settings, mainMenu)
	
	if len(test.base.wordsToEnter) == 0 {
		t.Error("Expected Korean tech content to be generated, but got empty content")
	}
	
	// Test Korean sentences
	settings.wordListCursor = 2
	test = initTimerBasedTest(settings, mainMenu)
	
	if len(test.base.wordsToEnter) == 0 {
		t.Error("Expected Korean sentence content to be generated, but got empty content")
	}
}

func TestInitTimerBasedTestWithEnglishContent(t *testing.T) {
	// Create test settings with English wordlist
	englishSelections := []WordsSelection{
		{name: "Common words", generatorKey: "Common words"},
		{name: "Frankenstein sentences", generatorKey: "Frankenstein sentences"},
	}
	
	settings := TimerBasedTestSettings{
		timeSelections:     []time.Duration{30 * time.Second},
		timeCursor:         0,
		wordListSelections: englishSelections,
		wordListCursor:     0, // Select Common words
		cursor:             0,
		enabled:            true,
	}
	
	// Create main menu with generators
	mainMenu := MainMenu{
		timeBasedGenerator:       words.NewGenerator([]string{}),
		koreanTimeBasedGenerator: words.NewKoreanGenerator([]string{}),
	}
	
	// Test English common words
	test := initTimerBasedTest(settings, mainMenu)
	
	// Verify that content was generated
	if len(test.base.wordsToEnter) == 0 {
		t.Error("Expected English content to be generated, but got empty content")
	}
	
	// Test English sentences
	settings.wordListCursor = 1
	test = initTimerBasedTest(settings, mainMenu)
	
	if len(test.base.wordsToEnter) == 0 {
		t.Error("Expected English sentence content to be generated, but got empty content")
	}
}

func TestInitTimerBasedTestKoreanRouting(t *testing.T) {
	// Test that Korean wordlists are properly routed to Korean generators
	tests := []struct {
		name         string
		generatorKey string
		expectKorean bool
		expectType   KoreanGenerationType
	}{
		{"Korean common words", "Korean common words", true, KoreanWords},
		{"Korean tech terms", "Korean tech terms", true, KoreanWords},
		{"Korean sentences", "Korean sentences", true, KoreanSentences},
		{"korean-common", "korean-common", true, KoreanWords},
		{"korean-tech", "korean-tech", true, KoreanWords},
		{"korean-sentences", "korean-sentences", true, KoreanSentences},
		{"English common", "Common words", false, ""},
		{"Frankenstein", "Frankenstein sentences", false, ""},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapping, isKorean := getKoreanGenerationType(tt.generatorKey)
			
			if isKorean != tt.expectKorean {
				t.Errorf("getKoreanGenerationType(%q) Korean detection = %v, expected %v", 
					tt.generatorKey, isKorean, tt.expectKorean)
			}
			
			if tt.expectKorean && mapping.GenerationType != tt.expectType {
				t.Errorf("getKoreanGenerationType(%q) type = %v, expected %v", 
					tt.generatorKey, mapping.GenerationType, tt.expectType)
			}
		})
	}
}