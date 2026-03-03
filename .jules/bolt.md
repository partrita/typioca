## 2026-03-03 - [Pre-compile regex for TUI rendering]
**Learning:** In terminal-based applications using frameworks like Bubble Tea, the `View` function and its helpers are called frequently (every tick or on every user input). Compiling regular expressions inside these functions using `regexp.MustCompile` is a significant performance bottleneck because regex compilation is an expensive operation.
**Action:** Always pre-compile regular expressions at the package level for functions in the TUI render loop or frequently executed paths.
