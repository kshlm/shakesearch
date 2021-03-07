## Current changes

### Backend

- Smartcase search - If query is a lowercase word, perform a case-insensitive search, if not to an exact string search
- Better match context - The full block of text containing a match is returned instead of a simple fixed length context.
- De-duplicated blocks - Blocks containing multiple matches are de-duplicated.
- Highlight match - All matches in a block are highlighted.

### Frontend

- Use `<pre>` blocks to display results - this makes it easier for the users to follow results.


## Possible changes (in order or priority)

### Results with metadata

Instead of sending a list of pre-highlighted strings, we can send a structured document with metadata attached to each result, and let the frontend handle display. For eg.,

```json
[
  {
    "blockText": "<textBlock>",
    "matches": [100, 122, 400]
  },
  {}
]
```

### More metadata with result

More metadata can be added to each result to make easier to identify it. This could include

- Line Number in the full text file
- The work the match is contained within
- Location of in the work (Act, Scene) etc,

Doing this will need pre-processing of the source text.

This metadata can be in the frontend to perform filtering of results based on work, act, scene etc.


### Multiline matches, Fuzzy matches, Regex search

The current search only does very simple matching and works best with single words. Adding support for regex searches, multiline matches and partial matches will be useful for users. Users searching Shakespeare's works are likely to search for long(er) phrases like dialogues in a play or verses in a poem.

### Search suggestions/Autocomplete/Spellcheck

Search suggestion, auto-completion, spellcheck features would improve usability for users not very familiar with Shakespeare's works. These could initially be implemented as a server side feature on the search API. But this would make the biggest impact if a streaming-api could be implemented which constantly provides these suggestions as the user types in the frontend.
