package cmd

import (
	"testing"
	"github.com/muesli/termenv"
)

func BenchmarkStyleAllRunes(b *testing.B) {
	runes := []rune("The quick brown fox jumps over the lazy dog. 1234567890. ")
	// Make it 1000 runes
	longRunes := make([]rune, 0, 1000)
	for len(longRunes) < 1000 {
		longRunes = append(longRunes, runes...)
	}
	longRunes = longRunes[:1000]

	styleFn := func(str string) termenv.Style {
		return termenv.String(str).Foreground(termenv.ANSI256Color(2))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		styleAllRunes(longRunes, styleFn)
	}
}
