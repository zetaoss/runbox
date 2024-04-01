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
			&Docker{DataDir: "/data", HistoryFile: "/data/history.txt"},
			"",
		},
		{
			Config{DataDir: "/tmp/data"},
			&Docker{DataDir: "/tmp/data", HistoryFile: "/tmp/data/history.txt"},
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
		t.Run("", func(t *testing.T) {
			fakeErr = tc.fakeErr
			defer func() {
				fakeErr = NoError
			}()
			err := docker.saveHistory(Options{}, "hello")
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
			Options{Image: "alpine", Command: ""},
			&Result{Logs: []Log{}, ExitCode: 0},
			"",
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
		t.Run(tc.name, func(t *testing.T) {
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

func TestRun_fakeErr(t *testing.T) {
	testCases := []struct {
		fakeErr    Error
		wantResult *Result
		wantError  string
	}{
		{NoError, &Result{Logs: []Log{}, ExitCode: 1}, ""},
		{ErrOpenFile, &Result{Logs: []Log{}, ExitCode: 1}, ""},
		{ErrCollectLogs, nil, "collectLogs err: %!w(<nil>)"},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			fakeErr = tc.fakeErr
			defer func() {
				fakeErr = NoError
			}()
			result, err := docker.Run(Options{})
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
		containerID string
		content     string
		wantLogs    []Log
		wantError   string
	}{
		{"", "", []Log{}, ""},
		{"hello", "", []Log{}, ""},
		{"foo", "bar", []Log{}, ""},
		{"foo", "{}", []Log{{Log: "", Stream: ""}}, ""},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			if tc.containerID != "" {
				logFileDir := fmt.Sprintf("/var/lib/docker/containers/%s", tc.containerID)
				logFilePath := fmt.Sprintf("/var/lib/docker/containers/%s/%s-json.log", tc.containerID, tc.containerID)
				os.MkdirAll(logFileDir, 0644)
				defer func() {
					_ = os.RemoveAll(logFileDir)
				}()
				err := os.WriteFile(logFilePath, []byte(tc.content), 0644)
				require.NoError(t, err)
			}
			logs, err := docker.collectLogs(tc.containerID)
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
			require.Equal(t, tc.wantLogs, logs)
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
	os.MkdirAll(logFileDir, 0644)
	defer func() {
		_ = os.RemoveAll(logFileDir)
	}()
	err := os.WriteFile(logFilePath, []byte("{"), 0644)
	require.NoError(t, err)
	logs, err := docker.collectLogs("foo")
	require.NoError(t, err)
	require.Equal(t, []Log{}, logs)
}
