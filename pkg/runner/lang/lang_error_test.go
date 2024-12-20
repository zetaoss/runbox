package lang

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zetaoss/runbox/pkg/runner/box"
	"github.com/zetaoss/runbox/pkg/testutil"
)

func TestRun_error(t *testing.T) {
	testCases := []struct {
		langInput Input
		wantError string
	}{
		{Input{Lang: "", Files: []box.File{}}, "no files"},
		{Input{Lang: "x", Files: []box.File{}}, "no files"},
		{Input{Lang: "go", Files: []box.File{}}, "no files"},
		{Input{Lang: "", Files: []box.File{{Body: `echo hello`}}}, "invalid language"},
		{Input{Lang: "x", Files: []box.File{{Body: `echo hello`}}}, "invalid language"},
	}
	for _, tc := range testCases {
		t.Run(testutil.Name(tc.langInput), func(t *testing.T) {
			output, err := lang1.Run(tc.langInput)
			require.Nil(t, output)
			require.EqualError(t, err, tc.wantError)
		})
	}
}

func TestRun_timeout(t *testing.T) {
	testCases := []struct {
		input Input
		want  *box.Result
	}{
		{
			Input{Lang: "bash", Files: []box.File{{Body: `echo hello; sleep 3`}}},
			&box.Result{
				Logs:     []box.Log{{Stream: 1, Log: "hello"}},
				Code:     0,
				CPU:      11040,
				MEM:      4648,
				Time:     2000,
				Timedout: true,
			},
		},
		{
			Input{Lang: "bash", Files: []box.File{{Body: `sleep 3; echo hello`}}},
			&box.Result{
				CPU:      10097,
				MEM:      788,
				Time:     2001,
				Timedout: true,
			},
		},
		{
			Input{Lang: "bash", Files: []box.File{{Body: `echo hello; echo world; sleep 3`}}},
			&box.Result{
				Logs: []box.Log{
					{Stream: 1, Log: "hello"},
					{Stream: 1, Log: "world"},
				},
				CPU:      9588,
				MEM:      796,
				Time:     2001,
				Timedout: true,
			},
		},
		{
			Input{Lang: "bash", Files: []box.File{{Body: `echo hello; sleep 3; echo world`}}},
			&box.Result{
				Logs:     []box.Log{{Stream: 1, Log: "hello"}},
				CPU:      9681,
				MEM:      804,
				Time:     2001,
				Timedout: true,
			},
		},
		{
			Input{Lang: "bash", Files: []box.File{{Body: `sleep 3; echo hello; echo world`}}},
			&box.Result{
				CPU:      9853,
				MEM:      800,
				Time:     2001,
				Timedout: true,
			},
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.input), func(t *testing.T) {
			got, err := lang1.Run(tc.input, map[string]int{"timeoutSeconds": 1})
			require.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}
