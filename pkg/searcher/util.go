package searcher

import (
	"regexp"
)

const (
	hlSep = "**"
)

var (
	lowercaseWord = regexp.MustCompile(`^[\p{Ll}]+$`)
)

// highlightMatch highlights a match in a text block
// Takes a string, a start index, and a length of section to be
// highlighted.
// Returns a string with the section higlighted.
// To highlight a section is wrapped with "**"
func highlightSection(text string, start, length int) string {
	if start+length < len(text) {
		return text[:start]+hlSep+text[start:start+length]+hlSep+text[start+length:]
	} else {
		return text[:start]+hlSep+text[start:]+hlSep
	}
}

func isLowercaseWord(s string) bool {
	return lowercaseWord.Match([]byte(s))
}
