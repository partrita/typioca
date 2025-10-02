package cmd

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/muesli/termenv"
)

func TestKoreanResultsDisplay(t *testing.T) {
	// Create a mock model with Korean results
	styles := Styles{
		correct:      func(s string) termenv.Style { return termenv.String(s) },
		toEnter:      func(s string) termenv.Style { return termenv.String(s) },
		mistakes:     func(s string) termenv.Style { return termenv.String(s) },
		cursor:       func(s string) termenv.Style { return termenv.String(s) },
		runningTimer: func(s string) termenv.Style { return termenv.String(s) },
		stoppedTimer: func(s string) termenv.Style { return termenv.String(s) },
		greener:      func(s string) termenv.Style { return termenv.String(s) },
		faintGreen:   func(s string) termenv.Style { return termenv.String(s) },
	}

	m := model{
		styles: styles,
		width:  80,
		height: 24,
	}

	// Test Korean Timer Based Test Results display
	t.Run("KoreanTimerBasedTestResults display", func(t *testing.T) {
		results := Results{
			wpm:      45,
			accuracy: 95.5,
			deltaWpm: 5.2,
			rawWpm:   50,
			time:     1 * time.Minute,
			wordList: "Korean Common Words",
		}

		koreanResults := KoreanTimerBasedTestResults{
			results:       results,
			wpmEachSecond: []float64{40.0, 42.0, 45.0},
		}

		m.state = koreanResults
		view := m.View()

		// Check that Korean indicator is present
		if !strings.Contains(view, "(한글)") {
			t.Error("Korean results display should contain '(한글)' indicator")
		}

		// Check that word list name is displayed
		if !strings.Contains(view, "Korean Common Words") {
			t.Error("Korean results display should contain word list name")
		}

		// Check that WPM is displayed
		if !strings.Contains(view, "wpm: 45") {
			t.Error("Korean results display should contain WPM")
		}

		// Check that accuracy is displayed
		if !strings.Contains(view, "accuracy: 95.5") {
			t.Error("Korean results display should contain accuracy")
		}
	})

	// Test Korean Word Count Based Test Results display
	t.Run("KoreanWordCountTestResults display", func(t *testing.T) {
		results := Results{
			wpm:      38,
			accuracy: 92.3,
			deltaWpm: -2.1,
			rawWpm:   42,
			time:     45 * time.Second,
			wordList: "Korean Tech Terms",
		}

		koreanResults := KoreanWordCountTestResults{
			results:       results,
			wordCnt:       25,
			wpmEachSecond: []float64{35.0, 37.0, 38.0},
		}

		m.state = koreanResults
		view := m.View()

		// Check that Korean indicator is present
		if !strings.Contains(view, "(한글)") {
			t.Error("Korean word count results display should contain '(한글)' indicator")
		}

		// Check that word list name is displayed
		if !strings.Contains(view, "Korean Tech Terms") {
			t.Error("Korean word count results display should contain word list name")
		}

		// Check that word count is displayed
		if !strings.Contains(view, "cnt: 25") {
			t.Error("Korean word count results display should contain word count")
		}

		// Check that "words:" label is used (not "sentences:")
		if !strings.Contains(view, "words:") {
			t.Error("Korean word count results display should contain 'words:' label")
		}
	})

	// Test Korean Sentence Count Based Test Results display
	t.Run("KoreanSentenceCountTestResults display", func(t *testing.T) {
		results := Results{
			wpm:      32,
			accuracy: 88.7,
			deltaWpm: 1.8,
			rawWpm:   36,
			time:     2 * time.Minute,
			wordList: "Korean Sentences",
		}

		koreanResults := KoreanSentenceCountTestResults{
			results:       results,
			sentenceCnt:   10,
			wpmEachSecond: []float64{30.0, 31.0, 32.0},
		}

		m.state = koreanResults
		view := m.View()

		// Check that Korean indicator is present
		if !strings.Contains(view, "(한글)") {
			t.Error("Korean sentence count results display should contain '(한글)' indicator")
		}

		// Check that word list name is displayed
		if !strings.Contains(view, "Korean Sentences") {
			t.Error("Korean sentence count results display should contain word list name")
		}

		// Check that sentence count is displayed
		if !strings.Contains(view, "cnt: 10") {
			t.Error("Korean sentence count results display should contain sentence count")
		}

		// Check that "sentences:" label is used (not "words:")
		if !strings.Contains(view, "sentences:") {
			t.Error("Korean sentence count results display should contain 'sentences:' label")
		}
	})
}

func TestKoreanResultsDisplayFormatting(t *testing.T) {
	// Test that Korean results display properly handles Korean text formatting
	styles := Styles{
		correct:      func(s string) termenv.Style { return termenv.String(s) },
		toEnter:      func(s string) termenv.Style { return termenv.String(s) },
		mistakes:     func(s string) termenv.Style { return termenv.String(s) },
		cursor:       func(s string) termenv.Style { return termenv.String(s) },
		runningTimer: func(s string) termenv.Style { return termenv.String(s) },
		stoppedTimer: func(s string) termenv.Style { return termenv.String(s) },
		greener:      func(s string) termenv.Style { return termenv.String(s) },
		faintGreen:   func(s string) termenv.Style { return termenv.String(s) },
	}

	m := model{
		styles: styles,
		width:  80,
		height: 24,
	}

	// Test with Korean word list name containing Korean characters
	results := Results{
		wpm:      40,
		accuracy: 90.0,
		deltaWpm: 0.0,
		rawWpm:   40,
		time:     1 * time.Minute,
		wordList: "일반 한글 단어", // Korean word list name
	}

	koreanResults := KoreanTimerBasedTestResults{
		results:       results,
		wpmEachSecond: []float64{40.0},
	}

	m.state = koreanResults
	view := m.View()

	// Check that Korean word list name is properly displayed
	if !strings.Contains(view, "일반 한글 단어") {
		t.Error("Korean results display should properly handle Korean word list names")
	}

	// Check that Korean indicator is appended correctly
	if !strings.Contains(view, "일반 한글 단어 (한글)") {
		t.Error("Korean results display should append '(한글)' to Korean word list names")
	}
}

func TestKoreanResultsDisplayMetrics(t *testing.T) {
	// Test that Korean-specific metrics are displayed correctly
	styles := Styles{
		correct:      func(s string) termenv.Style { return termenv.String(s) },
		toEnter:      func(s string) termenv.Style { return termenv.String(s) },
		mistakes:     func(s string) termenv.Style { return termenv.String(s) },
		cursor:       func(s string) termenv.Style { return termenv.String(s) },
		runningTimer: func(s string) termenv.Style { return termenv.String(s) },
		stoppedTimer: func(s string) termenv.Style { return termenv.String(s) },
		greener:      func(s string) termenv.Style { return termenv.String(s) },
		faintGreen:   func(s string) termenv.Style { return termenv.String(s) },
	}

	m := model{
		styles: styles,
		width:  80,
		height: 24,
	}

	// Test with various Korean WPM values
	testCases := []struct {
		name     string
		wpm      int
		accuracy float64
		deltaWpm float64
	}{
		{"High Korean WPM", 60, 98.5, 10.2},
		{"Average Korean WPM", 35, 92.0, -1.5},
		{"Low Korean WPM", 15, 85.3, -8.7},
		{"Perfect accuracy", 45, 100.0, 5.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			results := Results{
				wpm:      tc.wpm,
				accuracy: tc.accuracy,
				deltaWpm: tc.deltaWpm,
				rawWpm:   tc.wpm + 5,
				time:     1 * time.Minute,
				wordList: "Test Korean Words",
			}

			koreanResults := KoreanTimerBasedTestResults{
				results:       results,
				wpmEachSecond: []float64{float64(tc.wpm)},
			}

			m.state = koreanResults
			view := m.View()

			// Check that WPM is displayed correctly
			expectedWpm := fmt.Sprintf("wpm: %d", tc.wpm)
			if !strings.Contains(view, expectedWpm) {
				t.Errorf("Expected to find '%s' in Korean results display", expectedWpm)
			}

			// Check that accuracy is displayed correctly
			expectedAccuracy := fmt.Sprintf("accuracy: %.1f", tc.accuracy)
			if !strings.Contains(view, expectedAccuracy) {
				t.Errorf("Expected to find '%s' in Korean results display", expectedAccuracy)
			}
		})
	}
}