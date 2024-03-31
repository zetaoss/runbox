package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMD5(t *testing.T) {
	testcases := []struct {
		input string
		want  string
	}{
		{"hello", "5d41402abc4b2a76b9719d911017c592"},
		{"world", "7d793037a0760186574b0282f2f435e7"},
	}
	for _, tc := range testcases {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tc.want, MD5(tc.input))
		})
	}
}
