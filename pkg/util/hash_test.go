package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHash(t *testing.T) {
	testcases := []struct {
		input    string
		length   int
		wantHash string
	}{
		{"", 5, "lGqEi"},
		{"hello", 100, "iAfgxGtocqxgrfkoBwEuFxEwhepAdFiaiAfgxGtocqxgrfkoBwEuFxEwhepAdFiaiAfgxGtocqxgrfkoBwEuFxEwhepAdFiaiAfg"},
	}

	for _, tc := range testcases {
		t.Run(tc.input, func(t *testing.T) {
			got := Hash(tc.input, tc.length)
			require.Equal(t, tc.wantHash, got)
		})
	}
}

func TestNewHash(t *testing.T) {
	for i := 0; i < 100; i++ {
		h := NewHash(i)
		require.Equal(t, i, len(h))
	}
}
