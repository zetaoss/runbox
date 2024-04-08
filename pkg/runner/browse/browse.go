package browse

import (
	"fmt"

	"github.com/zetaoss/runbox/pkg/docker"
	"github.com/zetaoss/runbox/pkg/util/runid"
)

const (
	ImageTag = "4.19.1-20240402"
)

type Input struct {
	RunID string `json:"-"`
	URL   string `json:"url"`
}

type Output struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func toOutput(result *docker.Result) *Output {
	var text = ""
	for _, l := range result.Logs {
		if l.Stream == "stdout" {
			text += l.Log + "\n"
		}
	}
	return &Output{Text: text}
}

func Run(url string) (*Output, error) {
	if url == "" {
		return nil, ErrNoURL
	}
	cli, err := docker.New()
	if err != nil {
		return nil, fmt.Errorf("docker.New err: %w", err)
	}
	var opts = &docker.Options{
		RunID:          runid.New("browse"),
		Image:          "selenium/standalone-chrome:" + ImageTag,
		Command:        fmt.Sprintf("/opt/google/chrome/chrome --headless --dump-dom --disable-gpu --no-sandbox '%s'", url),
		PidsLimit:      300,
		TimeoutSeconds: 10,
		OutputLimit:    8000,
	}
	result, err := cli.Run(*opts)
	if err != nil {
		return nil, fmt.Errorf("Run err: %w", err)
	}
	return toOutput(result), nil
}
