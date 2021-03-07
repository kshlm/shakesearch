package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handleSearch(searcher))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("Listening on port %s...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

type Searcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}
		results := searcher.Search(query[0])
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
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
	results := []string{}
	qLen := len(query)
	for _, idx := range idxs {
		block, blockIdx := s.GetBlock(idx, qLen)
		results = append(results, highlightSection(block, blockIdx, qLen))
	}
	return results
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

// highlightMatch highlights a match in a text block
// Takes a block, a start index, and a length of section to be
// highlighted.
// Returns a string with the section higlighted.
// To highlight a section is wrapped with "**"
func highlightSection(block string, start, length int) string {
	return block[0:start]+"**"+block[start:start+length]+"**"+block[start+length:]
}
