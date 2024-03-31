package docker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/zetaoss/runbox/pkg/util"
)

var fakeErr = NoError

type Docker struct {
}

func New() *Docker {
	return &Docker{}
}

func (d *Docker) Run(opts Options) (*Result, error) {
	// default
	if opts.Name == "" {
		opts.Name = util.NewHash(10)
	}
	if opts.Shell == "" {
		opts.Shell = "sh"
	}
	if opts.TimeoutSeconds == 0 {
		opts.TimeoutSeconds = 9
	}
	if opts.PidsLimit == 0 {
		opts.PidsLimit = 15
	}

	command := "docker run"
	command += " --name=" + opts.Name
	command += " --pids-limit=" + fmt.Sprintf("%d", opts.PidsLimit)
	for _, bind := range opts.Binds {
		command += " -v " + bind
	}
	command += " " + opts.Image
	command += " " + opts.Command
	_, _, exitCode := util.Run(command)
	stdout, _, _ := util.Run("docker inspect  --format={{.Id}} " + opts.Name)
	containerID := strings.TrimRight(stdout, "\n")
	defer func() {
		util.Run("docker rm -f " + containerID)
	}()
	logs, err := d.collectLogs(containerID)
	if err != nil || fakeErr == ErrCollectLogs {
		return nil, fmt.Errorf("collectLogs err: %w", err)
	}
	return &Result{Logs: logs, ExitCode: exitCode}, nil
}

func (d *Docker) collectLogs(containerID string) ([]Log, error) {
	logFilePath := fmt.Sprintf("/var/lib/docker/containers/%s/%s-json.log", containerID, containerID)
	command := "cat " + logFilePath
	stdout, _, _ := util.Run(command)
	var logs = []Log{}
	scanner := bufio.NewScanner(strings.NewReader(stdout))
	for scanner.Scan() {
		var logLine Log
		if err := json.Unmarshal(scanner.Bytes(), &logLine); err != nil {
			log.Printf("warn: unmarshal: %s", err)
			continue
		}
		logs = append(logs, logLine)
	}
	if err := scanner.Err(); err != nil || fakeErr == ErrScanner {
		log.Printf("warn: scanner: %s", err)
	}
	return logs, nil
}
