## 2026-03-03 - [Pre-compile regex for TUI rendering]
**Learning:** In terminal-based applications using frameworks like Bubble Tea, the `View` function and its helpers are called frequently (every tick or on every user input). Compiling regular expressions inside these functions using `regexp.MustCompile` is a significant performance bottleneck because regex compilation is an expensive operation.
**Action:** Always pre-compile regular expressions at the package level for functions in the TUI render loop or frequently executed paths.

## 2026-03-04 - [Batch styling in TUI render loop]
**Learning:** In TUI applications, styling text character-by-character (e.g., applying a foreground color to each rune individually) is extremely inefficient. It results in $O(N)$ sets of ANSI escape sequences, significantly bloating the output string size and increasing terminal parsing overhead.
**Action:** Always batch style text wherever possible. Converting a slice of runes to a string and styling the entire string at once can improve performance by orders of magnitude (~98% in this case) and reduce the ANSI payload.
