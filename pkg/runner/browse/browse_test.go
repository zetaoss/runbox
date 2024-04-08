package browse

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		url        string
		wantOutput *Output
		wantError  string
	}{
		{"", nil, "ErrNoURL"},
		{"xxx", &Output{}, ""},
		{"http://example.com", &Output{Text: "<title>Example Domain</title>"}, ""},
		{"https://example.com", &Output{Text: "<title>Example Domain</title>"}, ""},
		{"https://google.com", &Output{Text: "<title>Google</title>"}, ""},
		{"https://zetawiki.com", &Output{Text: "<title>제타위키</title>"}, ""},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d_%v", i, tc.url), func(t *testing.T) {
			output, err := Run(tc.url)
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
			if tc.wantOutput == nil {
				require.Nil(t, output)
			} else {
				wantText := tc.wantOutput.Text
				require.Contains(t, output.Text, wantText)
				// ignore fields
				tc.wantOutput.Text = ""
				output.Text = ""
				require.Equal(t, tc.wantOutput, output)
			}
		})
	}
}
