package searcher

// highlightMatch highlights a match in a text block
// Takes a string, a start index, and a length of section to be
// highlighted.
// Returns a string with the section higlighted.
// To highlight a section is wrapped with "**"
func highlightSection(text string, start, length int) string {
	return text[0:start]+hlSep+text[start:start+length]+hlSep+text[start+length:]
}
