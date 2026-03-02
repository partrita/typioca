## 2025-03-02 - Regex Compilation Overhead
**Learning:** Calling `regexp.MustCompile` inside a function that is part of a rendering loop or frequent calculation (like `dropAnsiCodes` in a TUI app) introduces significant overhead (~70% in this case).
**Action:** Always pre-compile regular expressions at the package level for static patterns.

## 2025-03-02 - Efficient Rune Counting
**Learning:** `len([]rune(s))` allocates a new slice and copies data, while `utf8.RuneCountInString(s)` counts in place.
**Action:** Use `utf8.RuneCountInString` for character counting in Go strings, especially in performance-sensitive paths.
