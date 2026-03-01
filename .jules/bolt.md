## 2026-03-01 - [Pre-compiled Regex and Efficient Rune Counting]
**Learning:** In TUI applications built with Bubble Tea/Lipgloss, rendering logic (like `View()` and its helpers) is called frequently. Repeatedly compiling regex with `regexp.MustCompile` inside these functions creates a significant performance bottleneck. Additionally, using `len([]rune(s))` for character counting causes unnecessary allocations.
**Action:** Always pre-compile regular expressions as package-level variables. Use `utf8.RuneCountInString(s)` for efficient, allocation-free rune counting in Go.
