package googling

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zetaoss/runbox/pkg/runner/box"
	"github.com/zetaoss/runbox/pkg/testutil"
)

var googling1 *Googling

func init() {
	d := testutil.NewDocker()
	googling1 = New(box.New(d))
}

func TestRun(t *testing.T) {
	testCases := []struct {
		q    string
		want int
	}{
		{"hello", 4740000000},
		{"openai", 317000000},
		{"zetawiki", 137000},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.q), func(t *testing.T) {
			got, err := googling1.Run(tc.q)
			require.NoError(t, err)
			require.Greater(t, got, tc.want/10)
			require.Less(t, got, tc.want*10)
		})
	}
}
