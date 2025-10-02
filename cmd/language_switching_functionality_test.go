package cmd

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/bloznelis/typioca/cmd/words"
)

// TestLanguageSwitchingFromEnglishToKorean tests switching from English to Korean wordlists
func TestLanguageSwitchingFromEnglishToKorean(t *testing.T) {
	// Create mixed wordlist selections (English and Korean)
	mixedSelections := []WordsSelection{
		{name: "Common words", generatorKey: "Common words"},
		{name: "Frankenstein sentences", generatorKey: "Frankenstein sentences"},
		{name: "Korean common words", generatorKey: "Korean common words"},
		{name: "Korean tech terms", generatorKey: "Korean tech terms"},
		{name: "Korean sentences", generatorKey: "Korean sentences"},
	}
	
	// Create main menu with both generators
	mainMenu := MainMenu{
		timeBasedGenerator:       words.NewGenerator([]string{}),
		koreanTimeBasedGenerator: words.NewKoreanGenerator([]string{}),
	}
	
	// Start with English wordlist (Common words)
	settings := TimerBasedTestSettings{
		timeSelections:     []time.Duration{30 * time.Second},
		timeCursor:         0,
		wordListSelections: mixedSelections,
		wordListCursor:     0, // English: Common words
		cursor:             0,
		enabled:            true,
	}
	
	// Test English content generation
	englishTest := initTimerBasedTest(settings, mainMenu)
	if len(englishTest.base.wordsToEnter) == 0 {
		t.Error("Expected English content to be generated, but got empty content")
	}
	
	// Verify English content doesn't contain Korean characters
	englishContent := string(englishTest.base.wordsToEnter)
	if containsKoreanCharacters(englishContent) {
		t.Error("English wordlist generated Korean characters, which is unexpected")
	}
	
	// Switch to Korean wordlist (Korean common words)
	settings.wordListCursor = 2 // Korean common words
	koreanTest := initTimerBasedTest(settings, mainMenu)
	
	if len(koreanTest.base.wordsToEnter) == 0 {
		t.Error("Expected Korean content to be generated after switching from English, but got empty content")
	}
	
	// Verify Korean content contains Korean characters
	koreanContent := string(koreanTest.base.wordsToEnter)
	if !containsKoreanCharacters(koreanContent) {
		t.Error("Korean wordlist did not generate Korean characters after switching from English")
	}
	
	// Verify content is different between English and Korean
	if englishContent == koreanContent {
		t.Error("English and Korean content should be different")
	}
}

// TestLanguageSwitchingFromKoreanToEnglish tests switching from Korean to English wordlists
func TestLanguageSwitchingFromKoreanToEnglish(t *testing.T) {
	// Create mixed wordlist selections (Korean and English)
	mixedSelections := []WordsSelection{
		{name: "Korean common words", generatorKey: "Korean common words"},
		{name: "Korean tech terms", generatorKey: "Korean tech terms"},
		{name: "Common words", generatorKey: "Common words"},
		{name: "Frankenstein sentences", generatorKey: "Frankenstein sentences"},
	}
	
	// Create main menu with both generators
	mainMenu := MainMenu{
		timeBasedGenerator:       words.NewGenerator([]string{}),
		koreanTimeBasedGenerator: words.NewKoreanGenerator([]string{}),
	}
	
	// Start with Korean wordlist (Korean common words)
	settings := TimerBasedTestSettings{
		timeSelections:     []time.Duration{30 * time.Second},
		timeCursor:         0,
		wordListSelections: mixedSelections,
		wordListCursor:     0, // Korean: Korean common words
		cursor:             0,
		enabled:            true,
	}
	
	// Test Korean content generation
	koreanTest := initTimerBasedTest(settings, mainMenu)
	if len(koreanTest.base.wordsToEnter) == 0 {
		t.Error("Expected Korean content to be generated, but got empty content")
	}
	
	// Verify Korean content contains Korean characters
	koreanContent := string(koreanTest.base.wordsToEnter)
	if !containsKoreanCharacters(koreanContent) {
		t.Error("Korean wordlist did not generate Korean characters")
	}
	
	// Switch to English wordlist (Common words)
	settings.wordListCursor = 2 // Common words
	englishTest := initTimerBasedTest(settings, mainMenu)
	
	if len(englishTest.base.wordsToEnter) == 0 {
		t.Error("Expected English content to be generated after switching from Korean, but got empty content")
	}
	
	// Verify English content doesn't contain Korean characters
	englishContent := string(englishTest.base.wordsToEnter)
	if containsKoreanCharacters(englishContent) {
		t.Error("English wordlist generated Korean characters after switching from Korean, which is unexpected")
	}
	
	// Verify content is different between Korean and English
	if koreanContent == englishContent {
		t.Error("Korean and English content should be different")
	}
}

// TestRapidSwitchingBetweenKoreanWordlistTypes tests rapid switching between different Korean wordlist types
func TestRapidSwitchingBetweenKoreanWordlistTypes(t *testing.T) {
	// Create Korean-only wordlist selections
	koreanSelections := []WordsSelection{
		{name: "Korean common words", generatorKey: "Korean common words"},
		{name: "Korean tech terms", generatorKey: "Korean tech terms"},
		{name: "Korean sentences", generatorKey: "Korean sentences"},
	}
	
	// Create main menu with Korean generator
	mainMenu := MainMenu{
		timeBasedGenerator:       words.NewGenerator([]string{}),
		koreanTimeBasedGenerator: words.NewKoreanGenerator([]string{}),
	}
	
	settings := TimerBasedTestSettings{
		timeSelections:     []time.Duration{30 * time.Second},
		timeCursor:         0,
		wordListSelections: koreanSelections,
		wordListCursor:     0, // Start with Korean common words
		cursor:             0,
		enabled:            true,
	}
	
	// Test each Korean wordlist type and rapid switching
	var previousContent string
	
	for i := 0; i < len(koreanSelections); i++ {
		settings.wordListCursor = i
		test := initTimerBasedTest(settings, mainMenu)
		
		if len(test.base.wordsToEnter) == 0 {
			t.Errorf("Expected Korean content for wordlist %d (%s), but got empty content", 
				i, koreanSelections[i].name)
			continue
		}
		
		currentContent := string(test.base.wordsToEnter)
		
		// Verify Korean content contains Korean characters
		if !containsKoreanCharacters(currentContent) {
			t.Errorf("Korean wordlist %d (%s) did not generate Korean characters", 
				i, koreanSelections[i].name)
		}
		
		// Verify content is different from previous (except for first iteration)
		if i > 0 && currentContent == previousContent {
			t.Errorf("Korean wordlist %d (%s) generated same content as previous wordlist", 
				i, koreanSelections[i].name)
		}
		
		previousContent = currentContent
		
		// Log for debugging
		t.Logf("Korean wordlist %d (%s): Generated %d characters with Korean content", 
			i, koreanSelections[i].name, len(test.base.wordsToEnter))
	}
	
	// Test rapid back-and-forth switching
	for switchCount := 0; switchCount < 5; switchCount++ {
		// Switch to Korean common words
		settings.wordListCursor = 0
		commonTest := initTimerBasedTest(settings, mainMenu)
		
		// Switch to Korean sentences
		settings.wordListCursor = 2
		sentenceTest := initTimerBasedTest(settings, mainMenu)
		
		// Verify both generate content
		if len(commonTest.base.wordsToEnter) == 0 {
			t.Errorf("Rapid switch %d: Korean common words failed to generate content", switchCount)
		}
		if len(sentenceTest.base.wordsToEnter) == 0 {
			t.Errorf("Rapid switch %d: Korean sentences failed to generate content", switchCount)
		}
		
		// Verify both contain Korean characters
		commonContent := string(commonTest.base.wordsToEnter)
		sentenceContent := string(sentenceTest.base.wordsToEnter)
		
		if !containsKoreanCharacters(commonContent) {
			t.Errorf("Rapid switch %d: Korean common words lost Korean characters", switchCount)
		}
		if !containsKoreanCharacters(sentenceContent) {
			t.Errorf("Rapid switch %d: Korean sentences lost Korean characters", switchCount)
		}
	}
}

// TestLanguageSwitchingWithWordCountBasedTest tests language switching in word count based tests
func TestLanguageSwitchingWithWordCountBasedTest(t *testing.T) {
	// Create mixed wordlist selections
	mixedSelections := []WordsSelection{
		{name: "Common words", generatorKey: "Common words"},
		{name: "Korean common words", generatorKey: "Korean common words"},
		{name: "Korean tech terms", generatorKey: "Korean tech terms"},
	}
	
	// Create main menu with both generators
	mainMenu := MainMenu{
		wordCountGenerator:       words.NewGenerator([]string{}),
		koreanWordCountGenerator: words.NewKoreanGenerator([]string{}),
	}
	
	// Set word count for generators
	mainMenu.wordCountGenerator.Count = 25
	mainMenu.koreanWordCountGenerator.Count = 25
	
	settings := WordCountBasedTestSettings{
		wordCountSelections: []int{25, 50, 100},
		wordCountCursor:     0,
		wordListSelections:  mixedSelections,
		wordListCursor:      0, // Start with English
		cursor:              0,
		enabled:             true,
	}
	
	// Test English content generation
	englishTest := initWordCountBasedTest(settings, mainMenu)
	if len(englishTest.base.wordsToEnter) == 0 {
		t.Error("Expected English content in word count test, but got empty content")
	}
	
	englishContent := string(englishTest.base.wordsToEnter)
	if containsKoreanCharacters(englishContent) {
		t.Error("English word count test generated Korean characters")
	}
	
	// Switch to Korean wordlist
	settings.wordListCursor = 1 // Korean common words
	koreanTest := initWordCountBasedTest(settings, mainMenu)
	
	if len(koreanTest.base.wordsToEnter) == 0 {
		t.Error("Expected Korean content in word count test after switching, but got empty content")
	}
	
	koreanContent := string(koreanTest.base.wordsToEnter)
	if !containsKoreanCharacters(koreanContent) {
		t.Error("Korean word count test did not generate Korean characters after switching")
	}
	
	// Verify content is different
	if englishContent == koreanContent {
		t.Error("English and Korean word count content should be different")
	}
}

// TestLanguageSwitchingWithSentenceCountBasedTest tests language switching in sentence count based tests
func TestLanguageSwitchingWithSentenceCountBasedTest(t *testing.T) {
	// Create mixed sentence selections
	mixedSelections := []WordsSelection{
		{name: "Frankenstein sentences", generatorKey: "Frankenstein sentences"},
		{name: "Korean sentences", generatorKey: "Korean sentences"},
	}
	
	// Create main menu with both generators
	mainMenu := MainMenu{
		sentenceCountGenerator: words.NewGenerator([]string{}),
		koreanSentenceGenerator: words.NewKoreanGenerator([]string{}),
	}
	
	// Set sentence count for generators
	mainMenu.sentenceCountGenerator.Count = 5
	mainMenu.koreanSentenceGenerator.Count = 5
	
	settings := SentenceCountBasedTestSettings{
		sentenceCountSelections: []int{5, 10, 15},
		sentenceCountCursor:     0,
		sentenceListSelections:  mixedSelections,
		sentenceListCursor:      0, // Start with English
		cursor:                  0,
		enabled:                 true,
	}
	
	// Test English sentence generation
	englishTest := initSentenceCountBasedTest(settings, mainMenu)
	if len(englishTest.base.wordsToEnter) == 0 {
		t.Error("Expected English sentences, but got empty content")
	}
	
	englishContent := string(englishTest.base.wordsToEnter)
	if containsKoreanCharacters(englishContent) {
		t.Error("English sentence test generated Korean characters")
	}
	
	// Switch to Korean sentences
	settings.sentenceListCursor = 1 // Korean sentences
	koreanTest := initSentenceCountBasedTest(settings, mainMenu)
	
	if len(koreanTest.base.wordsToEnter) == 0 {
		t.Error("Expected Korean sentences after switching, but got empty content")
	}
	
	koreanContent := string(koreanTest.base.wordsToEnter)
	if !containsKoreanCharacters(koreanContent) {
		t.Error("Korean sentence test did not generate Korean characters after switching")
	}
	
	// Verify content is different
	if englishContent == koreanContent {
		t.Error("English and Korean sentence content should be different")
	}
}

// TestLanguageSwitchingErrorHandling tests error handling during language switching
func TestLanguageSwitchingErrorHandling(t *testing.T) {
	// Create selections with potentially problematic entries
	mixedSelections := []WordsSelection{
		{name: "Common words", generatorKey: "Common words"},
		{name: "Invalid Korean", generatorKey: "korean-invalid"},
		{name: "Korean common words", generatorKey: "Korean common words"},
	}
	
	// Create main menu with generators
	mainMenu := MainMenu{
		timeBasedGenerator:       words.NewGenerator([]string{}),
		koreanTimeBasedGenerator: words.NewKoreanGenerator([]string{}),
	}
	
	settings := TimerBasedTestSettings{
		timeSelections:     []time.Duration{30 * time.Second},
		timeCursor:         0,
		wordListSelections: mixedSelections,
		wordListCursor:     0,
		cursor:             0,
		enabled:            true,
	}
	
	// Test each wordlist and ensure no crashes occur
	for i := 0; i < len(mixedSelections); i++ {
		settings.wordListCursor = i
		
		// This should not panic or crash
		test := initTimerBasedTest(settings, mainMenu)
		
		// Should always generate some content (even if fallback)
		if len(test.base.wordsToEnter) == 0 {
			t.Errorf("Wordlist %d (%s) failed to generate any content, including fallback", 
				i, mixedSelections[i].name)
		}
		
		t.Logf("Wordlist %d (%s): Generated %d characters", 
			i, mixedSelections[i].name, len(test.base.wordsToEnter))
	}
}

// TestCursorPersistenceAcrossLanguageSwitches tests that cursor positions are properly maintained
func TestCursorPersistenceAcrossLanguageSwitches(t *testing.T) {
	// Create mixed wordlist selections
	mixedSelections := []WordsSelection{
		{name: "Common words", generatorKey: "Common words"},
		{name: "Korean common words", generatorKey: "Korean common words"},
		{name: "Frankenstein sentences", generatorKey: "Frankenstein sentences"},
		{name: "Korean sentences", generatorKey: "Korean sentences"},
	}
	
	// Create main menu
	mainMenu := MainMenu{
		timeBasedGenerator:       words.NewGenerator([]string{}),
		koreanTimeBasedGenerator: words.NewKoreanGenerator([]string{}),
	}
	
	settings := TimerBasedTestSettings{
		timeSelections:     []time.Duration{15 * time.Second, 30 * time.Second, 60 * time.Second},
		timeCursor:         1, // 30 seconds
		wordListSelections: mixedSelections,
		wordListCursor:     0, // Common words
		cursor:             0,
		enabled:            true,
	}
	
	// Test that time cursor persists across wordlist changes
	originalTimeCursor := settings.timeCursor
	
	// Switch wordlists multiple times
	for i := 0; i < len(mixedSelections); i++ {
		settings.wordListCursor = i
		test := initTimerBasedTest(settings, mainMenu)
		
		// Verify time cursor hasn't changed
		if settings.timeCursor != originalTimeCursor {
			t.Errorf("Time cursor changed from %d to %d when switching to wordlist %d", 
				originalTimeCursor, settings.timeCursor, i)
		}
		
		// Verify content was generated
		if len(test.base.wordsToEnter) == 0 {
			t.Errorf("No content generated for wordlist %d (%s)", i, mixedSelections[i].name)
		}
	}
}

// TestLanguageSwitchingPerformance tests that language switching doesn't cause performance degradation
func TestLanguageSwitchingPerformance(t *testing.T) {
	// Create mixed wordlist selections
	mixedSelections := []WordsSelection{
		{name: "Common words", generatorKey: "Common words"},
		{name: "Korean common words", generatorKey: "Korean common words"},
		{name: "Korean tech terms", generatorKey: "Korean tech terms"},
		{name: "Frankenstein sentences", generatorKey: "Frankenstein sentences"},
		{name: "Korean sentences", generatorKey: "Korean sentences"},
	}
	
	// Create main menu with both generators
	mainMenu := MainMenu{
		timeBasedGenerator:       words.NewGenerator([]string{}),
		koreanTimeBasedGenerator: words.NewKoreanGenerator([]string{}),
	}
	
	settings := TimerBasedTestSettings{
		timeSelections:     []time.Duration{30 * time.Second},
		timeCursor:         0,
		wordListSelections: mixedSelections,
		wordListCursor:     0,
		cursor:             0,
		enabled:            true,
	}
	
	// Measure performance of rapid switching
	start := time.Now()
	switchCount := 50 // Test 50 rapid switches
	
	for i := 0; i < switchCount; i++ {
		// Alternate between English and Korean
		if i%2 == 0 {
			settings.wordListCursor = 0 // English
		} else {
			settings.wordListCursor = 1 // Korean
		}
		
		test := initTimerBasedTest(settings, mainMenu)
		
		// Verify content was generated
		if len(test.base.wordsToEnter) == 0 {
			t.Errorf("Performance test iteration %d: No content generated", i)
		}
	}
	
	elapsed := time.Since(start)
	avgTimePerSwitch := elapsed / time.Duration(switchCount)
	
	// Performance should be reasonable (less than 100ms per switch)
	if avgTimePerSwitch > 100*time.Millisecond {
		t.Errorf("Language switching performance too slow: %v per switch (expected < 100ms)", avgTimePerSwitch)
	}
	
	t.Logf("Language switching performance: %v per switch (%d switches in %v)", 
		avgTimePerSwitch, switchCount, elapsed)
}

// TestLanguageSwitchingMemoryUsage tests that language switching doesn't cause memory leaks
func TestLanguageSwitchingMemoryUsage(t *testing.T) {
	// Create mixed wordlist selections
	mixedSelections := []WordsSelection{
		{name: "Common words", generatorKey: "Common words"},
		{name: "Korean common words", generatorKey: "Korean common words"},
		{name: "Korean sentences", generatorKey: "Korean sentences"},
	}
	
	// Create main menu with both generators
	mainMenu := MainMenu{
		timeBasedGenerator:       words.NewGenerator([]string{}),
		koreanTimeBasedGenerator: words.NewKoreanGenerator([]string{}),
	}
	
	settings := TimerBasedTestSettings{
		timeSelections:     []time.Duration{30 * time.Second},
		timeCursor:         0,
		wordListSelections: mixedSelections,
		wordListCursor:     0,
		cursor:             0,
		enabled:            true,
	}
	
	// Force garbage collection before test
	runtime.GC()
	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)
	
	// Perform many language switches
	for i := 0; i < 100; i++ {
		settings.wordListCursor = i % len(mixedSelections)
		test := initTimerBasedTest(settings, mainMenu)
		
		// Verify content was generated
		if len(test.base.wordsToEnter) == 0 {
			t.Errorf("Memory test iteration %d: No content generated", i)
		}
		
		// Clear the test to help with memory management
		test.base.wordsToEnter = nil
	}
	
	// Force garbage collection after test
	runtime.GC()
	runtime.ReadMemStats(&m2)
	
	// Check memory usage increase (handle potential underflow)
	var memIncrease uint64
	if m2.Alloc > m1.Alloc {
		memIncrease = m2.Alloc - m1.Alloc
	} else {
		// Memory decreased or stayed the same (due to GC)
		memIncrease = 0
	}
	
	// Memory increase should be reasonable (less than 10MB for 100 switches)
	if memIncrease > 10*1024*1024 {
		t.Errorf("Excessive memory usage during language switching: %d bytes increase", memIncrease)
	}
	
	t.Logf("Memory usage during language switching: %d bytes increase", memIncrease)
}

// TestLanguageSwitchingWithInvalidSelections tests switching with invalid or corrupted selections
func TestLanguageSwitchingWithInvalidSelections(t *testing.T) {
	// Create selections with some invalid entries
	invalidSelections := []WordsSelection{
		{name: "Common words", generatorKey: "Common words"},
		{name: "Invalid Korean", generatorKey: "korean-nonexistent"},
		{name: "Empty Korean", generatorKey: "korean-empty"},
		{name: "Korean common words", generatorKey: "Korean common words"},
		{name: "Invalid English", generatorKey: "nonexistent-english"},
	}
	
	// Create main menu with both generators
	mainMenu := MainMenu{
		timeBasedGenerator:       words.NewGenerator([]string{}),
		koreanTimeBasedGenerator: words.NewKoreanGenerator([]string{}),
	}
	
	settings := TimerBasedTestSettings{
		timeSelections:     []time.Duration{30 * time.Second},
		timeCursor:         0,
		wordListSelections: invalidSelections,
		wordListCursor:     0,
		cursor:             0,
		enabled:            true,
	}
	
	// Test each selection, including invalid ones
	for i := 0; i < len(invalidSelections); i++ {
		settings.wordListCursor = i
		
		// This should not panic or crash, even with invalid selections
		test := initTimerBasedTest(settings, mainMenu)
		
		// Should always generate some content (even if fallback)
		if len(test.base.wordsToEnter) == 0 {
			t.Errorf("Invalid selection test %d (%s): Failed to generate any content, including fallback", 
				i, invalidSelections[i].name)
		}
		
		// Content should be reasonable length (at least 10 characters)
		if len(test.base.wordsToEnter) < 10 {
			t.Errorf("Invalid selection test %d (%s): Generated content too short: %d characters", 
				i, invalidSelections[i].name, len(test.base.wordsToEnter))
		}
		
		t.Logf("Invalid selection test %d (%s): Generated %d characters", 
			i, invalidSelections[i].name, len(test.base.wordsToEnter))
	}
}

// TestLanguageSwitchingStateConsistency tests that generator state remains consistent across switches
func TestLanguageSwitchingStateConsistency(t *testing.T) {
	// Create mixed wordlist selections
	mixedSelections := []WordsSelection{
		{name: "Common words", generatorKey: "Common words"},
		{name: "Korean common words", generatorKey: "Korean common words"},
		{name: "Korean tech terms", generatorKey: "Korean tech terms"},
	}
	
	// Create main menu with both generators
	mainMenu := MainMenu{
		timeBasedGenerator:       words.NewGenerator([]string{}),
		koreanTimeBasedGenerator: words.NewKoreanGenerator([]string{}),
	}
	
	// Set specific counts for generators
	mainMenu.timeBasedGenerator.Count = 100
	mainMenu.koreanTimeBasedGenerator.Count = 100
	
	settings := TimerBasedTestSettings{
		timeSelections:     []time.Duration{30 * time.Second},
		timeCursor:         0,
		wordListSelections: mixedSelections,
		wordListCursor:     0,
		cursor:             0,
		enabled:            true,
	}
	
	// Test that generator counts remain consistent
	originalEnglishCount := mainMenu.timeBasedGenerator.Count
	originalKoreanCount := mainMenu.koreanTimeBasedGenerator.Count
	
	// Perform multiple switches
	for i := 0; i < 10; i++ {
		settings.wordListCursor = i % len(mixedSelections)
		test := initTimerBasedTest(settings, mainMenu)
		
		// Verify content was generated
		if len(test.base.wordsToEnter) == 0 {
			t.Errorf("State consistency test iteration %d: No content generated", i)
		}
		
		// Verify generator counts haven't been corrupted
		if mainMenu.timeBasedGenerator.Count != originalEnglishCount {
			t.Errorf("English generator count changed from %d to %d after switch %d", 
				originalEnglishCount, mainMenu.timeBasedGenerator.Count, i)
		}
		
		if mainMenu.koreanTimeBasedGenerator.Count != originalKoreanCount {
			t.Errorf("Korean generator count changed from %d to %d after switch %d", 
				originalKoreanCount, mainMenu.koreanTimeBasedGenerator.Count, i)
		}
	}
}

// TestLanguageSwitchingConcurrency tests language switching under concurrent access
func TestLanguageSwitchingConcurrency(t *testing.T) {
	// Create mixed wordlist selections
	mixedSelections := []WordsSelection{
		{name: "Common words", generatorKey: "Common words"},
		{name: "Korean common words", generatorKey: "Korean common words"},
		{name: "Korean sentences", generatorKey: "Korean sentences"},
	}
	
	// Create main menu with both generators
	mainMenu := MainMenu{
		timeBasedGenerator:       words.NewGenerator([]string{}),
		koreanTimeBasedGenerator: words.NewKoreanGenerator([]string{}),
	}
	
	// Run concurrent language switches
	var wg sync.WaitGroup
	errors := make(chan error, 10)
	
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			
			settings := TimerBasedTestSettings{
				timeSelections:     []time.Duration{30 * time.Second},
				timeCursor:         0,
				wordListSelections: mixedSelections,
				wordListCursor:     goroutineID % len(mixedSelections),
				cursor:             0,
				enabled:            true,
			}
			
			// Perform multiple switches in this goroutine
			for j := 0; j < 5; j++ {
				settings.wordListCursor = (goroutineID + j) % len(mixedSelections)
				test := initTimerBasedTest(settings, mainMenu)
				
				if len(test.base.wordsToEnter) == 0 {
					errors <- fmt.Errorf("goroutine %d, iteration %d: No content generated", goroutineID, j)
					return
				}
			}
		}(i)
	}
	
	wg.Wait()
	close(errors)
	
	// Check for any errors
	for err := range errors {
		t.Error(err)
	}
}

// Helper function to detect Korean characters in content
func containsKoreanCharacters(content string) bool {
	for _, r := range content {
		if isKoreanChar(r) {
			return true
		}
	}
	return false
}

// Note: isKoreanChar function is already defined in korean.go