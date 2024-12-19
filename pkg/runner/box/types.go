package box

import (
	"strings"
)

type Opts struct {
	CollectStats          *bool
	CollectImages         bool
	CollectImagesCount    int
	Command               string
	Env                   []string
	Files                 []File
	Image                 string
	PullImageIfNotPresent *bool
	Shell                 string
	Timeout               int
	User                  string
	WorkingDir            string
}

type Result struct {
	Logs     []Log    `json:"logs,omitempty"`
	Code     int      `json:"code,omitempty"`
	CPU      int      `json:"cpu,omitempty"`
	MEM      int      `json:"mem,omitempty"`
	Time     int      `json:"time,omitempty"`
	Timedout bool     `json:"timedout,omitempty"`
	Images   []string `json:"images,omitempty"`
}

type File struct {
	Name string `json:"name"`
	Body string `json:"body"`
}

type Log struct {
	Stream int
	Log    string
}

type logWriter struct {
	stream int
	logs   *[]Log
	buffer string
}

func (w *logWriter) Write(p []byte) (n int, err error) {
	w.buffer += string(p)
	lines := strings.Split(w.buffer, "\n")

	w.buffer = lines[len(lines)-1]
	lines = lines[:len(lines)-1]

	for _, line := range lines {
		log := Log{
			Stream: w.stream,
			Log:    strings.TrimRight(line, "\n"),
		}
		*w.logs = append(*w.logs, log)
	}
	return len(p), nil
}

func (w *logWriter) Close() {
	if w.buffer != "" {
		log := Log{
			Stream: w.stream,
			Log:    w.buffer,
		}
		*w.logs = append(*w.logs, log)
		w.buffer = ""
	}
}

func (r *Result) StreamStrings() (string, string) {
	var stdout, stderr string
	for _, l := range r.Logs {
		if l.Stream == 1 {
			stdout += l.Log + "\n"
		} else {
			stderr += l.Log + "\n"
		}
	}
	return stdout, stderr
}
