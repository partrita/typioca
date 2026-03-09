## 2026-03-03 - [Pre-compile regex for TUI rendering]
**Learning:** In terminal-based applications using frameworks like Bubble Tea, the `View` function and its helpers are called frequently (every tick or on every user input). Compiling regular expressions inside these functions using `regexp.MustCompile` is a significant performance bottleneck because regex compilation is an expensive operation.
**Action:** Always pre-compile regular expressions at the package level for functions in the TUI render loop or frequently executed paths.

## 2026-03-04 - [Batch styling in TUI render loop]
**Learning:** In TUI applications, styling text character-by-character (e.g., applying a foreground color to each rune individually) is extremely inefficient. It results in $O(N)$ sets of ANSI escape sequences, significantly bloating the output string size and increasing terminal parsing overhead.
**Action:** Always batch style text wherever possible. Converting a slice of runes to a string and styling the entire string at once can improve performance by orders of magnitude (~98% in this case) and reduce the ANSI payload.

## 2026-03-05 - [Optimize TUI rendering by reducing ANSI processing and using O(1) heuristics]
**Learning:** In performance-critical TUI render loops, stripping ANSI codes via regex to calculate string length is expensive due to frequent allocations and regex overhead. Additionally, calculating metrics like average line length over a large text block in every `View` call is $O(N)$ where $N$ is the number of lines.
**Action:** Implement manual rune counting that skips ANSI CSI sequences (terminating on 0x40-0x7E) to avoid allocations. Use $O(1)$ heuristics (like averaging only the first 3 lines) for layout centering in the `View` loop.

## 2026-03-06 - [Avoid unnecessary type conversions in Go 1.21+ logic]
**Learning:** Using `math.Max` or `math.Min` for integer clipping/clamping in Go requires casting to and from `float64`, which is clunky and slightly less efficient than using the built-in `max()` and `min()` functions introduced in Go 1.21. For generating repetitive terminal sequences like newlines, `strings.Repeat` is also more efficient than manual loops.
**Action:** Use built-in `max()`/`min()` for integer-based logic in projects using Go 1.21+. Use `strings.Repeat` for simple repetitive string generation.

## 2026-03-09 - [Pre-allocate slices for map-to-slice conversions]
**Learning:** Converting a map to a slice (e.g., extracting keys) without pre-allocating the slice results in multiple memory reallocations as the slice grows. In performance-critical paths like TUI rendering, these reallocations add up and increase GC pressure.
**Action:** Always pre-allocate slices with `make([]T, 0, len(m))` when converting maps to slices to ensure only a single allocation occurs.
