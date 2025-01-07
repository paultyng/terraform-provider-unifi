package markdown

// indentation tracks indentation data for a mardown writer.
type indentation struct {
	// Indentation is comprised of multiple sections of indentations.
	// indents tracks the combined full indentation,
	// and lengths tracks the length of each appended section.
	indents []byte
	lengths []int

	// index at which trailing spaces start in indents.
	trailSpaceIdx int
}

// Indent reports the fixed text prefix pushed so far.
//
// Invariant: This does not end with whitespace.
func (id *indentation) Indent() []byte {
	return id.indents[:id.trailSpaceIdx]
}

// Whitespace reports the trailing whitespace of the indentation pushed so far.
func (id *indentation) Whitespace() []byte {
	return id.indents[id.trailSpaceIdx:]
}

// Push adds a block of text to the indentation stack.
//
// Indent and Whitespace will report this in consescutive calls.
func (id *indentation) Push(bs []byte) {
	id.indents = append(id.indents, bs...)
	id.lengths = append(id.lengths, len(bs))
	id.trailSpaceIdx = trailingSpaceIdx(id.indents)
}

// Pop removes the last pushed block of text from the stack.
func (id *indentation) Pop() {
	count := len(id.lengths)
	if count == 0 {
		panic("bug: indentation.Pop called for empty indentation")
	}
	lastLen := id.lengths[count-1]

	id.lengths = id.lengths[:count-1]
	id.indents = id.indents[:len(id.indents)-lastLen]
	id.trailSpaceIdx = trailingSpaceIdx(id.indents)
}

// trailingSpaceIdx returns the index at which trailing space
// starts in the given byte slice.
//
// Returns 0 if the slice is entirely whitespace,
// and len(bs) if the slice is entirely non-whitespace.
func trailingSpaceIdx(bs []byte) int {
	for idx := len(bs); idx > 0; idx-- {
		if bs[idx-1] != ' ' {
			return idx
		}
	}
	return 0
}
