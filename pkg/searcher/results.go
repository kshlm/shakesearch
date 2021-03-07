package searcher

import (
	"sort"
)


const (
	hlSep = "**"
	emptyLine = "\r\n\r\n"
)

// Results holds a list of text blocks with the found matches
type Results struct {
	works string // The complete works of Shakespeare
	mLen int // Length of the match/query
	blocks map[int]*block // Matched blocks
}

// block holds a block of matched text
type block struct {
	start, end int // Start and end of block
	mLen int // Length of match
	text string // The text contained within the block
	matches []int // Start indices of all matches in block
}

// NewResults returns an empty Results struct
func NewResults(works string, qLen int) *Results {
	return &Results {
		works: works,
		mLen: qLen,
		blocks: make(map[int]*block),
	}
}

// AddMatch adds a found match to the Results
func (r *Results) AddMatch(idx int) {
	// If match is contained within a previously found block, add match to it
	if block := r.findBlockWithIndex(idx); block != nil {
		block.addMatch(idx)
		return
	}
	// Or create a new block for the match
	block := r.newBlock(idx)
	r.blocks[block.start] = block
	return
}

// findBlockWithIndex finds an existing block containing the given index
func (r *Results) findBlockWithIndex(idx int) *block {
	for _, b := range r.blocks {
		if b.start <= idx && idx < b.end {
			return b
		}
	}
	return nil
}

// newBlock creates a block for the given index
func (r *Results) newBlock(idx int) *block {
	var start, end int

	for start = idx; (start - 4) >= 0; start-- {
		if r.works[start-4:start] == emptyLine {
			break
		}
	}
	for end = idx+r.mLen; (end+5) < len(r.works) ; end++ {
		if r.works[end+1:end+5] == emptyLine {
			break
		}
	}

	block := &block{
		start: start,
		end: end,
		mLen: r.mLen,
		text: r.works[start:end],
	}
	block.addMatch(idx)

	return block
}

// Returns a list of blocks with the matches highlighted
func (r *Results) ToList() []string {
	var hlBlocks []string
	for _, block := range r.blocks {
		hlBlocks = append(hlBlocks, block.higlightMatches())
	}
	return hlBlocks
}

// addMatch adds a match to a block
func (b *block) addMatch(idx int) {
	relIdx := idx - b.start
	b.matches = append(b.matches, relIdx)
}

// higlightMatches highlights all the matches in a block and returns a highlighted string
func (b *block) higlightMatches() string {
	text := b.text

	sort.Ints(b.matches)
	for n, mi:= range b.matches {
		// adjust match index after each match is highlighted
		ami:= (n * 2 * len(hlSep)) + mi
		text = highlightSection(text, ami, b.mLen)
	}
	return text
}

// highlightMatch highlights a match in a text block
// Takes a string, a start index, and a length of section to be
// highlighted.
// Returns a string with the section higlighted.
// To highlight a section is wrapped with "**"
func highlightSection(text string, start, length int) string {
	return text[0:start]+hlSep+text[start:start+length]+hlSep+text[start+length:]
}
