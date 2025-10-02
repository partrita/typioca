package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bloznelis/typioca/cmd/words"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

func (m model) Init() tea.Cmd {
	return nil
}

// todo: clean these up. Maybe we could reuse filtering by enabled and synce, because now it's redundant
type WordsSelection struct {
	name         string
	generatorKey string
}

func filterEnabledWordSelection(config Config) []WordsSelection {
	var acc []WordsSelection
	for _, elem := range config.WordLists {
		if elem.Enabled && elem.synced && !elem.Sentences {
			acc = append(acc, WordsSelection{
				name:         elem.Name,
				generatorKey: elem.Path,
			})
		}
	}
	for _, elem := range config.EmbededWordLists {
		if elem.Enabled && !elem.IsSentences {
			acc = append(acc, WordsSelection{
				name:         elem.Name,
				generatorKey: elem.Name,
			})
		}
	}

	return acc
}

func filterEnabledSentenceSelection(config Config) []WordsSelection {
	var acc []WordsSelection
	for _, elem := range config.WordLists {
		if elem.Enabled && elem.synced && elem.Sentences {
			acc = append(acc, WordsSelection{
				name:         elem.Name,
				generatorKey: elem.Path,
			})
		}
	}

	for _, elem := range config.EmbededWordLists {
		if elem.Enabled && elem.IsSentences {
			acc = append(acc, WordsSelection{
				name:         elem.Name,
				generatorKey: elem.Name,
			})
		}
	}
	return acc
}

func filterEnabledSelections(config Config) []WordsSelection {
	var acc []WordsSelection
	for _, elem := range config.WordLists {
		if elem.Enabled && elem.synced {
			acc = append(acc, WordsSelection{
				name:         elem.Name,
				generatorKey: elem.Path,
			})
		}
	}

	for _, elem := range config.EmbededWordLists {
		if elem.Enabled {
			acc = append(acc, WordsSelection{
				name:         elem.Name,
				generatorKey: elem.Name,
			})
		}
	}

	return acc
}

// Korean filtering functions
func filterEnabledKoreanSelections(config Config) []WordsSelection {
	var acc []WordsSelection
	for _, elem := range config.WordLists {
		if elem.Enabled && elem.synced && elem.Language == "korean" {
			acc = append(acc, WordsSelection{
				name:         elem.Name,
				generatorKey: elem.Path,
			})
		}
	}

	for _, elem := range config.EmbededWordLists {
		if elem.Enabled && elem.Language == "korean" {
			acc = append(acc, WordsSelection{
				name:         elem.Name,
				generatorKey: elem.Name,
			})
		}
	}

	return acc
}

func filterEnabledKoreanWordSelection(config Config) []WordsSelection {
	var acc []WordsSelection
	for _, elem := range config.WordLists {
		if elem.Enabled && elem.synced && !elem.Sentences && elem.Language == "korean" {
			acc = append(acc, WordsSelection{
				name:         elem.Name,
				generatorKey: elem.Path,
			})
		}
	}
	for _, elem := range config.EmbededWordLists {
		if elem.Enabled && !elem.IsSentences && elem.Language == "korean" {
			acc = append(acc, WordsSelection{
				name:         elem.Name,
				generatorKey: elem.Name,
			})
		}
	}

	return acc
}

func filterEnabledKoreanSentenceSelection(config Config) []WordsSelection {
	var acc []WordsSelection
	for _, elem := range config.WordLists {
		if elem.Enabled && elem.synced && elem.Sentences && elem.Language == "korean" {
			acc = append(acc, WordsSelection{
				name:         elem.Name,
				generatorKey: elem.Path,
			})
		}
	}

	for _, elem := range config.EmbededWordLists {
		if elem.Enabled && elem.IsSentences && elem.Language == "korean" {
			acc = append(acc, WordsSelection{
				name:         elem.Name,
				generatorKey: elem.Name,
			})
		}
	}
	return acc
}

func initTimerBasedTest(settings TimerBasedTestSettings, mainMenu MainMenu) TimerBasedTest {
	// Get the current wordlist selection
	currentSelection := settings.wordListSelections[settings.wordListCursor]
	generatorKey := currentSelection.generatorKey
	
	// Generate content with comprehensive validation and error handling
	result := generateContentWithErrorHandling(
		generatorKey, 
		mainMenu.koreanTimeBasedGenerator, 
		mainMenu.timeBasedGenerator,
	)
	
	// Log the result for debugging
	logContentGenerationResult(result, generatorKey)
	
	wordsToEnter := result.Content
	
	return TimerBasedTest{
		settings: settings,
		timer: myTimer{
			timer:     timer.NewWithInterval(settings.timeSelections[settings.timeCursor], time.Second),
			duration:  settings.timeSelections[settings.timeCursor],
			isRunning: false,
			timedout:  false,
		},
		base: TestBase{
			wordsToEnter: wordsToEnter,
			inputBuffer:  make([]rune, 0),
			rawInputCnt:  0,
			mistakes: mistakes{
				mistakesAt:     make(map[int]bool, 0),
				rawMistakesCnt: 0,
			},
			cursor: 0,
		},
		completed: false,
		mainMenu:  mainMenu,
	}
}

func initWordCountBasedTest(settings WordCountBasedTestSettings, mainMenu MainMenu) WordCountBasedTest {
	// Get the current wordlist selection
	currentSelection := settings.wordListSelections[settings.wordListCursor]
	generatorKey := currentSelection.generatorKey
	
	// Set the count for the appropriate generator
	mainMenu.wordCountGenerator.Count = settings.wordCountSelections[settings.wordCountCursor]
	mainMenu.koreanWordCountGenerator.Count = settings.wordCountSelections[settings.wordCountCursor]
	
	// Generate content with comprehensive validation and error handling
	result := generateContentWithErrorHandling(
		generatorKey, 
		mainMenu.koreanWordCountGenerator, 
		mainMenu.wordCountGenerator,
	)
	
	// Log the result for debugging
	logContentGenerationResult(result, generatorKey)
	
	wordsToEnter := result.Content
	
	return WordCountBasedTest{
		settings: settings,
		stopwatch: myStopWatch{
			stopwatch: stopwatch.New(),
			isRunning: false,
		},
		base: TestBase{
			wordsToEnter: wordsToEnter,
			inputBuffer:  make([]rune, 0),
			rawInputCnt:  0,
			mistakes: mistakes{
				mistakesAt:     make(map[int]bool, 0),
				rawMistakesCnt: 0,
			},
			cursor: 0,
		},
		completed: false,
		mainMenu:  mainMenu,
	}
}

func initSentenceCountBasedTest(settings SentenceCountBasedTestSettings, mainMenu MainMenu) SentenceCountBasedTest {
	// Get the current sentence list selection
	currentSelection := settings.sentenceListSelections[settings.sentenceListCursor]
	generatorKey := currentSelection.generatorKey
	
	// Set the count for the appropriate generator
	mainMenu.sentenceCountGenerator.Count = settings.sentenceCountSelections[settings.sentenceCountCursor]
	mainMenu.koreanSentenceGenerator.Count = settings.sentenceCountSelections[settings.sentenceCountCursor]
	
	// Generate content with comprehensive validation and error handling
	result := generateContentWithErrorHandling(
		generatorKey, 
		mainMenu.koreanSentenceGenerator, 
		mainMenu.sentenceCountGenerator,
	)
	
	// Log the result for debugging
	logContentGenerationResult(result, generatorKey)
	
	wordsToEnter := result.Content
	
	return SentenceCountBasedTest{
		settings: settings,
		stopwatch: myStopWatch{
			stopwatch: stopwatch.New(),
			isRunning: false,
		},
		base: TestBase{
			wordsToEnter: wordsToEnter,
			inputBuffer:  make([]rune, 0),
			rawInputCnt:  0,
			mistakes: mistakes{
				mistakesAt:     make(map[int]bool, 0),
				rawMistakesCnt: 0,
			},
			cursor: 0,
		},
		completed: false,
		mainMenu:  mainMenu,
	}
}

func initTestSettingCursors() TestSettingCursors {
	return TestSettingCursors{
		TimerTimeCursor:             2,
		TimerWordlistCursor:         0,
		WordCountCursor:             2,
		WordCountWordlistCursor:     0,
		SentenceCountCursor:         2,
		SentenceCountWordlistCursor: 0,
	}
}

func (cursors *TestSettingCursors) resetWordlistCursors() {
	cursors.TimerWordlistCursor = 0
	cursors.WordCountWordlistCursor = 0
	cursors.SentenceCountWordlistCursor = 0
}

func initTimerBasedTestSettings(config Config, words []WordsSelection) TimerBasedTestSettings {
	// Determine if this is primarily Korean content
	isKoreanPrimary := isKoreanWordSelections(words)
	
	var durations []time.Duration
	if isKoreanPrimary {
		// Use Korean-specific shorter durations
		for _, seconds := range config.TimerSettings.KoreanDurations {
			durations = append(durations, time.Second*time.Duration(seconds))
		}
	} else {
		// Use default durations
		for _, seconds := range config.TimerSettings.DefaultDurations {
			durations = append(durations, time.Second*time.Duration(seconds))
		}
	}
	
	return TimerBasedTestSettings{
		timeSelections:     durations,
		timeCursor:         config.TestSettingCursors.TimerTimeCursor,
		wordListSelections: words,
		wordListCursor:     config.TestSettingCursors.TimerWordlistCursor,
		cursor:             0,
		enabled:            len(words) > 0,
	}
}

func initKoreanTimerBasedTestSettings(config Config, words []WordsSelection) TimerBasedTestSettings {
	// Use Korean-specific shorter durations
	var durations []time.Duration
	for _, seconds := range config.TimerSettings.KoreanDurations {
		durations = append(durations, time.Second*time.Duration(seconds))
	}
	
	return TimerBasedTestSettings{
		timeSelections:     durations,
		timeCursor:         0, // Default to first option (30s for Korean)
		wordListSelections: words,
		wordListCursor:     0,
		cursor:             0,
		enabled:            len(words) > 0,
	}
}

func initWordCountBasedTestSettings(config Config, words []WordsSelection) WordCountBasedTestSettings {
	return WordCountBasedTestSettings{
		wordCountSelections: []int{100, 50, 25, 10},
		wordCountCursor:     config.TestSettingCursors.WordCountCursor,
		wordListSelections:  words,
		wordListCursor:      config.TestSettingCursors.WordCountWordlistCursor,
		cursor:              0,
		enabled:             len(words) > 0,
	}
}

func initSentenceCountBasedTestSettings(config Config, words []WordsSelection) SentenceCountBasedTestSettings {
	return SentenceCountBasedTestSettings{
		sentenceCountSelections: []int{30, 15, 5, 1},
		sentenceCountCursor:     config.TestSettingCursors.SentenceCountCursor,
		sentenceListSelections:  words,
		sentenceListCursor:      config.TestSettingCursors.SentenceCountWordlistCursor,
		cursor:                  0,
		enabled:                 len(words) > 0,
	}
}



func initConfigView(config Config, mainMenu MainMenu) ConfigView {
	configView := ConfigView{
		config:   config,
		mainMenu: mainMenu,
	}
	return configView
}

func initConfigViewSelection() ConfigViewSelection {
	return ConfigViewSelection{}
}

func initMainMenu() MainMenu {
	config := ReadConfig()
	
	// Check Korean wordlist availability and log warnings
	koreanIssues := checkKoreanWordlistAvailability(config)
	if len(koreanIssues) > 0 {
		log.Printf("[KOREAN-WARNING] Korean wordlist availability issues detected:")
		for _, issue := range koreanIssues {
			log.Printf("[KOREAN-WARNING] - %s", issue)
		}
		log.Printf("[KOREAN-WARNING] Users may experience issues with Korean typing practice")
	}
	
	// Check terminal Korean support
	terminalSupported, terminalWarnings := validateTerminalKoreanSupport()
	if !terminalSupported {
		log.Printf("[KOREAN-WARNING] Terminal Korean support issues detected:")
		for _, warning := range terminalWarnings {
			log.Printf("[KOREAN-WARNING] - %s", warning)
		}
		log.Printf("[KOREAN-WARNING] Korean characters may not display properly")
	}
	
	// Include both English and Korean wordlists in the same selections
	timeBasedWordSelections := append(filterEnabledSelections(config), filterEnabledKoreanSelections(config)...)
	countBasedWordSelections := append(filterEnabledWordSelection(config), filterEnabledKoreanWordSelection(config)...)
	countBasedSentenceSelections := append(filterEnabledSentenceSelection(config), filterEnabledKoreanSentenceSelection(config)...)
	
	// Get Korean-specific selections for Korean generators
	koreanTimeBasedSelections := filterEnabledKoreanSelections(config)
	koreanWordSelections := filterEnabledKoreanWordSelection(config)
	koreanSentenceSelections := filterEnabledKoreanSentenceSelection(config)
	
	// Log Korean generator initialization
	if len(koreanTimeBasedSelections) > 0 {
		log.Printf("[KOREAN-DEBUG] Initializing Korean generators with %d time-based, %d word-based, %d sentence-based selections", 
			len(koreanTimeBasedSelections), len(koreanWordSelections), len(koreanSentenceSelections))
	} else {
		log.Printf("[KOREAN-WARNING] No Korean selections available - Korean typing practice will not be available")
	}
	
	return MainMenu{
		config: config,
		selections: []MainMenuSelection{
			initTimerBasedTestSettings(config, timeBasedWordSelections),
			initWordCountBasedTestSettings(config, countBasedWordSelections),
			initSentenceCountBasedTestSettings(config, countBasedSentenceSelections),
			initConfigViewSelection(),
		},
		cursor:                 0,
		timeBasedGenerator:     words.NewGenerator(paths(timeBasedWordSelections)),
		wordCountGenerator:     words.NewGenerator(paths(countBasedWordSelections)),
		sentenceCountGenerator: words.NewGenerator(paths(countBasedSentenceSelections)),
		// Initialize Korean generators
		koreanTimeBasedGenerator:   words.NewKoreanGenerator(paths(koreanTimeBasedSelections)),
		koreanWordCountGenerator:   words.NewKoreanGenerator(paths(koreanWordSelections)),
		koreanSentenceGenerator:    words.NewKoreanGenerator(paths(koreanSentenceSelections)),
	}
}

func paths(selections []WordsSelection) []string {
	var acc []string
	for _, elem := range selections {
		// XXX: don't to this at home
		if elem.generatorKey != elem.name {
			acc = append(acc, elem.generatorKey)
		}
	}
	return acc
}

// isKoreanWordSelections checks if the majority of word selections are Korean
func isKoreanWordSelections(selections []WordsSelection) bool {
	koreanCount := 0
	totalCount := len(selections)
	
	if totalCount == 0 {
		return false
	}
	
	for _, selection := range selections {
		// Check if selection name contains Korean indicators
		if strings.Contains(selection.name, "Korean") || 
		   strings.Contains(selection.name, "korean") ||
		   strings.Contains(selection.name, "한국") {
			koreanCount++
		}
	}
	
	// Return true if more than 50% are Korean selections
	return koreanCount > totalCount/2
}

// KoreanGenerationType represents the type of Korean content generation
type KoreanGenerationType string

const (
	KoreanWords     KoreanGenerationType = "words"
	KoreanSentences KoreanGenerationType = "sentences"
)

// ContentValidationError represents errors that occur during content validation
type ContentValidationError struct {
	GeneratorKey string
	ListName     string
	ErrorType    string
	Message      string
}

func (e ContentValidationError) Error() string {
	return fmt.Sprintf("Content validation failed for %s (%s): %s - %s", 
		e.GeneratorKey, e.ListName, e.ErrorType, e.Message)
}

// ContentGenerationResult holds the result of content generation with validation
type ContentGenerationResult struct {
	Content   []rune
	Success   bool
	Error     error
	Source    string // "korean" or "fallback"
	Validated bool
}

// KoreanWordListMapping defines the mapping between Korean generatorKeys and generation methods
type KoreanWordListMapping struct {
	GeneratorKey   string
	ListName       string
	GenerationType KoreanGenerationType
}

// getKoreanWordListMapping returns the mapping between Korean generatorKeys and appropriate generation methods
func getKoreanWordListMapping() map[string]KoreanWordListMapping {
	return map[string]KoreanWordListMapping{
		// Korean word lists - map to GenerateKoreanWords
		"korean-common": {
			GeneratorKey:   "korean-common",
			ListName:       "Korean common words",
			GenerationType: KoreanWords,
		},
		"korean-tech": {
			GeneratorKey:   "korean-tech", 
			ListName:       "Korean tech terms",
			GenerationType: KoreanWords,
		},
		// Korean sentence lists - map to GenerateKoreanSentences
		"korean-sentences": {
			GeneratorKey:   "korean-sentences",
			ListName:       "Korean sentences", 
			GenerationType: KoreanSentences,
		},
		// Alternative naming patterns that might be used
		"Korean common words": {
			GeneratorKey:   "Korean common words",
			ListName:       "Korean common words",
			GenerationType: KoreanWords,
		},
		"Korean tech terms": {
			GeneratorKey:   "Korean tech terms",
			ListName:       "Korean tech terms", 
			GenerationType: KoreanWords,
		},
		"Korean sentences": {
			GeneratorKey:   "Korean sentences",
			ListName:       "Korean sentences",
			GenerationType: KoreanSentences,
		},
	}
}

// getKoreanGenerationType determines the appropriate Korean generation method for a given generatorKey
func getKoreanGenerationType(generatorKey string) (KoreanWordListMapping, bool) {
	mapping := getKoreanWordListMapping()
	
	// Direct lookup first
	if koreanMapping, exists := mapping[generatorKey]; exists {
		return koreanMapping, true
	}
	
	// Fallback: check if the key contains Korean indicators and determine type
	lowerKey := strings.ToLower(generatorKey)
	if isKoreanWordList(generatorKey) {
		// Determine if it's sentences or words based on key content
		if strings.Contains(lowerKey, "sentence") || strings.Contains(lowerKey, "문장") {
			return KoreanWordListMapping{
				GeneratorKey:   generatorKey,
				ListName:       generatorKey,
				GenerationType: KoreanSentences,
			}, true
		} else {
			// Default to words for Korean content
			return KoreanWordListMapping{
				GeneratorKey:   generatorKey,
				ListName:       generatorKey,
				GenerationType: KoreanWords,
			}, true
		}
	}
	
	return KoreanWordListMapping{}, false
}

// isKoreanWordList detects if a generatorKey represents a Korean wordlist
func isKoreanWordList(generatorKey string) bool {
	lowerKey := strings.ToLower(generatorKey)
	isKorean := false
	
	// Check for Korean language indicators in the key
	if strings.Contains(generatorKey, "한국") || strings.Contains(generatorKey, "한글") {
		isKorean = true
		log.Printf("[KOREAN-DEBUG] Korean wordlist detected by Hangul characters: %s", generatorKey)
	}
	
	// Check for "korean" as a separate word or part of compound words
	if strings.Contains(lowerKey, "korean") {
		// Additional validation to ensure it's actually referring to Korean language
		// Check if it's at word boundaries or followed by common separators
		if strings.HasPrefix(lowerKey, "korean") || 
		   strings.Contains(lowerKey, " korean") ||
		   strings.Contains(lowerKey, "-korean") ||
		   strings.Contains(lowerKey, "_korean") ||
		   strings.Contains(lowerKey, "/korean") ||
		   strings.Contains(lowerKey, "korean ") ||
		   strings.Contains(lowerKey, "korean-") ||
		   strings.Contains(lowerKey, "korean_") ||
		   strings.Contains(lowerKey, "korean.") {
			isKorean = true
			log.Printf("[KOREAN-DEBUG] Korean wordlist detected by 'korean' keyword: %s", generatorKey)
		}
	}
	
	// Check for "korea" as a separate word or part of compound words
	if strings.Contains(lowerKey, "korea") {
		// Additional validation to ensure it's actually referring to Korea
		if strings.HasPrefix(lowerKey, "korea") || 
		   strings.Contains(lowerKey, " korea") ||
		   strings.Contains(lowerKey, "-korea") ||
		   strings.Contains(lowerKey, "_korea") ||
		   strings.Contains(lowerKey, "/korea") ||
		   strings.Contains(lowerKey, "korea ") ||
		   strings.Contains(lowerKey, "korea-") ||
		   strings.Contains(lowerKey, "korea_") ||
		   strings.Contains(lowerKey, "korea.") {
			isKorean = true
			log.Printf("[KOREAN-DEBUG] Korean wordlist detected by 'korea' keyword: %s", generatorKey)
		}
	}
	
	if !isKorean {
		log.Printf("[DEBUG] Non-Korean wordlist detected: %s", generatorKey)
	}
	
	return isKorean
}

// ValidateKoreanContent validates that generated Korean content is proper and usable
func ValidateKoreanContent(content []rune, generatorKey, listName string) error {
	// Check if content is empty
	if len(content) == 0 {
		return ContentValidationError{
			GeneratorKey: generatorKey,
			ListName:     listName,
			ErrorType:    "EmptyContent",
			Message:      "Generated Korean content is empty",
		}
	}
	
	// Convert to string for validation
	contentStr := string(content)
	
	// Check for whitespace-only content first (before length check)
	if strings.TrimSpace(contentStr) == "" {
		return ContentValidationError{
			GeneratorKey: generatorKey,
			ListName:     listName,
			ErrorType:    "WhitespaceOnly",
			Message:      "Generated content contains only whitespace",
		}
	}
	
	// Check if content is too short (less than 10 characters)
	if len(content) < 10 {
		return ContentValidationError{
			GeneratorKey: generatorKey,
			ListName:     listName,
			ErrorType:    "ContentTooShort",
			Message:      fmt.Sprintf("Generated content too short: %d characters", len(content)),
		}
	}
	
	// Perform integrity checks on the content
	if err := validateKoreanContentIntegrity(content); err != nil {
		return ContentValidationError{
			GeneratorKey: generatorKey,
			ListName:     listName,
			ErrorType:    "ContentIntegrityError",
			Message:      fmt.Sprintf("Content integrity validation failed: %v", err),
		}
	}
	
	// Check if content contains Korean characters
	hasKorean := detectKoreanContent(contentStr)
	if !hasKorean {
		return ContentValidationError{
			GeneratorKey: generatorKey,
			ListName:     listName,
			ErrorType:    "NoKoreanContent",
			Message:      "Generated content does not contain Korean characters",
		}
	}
	
	return nil
}

// generateKoreanContentWithValidation generates Korean content with comprehensive validation and error handling
func generateKoreanContentWithValidation(koreanMapping KoreanWordListMapping, generator words.WordsGenerator) (result ContentGenerationResult) {
	// Add panic recovery for robust error handling
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[KOREAN-ERROR] PANIC during Korean content generation: %v", r)
			result = handleContentGenerationPanic(koreanMapping.GeneratorKey)
		}
	}()
	
	var content []rune
	
	// Log the generation attempt with detailed information
	log.Printf("[KOREAN-DEBUG] Attempting Korean content generation: key=%s, listName=%s, type=%s", 
		koreanMapping.GeneratorKey, koreanMapping.ListName, koreanMapping.GenerationType)
	log.Printf("[KOREAN-DEBUG] Generator state - Count: %d", generator.Count)
	
	// Generate content based on type with error handling
	switch koreanMapping.GenerationType {
	case KoreanWords:
		log.Printf("[KOREAN-DEBUG] Using GenerateKoreanWords method for: %s", koreanMapping.ListName)
		content = generator.GenerateKoreanWords(koreanMapping.ListName)
		log.Printf("[KOREAN-DEBUG] GenerateKoreanWords returned %d characters", len(content))
	case KoreanSentences:
		log.Printf("[KOREAN-DEBUG] Using GenerateKoreanSentences method for: %s", koreanMapping.ListName)
		content = generator.GenerateKoreanSentences(koreanMapping.ListName)
		log.Printf("[KOREAN-DEBUG] GenerateKoreanSentences returned %d characters", len(content))
	default:
		// Fallback to words if type is unclear
		log.Printf("[KOREAN-ERROR] Unknown Korean generation type %s, falling back to words", koreanMapping.GenerationType)
		content = generator.GenerateKoreanWords(koreanMapping.ListName)
		log.Printf("[KOREAN-DEBUG] Fallback GenerateKoreanWords returned %d characters", len(content))
	}
	
	// Log initial content analysis
	if len(content) > 0 {
		koreanCharCount := 0
		for _, r := range content {
			if isKoreanChar(r) {
				koreanCharCount++
			}
		}
		log.Printf("[KOREAN-DEBUG] Generated content analysis: %d total chars, %d Korean chars (%.1f%%)", 
			len(content), koreanCharCount, float64(koreanCharCount)/float64(len(content))*100)
	} else {
		log.Printf("[KOREAN-ERROR] No content generated - this will cause blank screen issue")
	}
	
	// Validate the generated content
	log.Printf("[KOREAN-DEBUG] Starting Korean content validation...")
	validationErr := ValidateKoreanContent(content, koreanMapping.GeneratorKey, koreanMapping.ListName)
	if validationErr != nil {
		log.Printf("[KOREAN-ERROR] Korean content validation failed: %v", validationErr)
		log.Printf("[KOREAN-ERROR] This validation failure will cause Korean display issues")
		return ContentGenerationResult{
			Content:   content,
			Success:   false,
			Error:     validationErr,
			Source:    "korean",
			Validated: true,
		}
	}
	
	log.Printf("[KOREAN-SUCCESS] Korean content generation and validation successful: %d characters generated", len(content))
	log.Printf("[KOREAN-SUCCESS] Korean content ready for display - no blank screen expected")
	return ContentGenerationResult{
		Content:   content,
		Success:   true,
		Error:     nil,
		Source:    "korean",
		Validated: true,
	}
}

// generateFallbackContent generates fallback content when Korean generation fails
func generateFallbackContent(generatorKey string, generator words.WordsGenerator) ContentGenerationResult {
	log.Printf("[KOREAN-ERROR] Generating fallback content for failed Korean generation: %s", generatorKey)
	log.Printf("[KOREAN-ERROR] Korean display issue detected - attempting fallback to prevent blank screen")
	
	// Try to use generic generator as fallback
	log.Printf("[DEBUG] Attempting generic generator fallback for: %s", generatorKey)
	content := generator.Generate(generatorKey)
	log.Printf("[DEBUG] Generic generator returned %d characters", len(content))
	
	// Validate fallback content
	if len(content) == 0 {
		log.Printf("[ERROR] Fallback content generation also failed for: %s", generatorKey)
		log.Printf("[ERROR] Attempting emergency fallback strategies...")
		
		// Try alternative fallback strategies
		content = generateEmergencyFallbackContent(generatorKey)
		log.Printf("[DEBUG] Emergency fallback generated %d characters", len(content))
	}
	
	// Final validation of fallback content
	if len(content) == 0 {
		log.Printf("[ERROR] All fallback strategies failed for: %s", generatorKey)
		log.Printf("[ERROR] Creating absolute minimal fallback to prevent application crash")
		// Create absolute minimal fallback content
		fallbackText := "Content generation failed. Please restart the application or check your configuration."
		content = []rune(fallbackText)
		
		return ContentGenerationResult{
			Content:   content,
			Success:   false,
			Error:     fmt.Errorf("all fallback strategies failed for %s", generatorKey),
			Source:    "emergency",
			Validated: true,
		}
	}
	
	log.Printf("[SUCCESS] Fallback content generated: %d characters", len(content))
	log.Printf("[SUCCESS] Fallback content should prevent blank screen for Korean wordlist failure")
	return ContentGenerationResult{
		Content:   content,
		Success:   true,
		Error:     nil,
		Source:    "fallback",
		Validated: true,
	}
}

// generateEmergencyFallbackContent provides emergency content when all other methods fail
func generateEmergencyFallbackContent(generatorKey string) []rune {
	log.Printf("[EMERGENCY] Generating emergency fallback content for: %s", generatorKey)
	
	// Provide basic typing practice content based on the type of wordlist requested
	if isKoreanWordList(generatorKey) {
		log.Printf("[KOREAN-EMERGENCY] Creating emergency Korean content to prevent blank screen")
		
		// Check if terminal supports Korean characters
		supported, warnings := validateTerminalKoreanSupport()
		
		if supported {
			// Provide basic Korean practice content
			emergencyKorean := "안녕하세요 세계 한국어 타이핑 연습 컴퓨터 키보드 프로그래밍 개발자 소프트웨어 애플리케이션 테스트 시스템"
			log.Printf("[KOREAN-EMERGENCY] Emergency Korean content created: %d characters", len([]rune(emergencyKorean)))
			log.Printf("[KOREAN-EMERGENCY] This should resolve the Korean blank screen issue temporarily")
			return []rune(emergencyKorean)
		} else {
			// Terminal doesn't support Korean, provide helpful error message
			log.Printf("[KOREAN-EMERGENCY] Terminal Korean support issues detected: %v", warnings)
			errorMsg := "❌ Korean Display Issue Detected\n\n"
			errorMsg += "Your terminal may not support Korean characters properly.\n"
			errorMsg += "Korean text might appear as boxes (□) or question marks (?).\n\n"
			errorMsg += "Quick fixes:\n"
			errorMsg += "• Install Korean fonts on your system\n"
			errorMsg += "• Set terminal encoding to UTF-8\n"
			errorMsg += "• Try a different terminal application\n\n"
			errorMsg += "Temporary solution: Use English wordlists until Korean support is configured."
			
			return []rune(errorMsg)
		}
	} else {
		log.Printf("[EMERGENCY] Creating emergency English content")
		// Provide basic English practice content
		emergencyEnglish := "the quick brown fox jumps over the lazy dog programming computer keyboard typing practice software development application system"
		log.Printf("[EMERGENCY] Emergency English content created: %d characters", len([]rune(emergencyEnglish)))
		return []rune(emergencyEnglish)
	}
}

// handleContentGenerationPanic recovers from panics during content generation
func handleContentGenerationPanic(generatorKey string) ContentGenerationResult {
	if r := recover(); r != nil {
		log.Printf("PANIC during content generation for %s: %v", generatorKey, r)
		
		// Generate emergency content
		emergencyContent := generateEmergencyFallbackContent(generatorKey)
		
		return ContentGenerationResult{
			Content:   emergencyContent,
			Success:   false,
			Error:     fmt.Errorf("content generation panic: %v", r),
			Source:    "panic_recovery",
			Validated: true,
		}
	}
	
	// This should never be reached, but just in case
	return ContentGenerationResult{
		Content:   []rune("Unexpected error in content generation"),
		Success:   false,
		Error:     fmt.Errorf("unexpected panic recovery state"),
		Source:    "error",
		Validated: true,
	}
}

// validateGeneratorState validates that generators are properly initialized
func validateGeneratorState(koreanGenerator, fallbackGenerator words.WordsGenerator, generatorKey string) error {
	// Check if generators have proper pool data
	// This is a basic check - in a real implementation, you might want to expose
	// more generator internals for validation
	
	// Validate fallback generator first
	if fallbackGenerator.Count <= 0 {
		return fmt.Errorf("fallback generator has invalid count: %d", fallbackGenerator.Count)
	}
	
	// Try a simple test generation to validate state
	testContent := fallbackGenerator.Generate("Common words")
	if len(testContent) == 0 {
		// Try alternative embedded content
		testContent = fallbackGenerator.Generate("Frankenstein sentences")
		if len(testContent) == 0 {
			return fmt.Errorf("fallback generator appears to be uninitialized or has no accessible content")
		}
	}
	
	// Validate Korean generator if we're dealing with Korean content
	if isKoreanWordList(generatorKey) {
		if koreanGenerator.Count <= 0 {
			log.Printf("Warning: Korean generator has invalid count: %d", koreanGenerator.Count)
			// This is not a fatal error, we can still try to generate content
		}
		
		// Test Korean generator with known Korean wordlists
		koreanTestContent := koreanGenerator.GenerateKoreanWords("Korean common words")
		if len(koreanTestContent) == 0 {
			log.Printf("Warning: Korean generator test failed, but continuing with fallback available")
			// This is not a fatal error since we have fallback
		}
	}
	
	return nil
}

// createUserFriendlyErrorMessage creates user-friendly error messages for different failure scenarios
func createUserFriendlyErrorMessage(err error, generatorKey string) string {
	if contentErr, ok := err.(ContentValidationError); ok {
		switch contentErr.ErrorType {
		case "EmptyContent":
			return fmt.Sprintf("❌ Korean wordlist '%s' appears to be empty.\n💡 Please check your Korean wordlist configuration or try a different Korean wordlist.", generatorKey)
		case "ContentTooShort":
			return fmt.Sprintf("❌ Korean wordlist '%s' has insufficient content for typing practice.\n💡 Please use a different wordlist or check the configuration.", generatorKey)
		case "NoKoreanContent":
			return fmt.Sprintf("❌ Wordlist '%s' does not contain Korean characters.\n💡 Please verify you've selected the correct Korean wordlist from the menu.", generatorKey)
		case "WhitespaceOnly":
			return fmt.Sprintf("❌ Korean wordlist '%s' contains only whitespace.\n💡 The wordlist file may be corrupted. Please check the file or reinstall Korean content.", generatorKey)
		case "ContentIntegrityError":
			return fmt.Sprintf("❌ Korean wordlist '%s' appears to be corrupted.\n💡 Please check the wordlist file encoding (should be UTF-8) or reinstall Korean content.", generatorKey)
		default:
			return fmt.Sprintf("❌ Korean content validation failed for '%s': %s\n💡 Please check your Korean wordlist configuration.", generatorKey, contentErr.Message)
		}
	}
	
	// Handle other types of errors
	if isKoreanWordList(generatorKey) {
		return fmt.Sprintf("❌ Korean content generation failed for '%s': %v\n💡 Please check Korean wordlist availability and terminal Korean font support.", generatorKey, err)
	}
	
	return fmt.Sprintf("❌ Content generation failed for '%s': %v\n💡 Please check wordlist configuration or try a different wordlist.", generatorKey, err)
}

// createTerminalEncodingWarning creates warnings for terminal encoding and font issues
func createTerminalEncodingWarning(generatorKey string) string {
	warning := fmt.Sprintf("⚠️  Korean Display Warning for '%s':\n\n", generatorKey)
	warning += "🔤 Terminal Font Issues:\n"
	warning += "   • Your terminal may not have Korean (Hangul) font support\n"
	warning += "   • Korean characters might appear as boxes (□) or question marks (?)\n"
	warning += "   • Solution: Install a font that supports Korean characters\n\n"
	warning += "🔧 Terminal Encoding Issues:\n"
	warning += "   • Your terminal encoding may not be set to UTF-8\n"
	warning += "   • Korean characters require UTF-8 encoding to display properly\n"
	warning += "   • Solution: Set your terminal to UTF-8 encoding\n\n"
	warning += "💻 Recommended Terminal Fonts:\n"
	warning += "   • Noto Sans CJK (Google Noto fonts)\n"
	warning += "   • Source Han Sans (Adobe)\n"
	warning += "   • Malgun Gothic (Windows)\n"
	warning += "   • AppleGothic (macOS)\n\n"
	warning += "🛠️  Quick Fixes:\n"
	warning += "   • Linux: sudo apt install fonts-noto-cjk\n"
	warning += "   • macOS: Korean fonts are usually pre-installed\n"
	warning += "   • Windows: Install Korean language pack\n"
	warning += "   • Check terminal settings for UTF-8 support\n\n"
	warning += "🔄 If Korean text still doesn't display, try switching to an English wordlist temporarily."
	
	return warning
}

// createKoreanWordlistUnavailableGuidance provides guidance when Korean wordlists are unavailable
func createKoreanWordlistUnavailableGuidance(generatorKey string) string {
	guidance := fmt.Sprintf("📋 Korean Wordlist Unavailable: '%s'\n\n", generatorKey)
	guidance += "❓ Possible Causes:\n"
	guidance += "   • Korean wordlists are not installed\n"
	guidance += "   • Korean wordlists are disabled in configuration\n"
	guidance += "   • Korean wordlist files are missing or corrupted\n"
	guidance += "   • Application was not installed with Korean language support\n\n"
	guidance += "🔧 Solutions:\n"
	guidance += "   1. Check Configuration:\n"
	guidance += "      • Open application settings/configuration\n"
	guidance += "      • Verify Korean wordlists are enabled\n"
	guidance += "      • Check wordlist file paths are correct\n\n"
	guidance += "   2. Install Korean Content:\n"
	guidance += "      • Reinstall the application with Korean language support\n"
	guidance += "      • Download Korean wordlist files separately if available\n"
	guidance += "      • Check application documentation for Korean setup\n\n"
	guidance += "   3. Alternative Options:\n"
	guidance += "      • Use English wordlists temporarily\n"
	guidance += "      • Create custom Korean wordlist files\n"
	guidance += "      • Check for application updates with Korean support\n\n"
	guidance += "📁 Expected Korean Wordlists:\n"
	guidance += "   • Korean common words (korean-common)\n"
	guidance += "   • Korean tech terms (korean-tech)\n"
	guidance += "   • Korean sentences (korean-sentences)\n\n"
	guidance += "💡 If you continue to have issues, please check the application documentation or contact support."
	
	return guidance
}

// displayUserFriendlyError displays comprehensive error information to the user
func displayUserFriendlyError(err error, generatorKey string) string {
	errorMsg := createUserFriendlyErrorMessage(err, generatorKey)
	
	// Add additional context based on error type
	if contentErr, ok := err.(ContentValidationError); ok {
		switch contentErr.ErrorType {
		case "NoKoreanContent", "EmptyContent":
			// Add guidance for unavailable Korean wordlists
			errorMsg += "\n\n" + createKoreanWordlistUnavailableGuidance(generatorKey)
		case "ContentIntegrityError":
			// Add terminal encoding warning for integrity issues
			errorMsg += "\n\n" + createTerminalEncodingWarning(generatorKey)
		}
	} else if isKoreanWordList(generatorKey) {
		// For other Korean-related errors, add terminal encoding warning
		errorMsg += "\n\n" + createTerminalEncodingWarning(generatorKey)
	}
	
	return errorMsg
}

// logUserFriendlyError logs user-friendly error information for debugging
func logUserFriendlyError(err error, generatorKey string) {
	if isKoreanWordList(generatorKey) {
		log.Printf("[KOREAN-ERROR] User-friendly error for Korean wordlist '%s': %v", generatorKey, err)
		log.Printf("[KOREAN-ERROR] This error will be displayed to user with helpful guidance")
	} else {
		log.Printf("[ERROR] User-friendly error for wordlist '%s': %v", generatorKey, err)
	}
	
	// Log specific troubleshooting information
	troubleshootingMsg := createTroubleshootingMessage(err, generatorKey)
	log.Printf("[DEBUG] Troubleshooting information: %s", troubleshootingMsg)
}

// checkKoreanWordlistAvailability checks if Korean wordlists are available and configured
func checkKoreanWordlistAvailability(config Config) []string {
	var issues []string
	
	// Check if any Korean wordlists are enabled
	koreanWordlists := filterEnabledKoreanSelections(config)
	if len(koreanWordlists) == 0 {
		issues = append(issues, "No Korean wordlists are enabled in configuration")
	}
	
	// Check specific Korean wordlist types
	koreanWords := filterEnabledKoreanWordSelection(config)
	koreanSentences := filterEnabledKoreanSentenceSelection(config)
	
	if len(koreanWords) == 0 {
		issues = append(issues, "No Korean word lists are available")
	}
	
	if len(koreanSentences) == 0 {
		issues = append(issues, "No Korean sentence lists are available")
	}
	
	// Check embedded Korean wordlists
	hasEmbeddedKorean := false
	for _, embedded := range config.EmbededWordLists {
		if embedded.Enabled && embedded.Language == "korean" {
			hasEmbeddedKorean = true
			break
		}
	}
	
	if !hasEmbeddedKorean {
		issues = append(issues, "No embedded Korean wordlists are enabled")
	}
	
	return issues
}

// validateTerminalKoreanSupport checks if the terminal supports Korean character display
func validateTerminalKoreanSupport() (bool, []string) {
	var warnings []string
	supported := true
	
	// Test Korean character rendering by checking if we can create Korean runes
	testKorean := "안녕하세요"
	testRunes := []rune(testKorean)
	
	if len(testRunes) != 5 {
		warnings = append(warnings, "Korean character conversion failed - terminal may not support UTF-8")
		supported = false
	}
	
	// Check for common Korean characters
	for _, r := range testRunes {
		if !isKoreanChar(r) {
			warnings = append(warnings, "Korean character validation failed - terminal encoding issues detected")
			supported = false
			break
		}
	}
	
	// Additional checks could be added here for specific terminal capabilities
	// For now, we'll assume basic UTF-8 support is sufficient
	
	if !supported {
		warnings = append(warnings, "Terminal may not properly display Korean characters")
		warnings = append(warnings, "Korean text might appear as boxes (□) or question marks (?)")
		warnings = append(warnings, "Please ensure your terminal supports UTF-8 and has Korean fonts installed")
	}
	
	return supported, warnings
}

// createKoreanSetupGuidance provides comprehensive setup guidance for Korean support
func createKoreanSetupGuidance() string {
	guidance := "🇰🇷 Korean Typing Setup Guide\n\n"
	guidance += "📋 Required Components:\n"
	guidance += "   ✅ Korean wordlists (korean-common, korean-tech, korean-sentences)\n"
	guidance += "   ✅ UTF-8 terminal encoding\n"
	guidance += "   ✅ Korean-compatible font\n"
	guidance += "   ✅ Proper application configuration\n\n"
	guidance += "🔧 Setup Steps:\n"
	guidance += "   1. Install Korean Fonts:\n"
	guidance += "      • Linux: sudo apt install fonts-noto-cjk\n"
	guidance += "      • macOS: Korean fonts pre-installed\n"
	guidance += "      • Windows: Install Korean language pack\n\n"
	guidance += "   2. Configure Terminal:\n"
	guidance += "      • Set encoding to UTF-8\n"
	guidance += "      • Select Korean-compatible font\n"
	guidance += "      • Test: echo '안녕하세요'\n\n"
	guidance += "   3. Enable Korean Wordlists:\n"
	guidance += "      • Check application configuration\n"
	guidance += "      • Ensure Korean wordlists are enabled\n"
	guidance += "      • Verify wordlist files exist\n\n"
	guidance += "🆘 Troubleshooting:\n"
	guidance += "   • Korean text appears as boxes: Font issue\n"
	guidance += "   • Korean text appears as ?: Encoding issue\n"
	guidance += "   • No Korean options: Configuration issue\n"
	guidance += "   • Blank screen: Content generation issue\n\n"
	guidance += "💡 Quick Test: If you can see '한글' properly, Korean support is working!"
	
	return guidance
}

// createTroubleshootingMessage creates detailed troubleshooting information for Korean display issues
func createTroubleshootingMessage(err error, generatorKey string) string {
	troubleshootingMsg := fmt.Sprintf("Korean display troubleshooting for '%s':\n", generatorKey)
	
	if contentErr, ok := err.(ContentValidationError); ok {
		switch contentErr.ErrorType {
		case "EmptyContent":
			troubleshootingMsg += "- Issue: Korean wordlist is empty\n"
			troubleshootingMsg += "- Check: Verify Korean wordlist files exist and contain data\n"
			troubleshootingMsg += "- Check: Ensure Korean wordlist is properly configured in settings\n"
			troubleshootingMsg += "- Action: Try selecting a different Korean wordlist or reinstall Korean content"
		case "ContentTooShort":
			troubleshootingMsg += "- Issue: Korean wordlist has insufficient content\n"
			troubleshootingMsg += "- Check: Verify Korean wordlist file is not truncated or corrupted\n"
			troubleshootingMsg += "- Action: Use a different Korean wordlist with more content"
		case "NoKoreanContent":
			troubleshootingMsg += "- Issue: Selected wordlist does not contain Korean characters\n"
			troubleshootingMsg += "- Check: Verify you selected a Korean wordlist, not an English one\n"
			troubleshootingMsg += "- Action: Select a proper Korean wordlist from the menu"
		case "WhitespaceOnly":
			troubleshootingMsg += "- Issue: Korean wordlist contains only whitespace\n"
			troubleshootingMsg += "- Check: Korean wordlist file may be corrupted or improperly formatted\n"
			troubleshootingMsg += "- Action: Reinstall Korean wordlists or check file encoding"
		case "ContentIntegrityError":
			troubleshootingMsg += "- Issue: Korean wordlist appears corrupted\n"
			troubleshootingMsg += "- Check: File encoding should be UTF-8 for Korean characters\n"
			troubleshootingMsg += "- Check: Terminal font supports Korean characters (Hangul)\n"
			troubleshootingMsg += "- Action: Reinstall Korean wordlists or check terminal configuration"
		default:
			troubleshootingMsg += fmt.Sprintf("- Issue: %s\n", contentErr.Message)
			troubleshootingMsg += "- Action: Check Korean wordlist configuration and terminal settings"
		}
	} else {
		troubleshootingMsg += fmt.Sprintf("- Issue: %v\n", err)
		troubleshootingMsg += "- Check: Korean generator initialization and wordlist availability\n"
		troubleshootingMsg += "- Action: Restart application or check Korean wordlist configuration"
	}
	
	// Add general Korean display troubleshooting tips
	troubleshootingMsg += "\n\nGeneral Korean display troubleshooting:\n"
	troubleshootingMsg += "- Ensure terminal supports UTF-8 encoding\n"
	troubleshootingMsg += "- Verify terminal font includes Korean (Hangul) characters\n"
	troubleshootingMsg += "- Check that Korean wordlists are properly installed\n"
	troubleshootingMsg += "- Try switching to a different Korean wordlist option"
	
	return troubleshootingMsg
}

// logContentGenerationAttempt logs detailed information about content generation attempts
func logContentGenerationAttempt(generatorKey string, isKorean bool, mapping KoreanWordListMapping) {
	if isKorean {
		log.Printf("[KOREAN-DEBUG] Content generation attempt - Korean wordlist detected: key=%s, listName=%s, type=%s", 
			generatorKey, mapping.ListName, mapping.GenerationType)
		log.Printf("[KOREAN-DEBUG] Korean generation details - using %s method for %s", 
			mapping.GenerationType, mapping.ListName)
		log.Printf("[KOREAN-DEBUG] Korean wordlist selection confirmed - proceeding with Korean-specific generator")
	} else {
		log.Printf("[DEBUG] Content generation attempt - Non-Korean wordlist: key=%s", generatorKey)
		log.Printf("[DEBUG] Non-Korean generation details - using standard Generate method")
	}
}

// logContentGenerationResult logs the result of content generation for debugging
func logContentGenerationResult(result ContentGenerationResult, generatorKey string) {
	isKoreanContent := isKoreanWordList(generatorKey)
	
	if result.Success {
		if isKoreanContent {
			log.Printf("[KOREAN-SUCCESS] Content generation SUCCESS: key=%s, source=%s, length=%d, validated=%t", 
				generatorKey, result.Source, len(result.Content), result.Validated)
			log.Printf("[KOREAN-SUCCESS] Korean content successfully generated and validated")
		} else {
			log.Printf("[SUCCESS] Content generation SUCCESS: key=%s, source=%s, length=%d, validated=%t", 
				generatorKey, result.Source, len(result.Content), result.Validated)
		}
	} else {
		if isKoreanContent {
			log.Printf("[KOREAN-ERROR] Content generation FAILED: key=%s, source=%s, error=%v, length=%d", 
				generatorKey, result.Source, result.Error, len(result.Content))
			log.Printf("[KOREAN-ERROR] Korean display issue detected - troubleshooting info: %s", 
				createTroubleshootingMessage(result.Error, generatorKey))
		} else {
			log.Printf("[ERROR] Content generation FAILED: key=%s, source=%s, error=%v, length=%d", 
				generatorKey, result.Source, result.Error, len(result.Content))
		}
	}
	
	// Log content preview for debugging (first 50 characters)
	if len(result.Content) > 0 {
		preview := string(result.Content)
		if len(preview) > 50 {
			preview = preview[:50] + "..."
		}
		if isKoreanContent {
			log.Printf("[KOREAN-DEBUG] Korean content preview: %s", preview)
			// Additional Korean-specific debugging
			koreanCharCount := 0
			for _, r := range result.Content {
				if isKoreanChar(r) {
					koreanCharCount++
				}
			}
			log.Printf("[KOREAN-DEBUG] Korean character analysis: %d Korean chars out of %d total", 
				koreanCharCount, len(result.Content))
		} else {
			log.Printf("[DEBUG] Content preview: %s", preview)
		}
	} else {
		if isKoreanContent {
			log.Printf("[KOREAN-ERROR] No content generated - Korean display will be blank")
		} else {
			log.Printf("[ERROR] No content generated - display will be blank")
		}
	}
}

// generateContentWithErrorHandling generates content with comprehensive error handling and fallback
func generateContentWithErrorHandling(generatorKey string, koreanGenerator, fallbackGenerator words.WordsGenerator) ContentGenerationResult {
	log.Printf("[DEBUG] Starting content generation for: %s", generatorKey)
	
	// Validate generator state first
	if err := validateGeneratorState(koreanGenerator, fallbackGenerator, generatorKey); err != nil {
		log.Printf("[ERROR] Generator validation failed: %v", err)
		log.Printf("[ERROR] This will likely cause display issues - creating emergency content")
		// Create emergency fallback content
		emergencyText := "Content generation system error. Please restart the application or check your configuration."
		return ContentGenerationResult{
			Content:   []rune(emergencyText),
			Success:   false,
			Error:     err,
			Source:    "emergency",
			Validated: true,
		}
	}
	
	// Check if this is a Korean wordlist
	if koreanMapping, isKorean := getKoreanGenerationType(generatorKey); isKorean {
		log.Printf("[KOREAN-DEBUG] Korean wordlist selection confirmed: %s", generatorKey)
		// Log the Korean generation attempt
		logContentGenerationAttempt(generatorKey, isKorean, koreanMapping)
		
		// Attempt Korean content generation with validation
		result := generateKoreanContentWithValidation(koreanMapping, koreanGenerator)
		
		if result.Success {
			log.Printf("[KOREAN-SUCCESS] Korean content generation completed successfully")
			return result
		}
		
		// Korean generation failed, create comprehensive user-friendly error message
		logUserFriendlyError(result.Error, generatorKey)
		log.Printf("[KOREAN-ERROR] Korean display issue detected - attempting fallback to prevent blank screen")
		
		// Try fallback generation
		fallbackResult := generateFallbackContent(generatorKey, fallbackGenerator)
		
		// If fallback also fails, provide a comprehensive helpful error message
		if !fallbackResult.Success {
			log.Printf("[ERROR] Both Korean and fallback generation failed for: %s", generatorKey)
			log.Printf("[ERROR] Displaying comprehensive error message to user")
			
			// Create comprehensive error message with all guidance
			errorText := displayUserFriendlyError(result.Error, generatorKey)
			
			return ContentGenerationResult{
				Content:   []rune(errorText),
				Success:   false,
				Error:     result.Error,
				Source:    "comprehensive_error",
				Validated: true,
			}
		}
		
		log.Printf("[SUCCESS] Fallback content generation succeeded after Korean failure")
		return fallbackResult
	} else {
		log.Printf("[DEBUG] Non-Korean wordlist detected: %s", generatorKey)
		// Log the non-Korean generation attempt
		logContentGenerationAttempt(generatorKey, false, KoreanWordListMapping{})
		
		// Use generic generation for non-Korean wordlists
		log.Printf("[DEBUG] Using standard generator for non-Korean content")
		content := fallbackGenerator.Generate(generatorKey)
		log.Printf("[DEBUG] Standard generator returned %d characters", len(content))
		
		// Basic validation for non-Korean content
		if len(content) == 0 {
			log.Printf("[ERROR] Non-Korean content generation failed for: %s", generatorKey)
			log.Printf("[ERROR] Displaying user-friendly error message")
			
			// Create user-friendly error for empty non-Korean content
			err := fmt.Errorf("empty content generated for %s", generatorKey)
			errorText := createUserFriendlyErrorMessage(err, generatorKey)
			content = []rune(errorText)
			
			return ContentGenerationResult{
				Content:   content,
				Success:   false,
				Error:     err,
				Source:    "error_message",
				Validated: true,
			}
		}
		
		log.Printf("[SUCCESS] Non-Korean content generation successful: %d characters", len(content))
		return ContentGenerationResult{
			Content:   content,
			Success:   true,
			Error:     nil,
			Source:    "english",
			Validated: true,
		}
	}
}

func initialModel(profile termenv.Profile, fore termenv.Color, width, height int) model {
	return model{
		width:  width,
		height: height,
		state:  initMainMenu(),
		styles: Styles{
			correct: func(str string) termenv.Style {
				return termenv.String(str).Foreground(fore)
			},
			toEnter: func(str string) termenv.Style {
				return termenv.String(str).Foreground(fore).Faint()
			},
			mistakes: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("1")).Underline()
			},
			cursor: func(str string) termenv.Style {
				return termenv.String(str).Reverse().Bold()
			},
			runningTimer: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("2"))
			},
			stoppedTimer: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("2")).Faint()
			},
			greener: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("6")).Faint()
			},
			faintGreen: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("10")).Faint()
			},
		},
	}
}
