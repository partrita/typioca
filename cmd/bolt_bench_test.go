package cmd

import (
	"testing"

	"github.com/muesli/termenv"
)

func BenchmarkParagraphView(b *testing.B) {
	wordsToEnter := []rune("The quick brown fox jumps over the lazy dog. Programming is fun and rewarding. Terminal applications are powerful.")
	// Make it longer
	for i := 0; i < 5; i++ {
		wordsToEnter = append(wordsToEnter, wordsToEnter...)
	}

	base := &TestBase{
		wordsToEnter: wordsToEnter,
		inputBuffer:  []rune("The quick brown fox"),
		mistakes: mistakes{
			mistakesAt: make(map[int]bool),
		},
	}

	styles := Styles{
		correct: func(str string) termenv.Style {
			return termenv.String(str).Foreground(termenv.ANSI256Color(2))
		},
		toEnter: func(str string) termenv.Style {
			return termenv.String(str).Foreground(termenv.ANSI256Color(8))
		},
		mistakes: func(str string) termenv.Style {
			return termenv.String(str).Foreground(termenv.ANSI256Color(1)).Underline()
		},
		cursor: func(str string) termenv.Style {
			return termenv.String(str).Reverse()
		},
	}

	lineLimit := 40

	b.Run("NoCache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// Invalidate cache by changing rawInputCnt
			base.rawInputCnt++
			base.paragraphView(lineLimit, styles)
		}
	} )

	b.Run("WithCache", func(b *testing.B) {
		// Populate cache
		base.paragraphView(lineLimit, styles)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			base.paragraphView(lineLimit, styles)
		}
	})
}
