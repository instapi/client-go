package csvutil

import (
	"bufio"
	"io"
)

var _ io.Reader = (*LineLimitReader)(nil)

// LineLimitReader limits reading to a given number of lines.
type LineLimitReader struct {
	r    *bufio.Reader
	n    int
	line int
}

// NewLineLimitReader creates a new `LineLimitReader` instance.
func NewLineLimitReader(r io.Reader, n int) *LineLimitReader {
	return &LineLimitReader{r: bufio.NewReader(r), n: n}
}

func (l *LineLimitReader) Read(p []byte) (int, error) {
	if l.line == l.n {
		return 0, io.EOF
	}

	b, err := l.r.ReadBytes('\n')

	if err != nil {
		return 0, err
	}

	l.line++

	return copy(p, b), nil
}
