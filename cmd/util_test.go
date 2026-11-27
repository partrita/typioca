package cmd

import (
	"testing"

	"github.com/muesli/termenv"
)

func plainStyle() StringStyle {
	return func(str string) termenv.Style {
		return termenv.String(str)
	}
}

func TestAverageLineLenDropsTheLastLine(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  int
	}{
		{
			name:  "no lines",
			lines: nil,
			want:  0,
		},
		{
			name:  "single line is used as is",
			lines: []string{"abcdefghij"},
			want:  10,
		},
		{
			name:  "short last line is ignored",
			lines: []string{"abcdefghij", "abcdefghij", "abc"},
			want:  10,
		},
		{
			name:  "only the first three full lines are considered",
			lines: []string{"aaaaaaaaaa", "bbbbbbbbbb", "cccccccccc", "dddddddddd", "ee"},
			want:  10,
		},
		{
			name:  "mixed lengths average over the full lines",
			lines: []string{"aaaaaaaa", "bbbb", "cccccccccc", "dd"},
			want:  7,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := averageLineLen(test.lines)
			if got != test.want {
				t.Fatalf("averageLineLen(%q) = %d, want %d", test.lines, got, test.want)
			}
		})
	}
}

func TestAverageLineLenIgnoresAnsiCodes(t *testing.T) {
	lines := []string{
		"\x1b[32mabcdefghij\x1b[0m",
		"\x1b[32mabcdefghij\x1b[0m",
		"\x1b[32mabc\x1b[0m",
	}

	if got := averageLineLen(lines); got != 10 {
		t.Fatalf("averageLineLen(%q) = %d, want 10", lines, got)
	}
}

func TestFormatDeltaWpm(t *testing.T) {
	styles := Styles{greener: plainStyle()}

	tests := []struct {
		name  string
		delta float64
		want  string
	}{
		{name: "positive delta keeps its sign", delta: 12.3456, want: "Δavg: +12.35%"},
		{name: "negative delta keeps its sign", delta: -8.5, want: "Δavg: -8.50%"},
		{name: "large delta is clamped", delta: 250, want: "Δavg: +100.00%"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := formatDeltaWpm(test.delta, styles); got != test.want {
				t.Fatalf("formatDeltaWpm(%v) = %q, want %q", test.delta, got, test.want)
			}
		})
	}
}
