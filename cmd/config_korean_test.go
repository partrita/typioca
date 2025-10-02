package cmd

import (
	"strings"
	"testing"
	
	"github.com/muesli/termenv"
)

func TestKoreanConfigurationDisplay(t *testing.T) {
	// Create a config with Korean embedded word lists
	config := Config{
		EmbededWordLists: []EmbededWordList{
			{"Common words", false, true, "english"},
			{"일반 한글 단어", false, true, "korean"},
			{"한글 기술 용어", false, true, "korean"},
		},
		WordLists: []WordList{
			{Name: "Custom English", Language: "english", Enabled: true, synced: true},
			{Name: "Custom Korean", Language: "korean", Enabled: true, synced: true},
		},
	}

	// Create a config view
	configView := ConfigView{
		config: config,
		cursor: 0,
	}

	// Create a model with the config view
	m := model{
		state:  configView,
		styles: createTestStyles(),
		width:  80,
		height: 24,
	}

	// Get the view output
	view := m.View()

	// Verify that Korean language indicators are displayed
	if !strings.Contains(view, "KR") {
		t.Error("Expected Korean language indicator 'KR' to be displayed in config view")
	}

	// Verify that English language indicators are displayed
	if !strings.Contains(view, "EN") {
		t.Error("Expected English language indicator 'EN' to be displayed in config view")
	}

	// Verify that Korean word list names are displayed
	if !strings.Contains(view, "일반 한글 단어") {
		t.Error("Expected Korean embedded word list '일반 한글 단어' to be displayed")
	}

	if !strings.Contains(view, "한글 기술 용어") {
		t.Error("Expected Korean embedded word list '한글 기술 용어' to be displayed")
	}

	// Verify that custom Korean files are displayed
	if !strings.Contains(view, "Custom Korean") {
		t.Error("Expected custom Korean word list 'Custom Korean' to be displayed")
	}
}

func TestKoreanConfigurationToggle(t *testing.T) {
	// Create a config with Korean embedded word lists
	config := Config{
		EmbededWordLists: []EmbededWordList{
			{"일반 한글 단어", false, false, "korean"}, // Initially disabled
		},
		WordLists: []WordList{
			{Name: "Custom Korean", Language: "korean", Enabled: false, synced: true}, // Initially disabled
		},
		TestSettingCursors: TestSettingCursors{},
	}

	// Test toggling embedded Korean word list
	config.EmbededWordLists[0].toggleEnabled()
	if !config.EmbededWordLists[0].Enabled {
		t.Error("Expected Korean embedded word list to be enabled after toggle")
	}

	// Test toggling custom Korean word list
	config.WordLists[0].toggleEnabled()
	if !config.WordLists[0].Enabled {
		t.Error("Expected Korean custom word list to be enabled after toggle")
	}
}

func TestKoreanConfigurationFiltering(t *testing.T) {
	// Create a config with mixed language word lists
	config := Config{
		EmbededWordLists: []EmbededWordList{
			{"Common words", false, true, "english"},
			{"일반 한글 단어", false, true, "korean"},
			{"한글 기술 용어", false, false, "korean"}, // Disabled
		},
		WordLists: []WordList{
			{Name: "Custom English", Language: "english", Enabled: true, synced: true},
			{Name: "Custom Korean", Language: "korean", Enabled: true, synced: true},
			{Name: "Disabled Korean", Language: "korean", Enabled: false, synced: true},
		},
	}

	// Test Korean word filtering
	koreanSelections := filterEnabledKoreanSelections(config)
	
	// Should include enabled Korean embedded and custom lists
	expectedCount := 2 // "일반 한글 단어" and "Custom Korean"
	if len(koreanSelections) != expectedCount {
		t.Errorf("Expected %d Korean selections, got %d", expectedCount, len(koreanSelections))
	}

	// Verify the correct Korean lists are included
	found := make(map[string]bool)
	for _, selection := range koreanSelections {
		found[selection.name] = true
	}

	if !found["일반 한글 단어"] {
		t.Error("Expected '일반 한글 단어' to be in Korean selections")
	}

	if !found["Custom Korean"] {
		t.Error("Expected 'Custom Korean' to be in Korean selections")
	}

	if found["한글 기술 용어"] {
		t.Error("Did not expect disabled '한글 기술 용어' to be in Korean selections")
	}

	if found["Disabled Korean"] {
		t.Error("Did not expect disabled 'Disabled Korean' to be in Korean selections")
	}
}

func TestDefaultConfigIncludesKoreanWordLists(t *testing.T) {
	config := defaultConfig()
	
	// Check that Korean embedded word lists are included
	koreanEmbeddedCount := 0
	for _, embedded := range config.EmbededWordLists {
		if embedded.Language == "korean" {
			koreanEmbeddedCount++
		}
	}
	
	if koreanEmbeddedCount == 0 {
		t.Error("Expected default config to include Korean embedded word lists")
	}
	
	// Verify specific Korean word lists exist
	found := make(map[string]bool)
	for _, embedded := range config.EmbededWordLists {
		if embedded.Language == "korean" {
			found[embedded.Name] = true
		}
	}
	
	expectedKoreanLists := []string{"Korean common words", "Korean tech terms", "Korean sentences"}
	for _, expected := range expectedKoreanLists {
		if !found[expected] {
			t.Errorf("Expected Korean embedded word list '%s' to be in default config", expected)
		}
	}
}

// Helper function to create test styles
func createTestStyles() Styles {
	return Styles{
		correct:      func(s string) termenv.Style { return termenv.String(s) },
		toEnter:      func(s string) termenv.Style { return termenv.String(s) },
		mistakes:     func(s string) termenv.Style { return termenv.String(s) },
		cursor:       func(s string) termenv.Style { return termenv.String(s) },
		runningTimer: func(s string) termenv.Style { return termenv.String(s) },
		stoppedTimer: func(s string) termenv.Style { return termenv.String(s) },
		greener:      func(s string) termenv.Style { return termenv.String(s) },
		faintGreen:   func(s string) termenv.Style { return termenv.String(s) },
	}
}