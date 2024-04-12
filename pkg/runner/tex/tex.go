package tex

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/zetaoss/runbox/pkg/docker"
	"github.com/zetaoss/runbox/pkg/runner/notebook/nbformat"
	"github.com/zetaoss/runbox/pkg/types"
	"github.com/zetaoss/runbox/pkg/util/runid"
)

type Input struct {
	RunID string `json:"-"`
	Lang  string `json:"lang"`
	Text  string `json:"text"`
}

type Output struct {
	Metadata    nbformat.Metadata   `json:"metadata"`
	CellOutputs [][]nbformat.Output `json:"cellOutputs"`
}

const (
	ErrInvalidLanguage = types.Error("ErrInvalidLanguage")
)

func toOutput(result *docker.Result) (*Output, error) {
	jsonText := ""
	for _, l := range result.Logs {
		if l.Stream == "stdout" {
			jsonText += l.Log
		}
	}
	var nb nbformat.Notebook
	if err := json.Unmarshal([]byte(jsonText), &nb); err != nil {
		return nil, err
	}
	var cellOutputs [][]nbformat.Output
	for _, cell := range nb.Cells {
		cellOutputs = append(cellOutputs, cell.Outputs)
	}
	return &Output{Metadata: nb.Metadata, CellOutputs: cellOutputs}, nil
}

func writeTexFile(input Input) ([]string, error) {
	bindSrcRoot := "/tmp/runbox/files/" + input.RunID
	if err := os.MkdirAll(bindSrcRoot, 0777); err != nil {
		return nil, fmt.Errorf("MkdirAll err: %w, name: %s", err, bindSrcRoot)
	}
	bindSrcFile := bindSrcRoot + "/my.tex"
	bindDstFile := "/home/user01/my.tex"
	if err := os.WriteFile(bindSrcFile, []byte(input.Text), 0644); err != nil {
		return nil, fmt.Errorf("WriteFile err: %w, name: %s", err, bindSrcFile)
	}
	return []string{bindSrcFile + ":" + bindDstFile}, nil
}

func Run(input Input) (*Output, error) {
	if input.Lang != "tex" && input.Lang != "latex" {
		return nil, ErrInvalidLanguage
	}
	input.RunID = runid.New("notebook", input.Lang)
	binds, err := writeTexFile(input)
	if err != nil {
		return nil, err
	}
	command := "touch oblivoir.sty && pdflatex -halt-on-error my.tex && echo %%%%$key%%%% && convert zeta.pdf p%d.png && find p?.png | xargs -i sh -c \"echo; base64 -w0 {}\""
	opts := docker.Options{
		RunID:     input.RunID,
		Image:     "jmnote/runbox:tex",
		Binds:     binds,
		PidsLimit: 40,
		Command:   command,
	}
	cli, err := docker.New()
	if err != nil {
		return nil, fmt.Errorf("docker new err: %w", err)
	}
	result, err := cli.Run(opts)
	if err != nil {
		return nil, fmt.Errorf("Run err: %w", err)
	}
	output, err := toOutput(result)
	if err != nil {
		return nil, err
	}
	return output, nil
}
