package docker

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var docker *Docker

func init() {
	var err error
	docker, err = New()
	if err != nil {
		panic(err)
	}
}

func TestNewWithConfig(t *testing.T) {
	testCases := []struct {
		config     Config
		wantDocker *Docker
		wantError  string
	}{
		{
			Config{},
			&Docker{DataDir: "/tmp/runbox", HistoryFile: "/tmp/runbox/history.txt"},
			"",
		},
		{
			Config{DataDir: "/tmp/runbox"},
			&Docker{DataDir: "/tmp/runbox", HistoryFile: "/tmp/runbox/history.txt"},
			"",
		},
		{
			Config{DataDir: "/etc/os-release"},
			nil,
			"MkdirAll err",
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			d, err := NewWithConfig(tc.config)
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.wantError)
			}
			require.Equal(t, tc.wantDocker, d)
		})
	}
}

func TestSaveHistory(t *testing.T) {
	testCases := []struct {
		fakeErr   Error
		wantError string
	}{
		{NoError, ""},
		{ErrOpenFile, "OpenFile err: %!w(<nil>)"},
		{ErrWriteString, "WriteString err: %!w(<nil>)"},
		{ErrClose, "close err: %!w(<nil>)"},
	}
	for _, tc := range testCases {
		t.Run("_"+string(tc.fakeErr), func(t *testing.T) {
			fakeErr = tc.fakeErr
			defer func() {
				fakeErr = NoError
			}()
			err := docker.saveHistory("hello")
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
		})
	}
}

func TestRun(t *testing.T) {
	testCases := []struct {
		name       string
		options    Options
		wantResult *Result
		wantError  string
	}{
		{
			"",
			Options{},
			nil,
			"no image",
		},
		{
			"",
			Options{Command: "echo hello"},
			nil,
			"no image",
		},
		{
			"",
			Options{Image: "alpine", Command: ""},
			nil,
			"no command",
		},
		{
			"etc",
			Options{Image: "alpine", Command: "/etc"},
			&Result{Logs: []Log{}, ExitCode: 126},
			"",
		},
		// echo
		{
			"echo",
			Options{Image: "alpine", Command: "echo foo"},
			&Result{Logs: []Log{{Log: "foo\n", Stream: "stdout"}}, ExitCode: 0},
			"",
		},
		{
			"echo",
			Options{Image: "alpine", Command: "echo foo bar"},
			&Result{Logs: []Log{{Log: "foo bar\n", Stream: "stdout"}}, ExitCode: 0},
			"",
		},
		{
			"echo",
			Options{Image: "alpine", Command: "echo foo; echo bar"},
			&Result{Logs: []Log{{Log: "foo\n", Stream: "stdout"}}, ExitCode: 0},
			"",
		},
		{
			"echo",
			Options{Image: "alpine", Command: "echo foo bar", Binds: []string{"foo", "bar"}},
			&Result{Logs: []Log{{Log: "foo bar\n", Stream: "stdout"}}, ExitCode: 0},
			"",
		},
		{
			"sleep",
			Options{Image: "alpine", Command: "sleep 1"},
			&Result{Logs: []Log{}, ExitCode: 0},
			"",
		},
		{
			"sleep",
			Options{Image: "alpine", Command: "sleep 1; echo foo"},
			&Result{Logs: []Log{}, ExitCode: 0},
			"",
		},
	}
	for _, tc := range testCases {
		t.Run("_"+tc.name, func(t *testing.T) {
			result, err := docker.Run(tc.options)
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
			require.Equal(t, tc.wantResult, result)
		})
	}
}

func TestRun_OutputLimitReached(t *testing.T) {
	testCases := []struct {
		name           string
		options        Options
		wantResult     *Result
		wantLogLengths []int
		wantError      string
	}{
		{
			"python",
			Options{Image: "jmnote/runbox:python", Command: `python -c "print(10*'HelloWorld')"`, OutputLimit: 0},
			&Result{ExitCode: 0},
			[]int{101},
			"",
		},
		{
			"python",
			Options{Image: "jmnote/runbox:python", Command: `python -c "print(100*'HelloWorld')"`, OutputLimit: 0},
			&Result{ExitCode: 0},
			[]int{1001},
			"",
		},
		{
			"python",
			Options{Image: "jmnote/runbox:python", Command: `python -c "print(10*'HelloWorld')"`, OutputLimit: 500},
			&Result{ExitCode: 0},
			[]int{101},
			"",
		},
		{
			"python",
			Options{Image: "jmnote/runbox:python", Command: `python -c "print(100*'HelloWorld')"`, OutputLimit: 500},
			&Result{ExitCode: 0, OutputLimitReached: true},
			[]int{500},
			"",
		},
	}
	for _, tc := range testCases {
		t.Run("_"+tc.name, func(t *testing.T) {
			result, err := docker.Run(tc.options)
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
			logLengths := []int{}
			for _, l := range result.Logs {
				logLengths = append(logLengths, len(l.Log))
			}
			// ignore fields
			result.Logs = nil
			require.Equal(t, tc.wantResult, result)
			require.Equal(t, tc.wantLogLengths, logLengths)
		})
	}
}

func TestRun_fakeErr(t *testing.T) {
	testCases := []struct {
		fakeErr    Error
		wantResult *Result
		wantError  string
	}{
		{NoError, &Result{Logs: []Log{{Log: "hello\n", Stream: "stdout"}}, ExitCode: 0}, ""},
		{ErrOpenFile, &Result{Logs: []Log{{Log: "hello\n", Stream: "stdout"}}, ExitCode: 0}, ""},
		{ErrCollectLogs, nil, "collectLogs err: %!w(<nil>)"},
	}
	for _, tc := range testCases {
		t.Run("_"+string(tc.fakeErr), func(t *testing.T) {
			fakeErr = tc.fakeErr
			defer func() {
				fakeErr = NoError
			}()
			result, err := docker.Run(Options{Image: "alpine", Command: "echo hello"})
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
			require.Equal(t, tc.wantResult, result)
		})
	}
}

func TestCollectLogs(t *testing.T) {
	testCases := []struct {
		containerID      string
		content          string
		outputLimit      int
		wantLogs         []Log
		wantLimitReached bool
		wantError        string
	}{
		{"", "", 0, nil, false, "open /var/lib/docker/containers//-json.log: no such file or directory"},
		{"hello", "", 0, []Log{}, false, ""},
		{"foo", "bar", 0, []Log{}, false, ""},
		{"foo", "{}", 0, []Log{{Log: "", Stream: ""}}, false, ""},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			if tc.containerID != "" {
				logFileDir := fmt.Sprintf("/var/lib/docker/containers/%s", tc.containerID)
				logFilePath := fmt.Sprintf("/var/lib/docker/containers/%s/%s-json.log", tc.containerID, tc.containerID)
				_ = os.MkdirAll(logFileDir, 0644)
				defer func() {
					_ = os.RemoveAll(logFileDir)
				}()
				err := os.WriteFile(logFilePath, []byte(tc.content), 0644)
				require.NoError(t, err)
			}
			logs, limitReached, err := docker.collectLogs(tc.containerID, tc.outputLimit)
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
			require.Equal(t, tc.wantLogs, logs)
			require.Equal(t, tc.wantLimitReached, limitReached)
		})
	}
}

func TestCollectLogs_fakeErr(t *testing.T) {
	fakeErr = ErrScanner
	defer func() {
		fakeErr = NoError
	}()
	containerID := "foo"
	logFileDir := fmt.Sprintf("/var/lib/docker/containers/%s", containerID)
	logFilePath := fmt.Sprintf("/var/lib/docker/containers/%s/%s-json.log", containerID, containerID)
	_ = os.MkdirAll(logFileDir, 0644)
	defer func() {
		_ = os.RemoveAll(logFileDir)
	}()
	err := os.WriteFile(logFilePath, []byte("{"), 0644)
	require.NoError(t, err)
	logs, limitReached, err := docker.collectLogs("foo", 0)
	require.NoError(t, err)
	require.False(t, limitReached)
	require.Equal(t, []Log{}, logs)
}
