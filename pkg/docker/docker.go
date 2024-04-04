package docker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zetaoss/runbox/pkg/util"
	"github.com/zetaoss/runbox/pkg/util/runid"
	"k8s.io/klog/v2"
)

var (
	fakeErr = NoError
)

type Config struct {
	DataDir string
}

type Docker struct {
	DataDir     string
	HistoryFile string
}

func New() (*Docker, error) {
	return NewWithConfig(Config{})
}

func NewWithConfig(cfg Config) (*Docker, error) {
	dataDir := cfg.DataDir
	if dataDir == "" {
		dataDir = "/data"
	}
	if err := os.MkdirAll(dataDir, 0644); err != nil {
		return nil, fmt.Errorf("MkdirAll err: %w", err)
	}
	historyFile := dataDir + "/history.txt"
	return &Docker{
		DataDir:     dataDir,
		HistoryFile: historyFile,
	}, nil
}

func (d *Docker) saveHistory(command string) error {
	f, err := os.OpenFile(d.HistoryFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil || fakeErr == ErrOpenFile {
		return fmt.Errorf("OpenFile err: %w", err)
	}
	defer f.Close()
	if _, err = f.WriteString(command + "\n"); err != nil || fakeErr == ErrWriteString {
		return fmt.Errorf("WriteString err: %w", err)
	}
	if err := f.Close(); err != nil || fakeErr == ErrClose {
		return fmt.Errorf("close err: %w", err)
	}
	return nil
}

func (d *Docker) Run(opts Options) (*Result, error) {
	if opts.Image == "" {
		return nil, fmt.Errorf("no image")
	}
	if opts.Command == "" {
		return nil, fmt.Errorf("no command")
	}
	// default
	if opts.RunID == "" {
		opts.RunID = runid.New("docker")
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
	command += " --name=" + opts.RunID
	command += " --pids-limit=" + fmt.Sprintf("%d", opts.PidsLimit)
	for _, bind := range opts.Binds {
		command += " -v " + bind
	}
	command += " " + opts.Image
	command += " " + opts.Command

	// save history
	if err := d.saveHistory(command); err != nil {
		klog.Warningf("saveHistory err: %s", err.Error())
	}
	_, _, exitCode := util.Run(command)
	stdout, _, _ := util.Run("docker inspect --format={{.Id}} " + opts.RunID)
	containerID := strings.TrimRight(stdout, "\n")
	defer func() {
		util.Run("docker rm -f " + containerID)
	}()
	logs, outputLimitReached, err := d.collectLogs(containerID, opts.OutputLimit)
	if err != nil || fakeErr == ErrCollectLogs {
		return nil, fmt.Errorf("collectLogs err: %w", err)
	}
	return &Result{ExitCode: exitCode, Logs: logs, OutputLimitReached: outputLimitReached}, nil
}

func (d *Docker) collectLogs(containerID string, outputLimit int) ([]Log, bool, error) {
	logFilePath := fmt.Sprintf("/var/lib/docker/containers/%s/%s-json.log", containerID, containerID)

	file, err := os.Open(logFilePath)
	if err != nil {
		return nil, false, err
	}
	defer file.Close()

	var logs []Log = []Log{}
	var totalLength int
	var limitReached bool

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var logLine Log
		if err := json.Unmarshal(scanner.Bytes(), &logLine); err != nil {
			log.Printf("warn: unmarshal: %s", err)
			continue
		}

		lineLength := len([]rune(logLine.Log))
		if outputLimit != 0 && totalLength+lineLength > outputLimit {
			limitReached = true
			logLine.Log = logLine.Log[:outputLimit-totalLength]
			logs = append(logs, logLine)
			break
		}

		totalLength += lineLength
		logs = append(logs, logLine)
	}

	if err := scanner.Err(); err != nil || fakeErr == ErrScanner {
		log.Printf("warn: scanner: %s", err)
	}
	return logs, limitReached, nil
}
