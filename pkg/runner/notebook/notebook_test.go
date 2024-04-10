package notebook

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		input      Input
		wantOutput *Output
		wantError  string
	}{
		{Input{}, nil, "ErrInvalidLanguage"},
		{Input{Lang: "xxx"}, nil, "ErrInvalidLanguage"},
		{Input{Lang: "r"}, &Output{}, ""},
		{Input{Lang: "python"}, &Output{}, ""},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			output, err := Run(tc.input)
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
			require.Equal(t, tc.wantOutput, output)
		})
	}
}
