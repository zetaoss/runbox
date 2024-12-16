package googling

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/zetaoss/runbox/pkg/runner/box"
	"k8s.io/utils/ptr"
)

type Googling struct {
	box *box.Box
}

func New(box *box.Box) *Googling {
	return &Googling{box}
}

func (g *Googling) Run(q string) (int, error) {
	command := fmt.Sprintf(`/opt/google/chrome/chrome --disable-gpu --dump-dom --headless --no-sandbox "https://google.com/search?q=%s" | grep -Po '(About|검색결과 약) \K([0-9,box]+)' | head -1`, url.QueryEscape(q))
	opts := &box.Opts{
		CollectStats:  ptr.To(false),
		CollectImages: false,
		Command:       command,
		Env:           []string{},
		Files:         []box.File{},
		Image:         "selenium/standalone-chrome:3.141.59",
	}
	result, err := g.box.Run(opts)
	if err != nil {
		return 0, fmt.Errorf("box.Run err: %w", err)
	}
	stdout, _ := result.StreamStrings()
	stdout = strings.TrimRight(stdout, "\n")
	stdout = strings.ReplaceAll(stdout, ",", "")
	number, err := strconv.Atoi(stdout)
	if err != nil {
		return 0, fmt.Errorf("strconv.Atoi err: %w", err)
	}
	return number, nil
}
