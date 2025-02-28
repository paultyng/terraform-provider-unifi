package markdown

import (
	"io"
)

// lineIndentWriter wraps io.Writer and adds given indent everytime new line is created .
type lineIndentWriter struct {
	io.Writer

	id                    indentation
	firstWriteExtraIndent []byte

	previousCharWasNewLine bool
}

func wrapWithLineIndentWriter(w io.Writer) *lineIndentWriter {
	return &lineIndentWriter{Writer: w, previousCharWasNewLine: true}
}

func (l *lineIndentWriter) PushIndent(indent []byte) {
	l.id.Push(indent)
}

func (l *lineIndentWriter) PopIndent() {
	l.id.Pop()
}

func (l *lineIndentWriter) AddIndentOnFirstWrite(add []byte) {
	l.firstWriteExtraIndent = append(l.firstWriteExtraIndent, add...)
}

func (l *lineIndentWriter) DelIndentOnFirstWrite(del []byte) {
	l.firstWriteExtraIndent = l.firstWriteExtraIndent[:len(l.firstWriteExtraIndent)-len(del)]
}

func (l *lineIndentWriter) WasIndentOnFirstWriteWritten() bool {
	return len(l.firstWriteExtraIndent) == 0
}

func (l *lineIndentWriter) Write(b []byte) (n int, _ error) {
	if len(b) == 0 {
		return 0, nil
	}

	writtenFromB := 0
	for i, c := range b {
		if l.previousCharWasNewLine {
			ns, err := l.Writer.Write(l.id.Indent())
			n += ns
			if err != nil {
				return n, err
			}
		}

		if c == newLineChar[0] {
			if !l.WasIndentOnFirstWriteWritten() {
				ns, err := l.Writer.Write(l.firstWriteExtraIndent)
				n += ns
				if err != nil {
					return n, err
				}
				l.firstWriteExtraIndent = nil
			}

			ns, err := l.Writer.Write(b[writtenFromB : i+1])
			n += ns
			writtenFromB += ns
			if err != nil {
				return n, err
			}
			l.previousCharWasNewLine = true
			continue
		}

		// Not a newline, make a space if indent was created.
		if l.previousCharWasNewLine {
			ws := l.id.Whitespace()
			if len(ws) > 0 {
				ns, err := l.Writer.Write(ws)
				n += ns
				if err != nil {
					return n, err
				}
			}
		}
		l.previousCharWasNewLine = false
	}

	if writtenFromB >= len(b) {
		return n, nil
	}

	if !l.WasIndentOnFirstWriteWritten() {
		ns, err := l.Writer.Write(l.firstWriteExtraIndent)
		n += ns
		if err != nil {
			return n, err
		}
		l.firstWriteExtraIndent = nil
	}

	ns, err := l.Writer.Write(b[writtenFromB:])
	n += ns
	return n, err
}
