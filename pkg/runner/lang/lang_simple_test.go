package lang

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zetaoss/runbox/pkg/runner/box"
	"github.com/zetaoss/runbox/pkg/testutil"
	"k8s.io/utils/ptr"
)

var lang1 *Lang

func init() {
	d := testutil.NewDocker()
	lang1 = New(box.New(d))
}

func equalResult(t *testing.T, want, got *box.Result) {
	t.Helper()

	assert.Greater(t, got.CPU, want.CPU/100, "want.CPU", want.CPU)
	assert.Greater(t, got.MEM, want.MEM/1000, "want.MEM", want.MEM)
	assert.Less(t, got.CPU, want.CPU*100, "want.CPU", want.CPU)
	assert.Less(t, got.MEM, want.MEM*1000, "want.MEM", want.MEM)
	want.CPU = got.CPU
	want.MEM = got.MEM

	assert.Greater(t, got.Time, want.Time/100, "want.Time", want.Time)
	assert.Less(t, got.Time, want.Time*100, "want.Time", want.Time)
	want.Time = got.Time

	assert.Equal(t, want, got)
}

func TestToLangOpts(t *testing.T) {
	testcases := []struct {
		input     Input
		want      *LangOpts
		wantError string
	}{
		{
			Input{
				Lang: "bash",
				Files: []box.File{
					{Name: "greet.txt", Body: "hello"},
					{Body: "cat greet.txt"},
				},
				Main: 1,
			},
			&LangOpts{
				Input:          Input{Lang: "bash", Files: []box.File{{Name: "greet.txt", Body: "hello"}, {Body: "cat greet.txt"}}, Main: 1},
				Command:        "/bin/bash runbox.sh",
				FileName:       "runbox",
				FileExt:        "sh",
				Shell:          "bash",
				TimeoutSeconds: 10,
				WorkingDir:     "/home/user01",
			},
			"",
		},
	}
	for _, tc := range testcases {
		t.Run("", func(t *testing.T) {
			got, err := toLangOpts(tc.input)
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
			require.Equal(t, tc.want, got)
		})
	}
}

func TestToBoxOpts(t *testing.T) {
	testcases := []struct {
		langOpts LangOpts
		want     box.Opts
	}{
		{
			LangOpts{
				Input:          Input{Lang: "bash", Files: []box.File{{Name: "greet.txt", Body: "hello"}, {Body: "cat greet.txt"}}, Main: 1},
				Command:        "/bin/bash runbox.sh",
				FileName:       "runbox",
				FileExt:        "sh",
				Shell:          "bash",
				TimeoutSeconds: 10,
				WorkingDir:     "/home/user01",
			},
			box.Opts{
				CollectStats:  ptr.To(true),
				CollectImages: true,
				Command:       "/bin/bash runbox.sh",
				Env:           nil,
				Files: []box.File{
					{Name: "/home/user01/greet.txt", Body: "hello"},
					{Name: "/home/user01/runbox.sh", Body: "cat greet.txt"},
				},
				Image:      "ghcr.io/zetaoss/runcontainers/bash",
				Shell:      "bash",
				Timeout:    10000,
				WorkingDir: "/home/user01",
			},
		},
	}
	for _, tc := range testcases {
		t.Run("", func(t *testing.T) {
			got := toBoxOpts(tc.langOpts)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestRun_simple(t *testing.T) {
	testcases := []struct {
		input     Input
		want      *box.Result
		wantError string
	}{
		{
			Input{
				Lang: "bash",
				Files: []box.File{
					{Name: "greet.txt", Body: "hello"},
					{Body: "cat greet.txt"},
				},
				Main: 1,
			},
			&box.Result{
				Logs:     []box.Log{{Stream: 1, Log: "hello"}},
				CPU:      9183,
				MEM:      676,
				Time:     196,
				Timedout: false,
			},
			"",
		},
	}
	for _, tc := range testcases {
		t.Run("", func(t *testing.T) {
			got, err := lang1.Run(tc.input)
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
			equalResult(t, tc.want, got)
		})
	}
}
