package cmd

import (
	"unicode/utf8"
)

func longestStringLen(strings []string) int {
	var longest int
	for _, elem := range strings {
		length := utf8.RuneCountInString(elem)
		if length > longest {
			longest = length
		}
	}

	return longest
}

func names(wordList []WordList) []string {
	var acc []string

	for _, elem := range wordList {
		acc = append(acc, elem.Name)
	}

	return acc
}

func averageStringLen(strings []string) int {
	var totalLen int = 0
	var cnt int = 0

	for _, str := range strings {
		currentLen := runeCountIgnoringAnsi(str)
		totalLen += currentLen
		cnt += 1
	}

	if cnt == 0 {
		cnt = 1
	}

	return totalLen / cnt
}

func averageLineLenFast(lines []string) int {
	linesLen := len(lines)
	linesToConsider := min(linesLen, 3)
	return averageStringLen(lines[:linesToConsider])
}

func averageLineLen(lines []string) int {
	linesLen := len(lines)
	if linesLen > 1 {
		lines = lines[:linesLen-1] //Drop last line, as it might skew up average length
	}

	return averageStringLen(lines)
}

// runeCountIgnoringAnsi counts the number of runes in a string while skipping
// ANSI escape sequences (CSI). This is more efficient than stripping ANSI codes
// and then counting runes because it avoids unnecessary string allocations.
func runeCountIgnoringAnsi(s string) int {
	count := 0
	lenS := len(s)
	for i := 0; i < lenS; {
		if s[i] == '\x1b' && i+1 < lenS && s[i+1] == '[' {
			i += 2
			// Skip CSI sequence: parameter bytes, intermediate bytes, and a final byte (0x40–0x7E)
			for i < lenS {
				b := s[i]
				i++
				if b >= 0x40 && b <= 0x7E {
					break
				}
			}
			continue
		}
		_, size := utf8.DecodeRuneInString(s[i:])
		count++
		i += size
	}
	return count
}

func floor(value int) int32 {
	return int32(max(0, value))
}

func dropLastString(strings []string) []string {
	le := len(strings)
	if le != 0 {
		return strings[:le-1]
	} else {
		return strings
	}
}

func dropLastRune(runes []rune) []rune {
	le := len(runes)
	if le != 0 {
		return runes[:le-1]
	} else {
		return runes
	}
}

func toKeysSlice(mp map[int]bool) []int {
	acc := []int{}
	for key := range mp {
		acc = append(acc, key)
	}
	return acc
}

func reverse(runes []rune) []rune {
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return runes
}
