package csvutil

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLineLimitReader(t *testing.T) {
	const testdata = "A\nB\nC\nD\nE\nF"

	sr := strings.NewReader(testdata)
	lr := NewLineLimitReader(sr, 3)

	b, err := io.ReadAll(lr)

	require.NoError(t, err)
	require.Equal(t, []byte("A\nB\nC\n"), b)
}
