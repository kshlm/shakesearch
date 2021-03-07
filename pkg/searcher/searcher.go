package searcher

import (
	"fmt"
	"index/suffixarray"
	"io/ioutil"
)

type Searcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
}

func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.CompleteWorks = string(dat)
	s.SuffixArray = suffixarray.New(dat)
	return nil
}

func (s *Searcher) Search(query string) []string {
	idxs := s.SuffixArray.Lookup([]byte(query), -1)
	results := NewResults(s.CompleteWorks, len(query))
	for _, idx := range idxs {
		results.AddMatch(idx)
	}
	return results.ToList()
}

// GetBlock returns the block of text the given index is within, and the
// relative index within the block.
// A block of text is defined as a set of lines bounded by '\r\n' lines.
func (s *Searcher) GetBlock(index, qlen int) (string, int) {

	var start, end int
	const emptyLine = "\r\n\r\n"

	for start = index; (start - 4) >= 0; start-- {
		if s.CompleteWorks[start-4:start] == emptyLine {
			break
		}
	}
	for end = index+qlen; (end+5) < len(s.CompleteWorks) ; end++ {
		if s.CompleteWorks[end+1:end+5] == emptyLine {
			break
		}
	}

	return s.CompleteWorks[start:end], (index - start)
}
