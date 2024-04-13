package notebook

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/zetaoss/runbox/pkg/docker"
	"github.com/zetaoss/runbox/pkg/runner/notebook/nbformat"
	"github.com/zetaoss/runbox/pkg/util/runid"
)

type Input struct {
	RunID     string     `json:"-"`
	Lang      string     `json:"lang"`
	CellTexts [][]string `json:"cellTexts"`
}

type Output struct {
	Metadata    nbformat.Metadata   `json:"metadata"`
	CellOutputs [][]nbformat.Output `json:"cellOutputs"`
}

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

func writeNotebookFile(nb nbformat.Notebook, runID string) ([]string, error) {
	bindSrcRoot := "/data/files/" + runID
	if err := os.MkdirAll(bindSrcRoot, 0777); err != nil {
		return nil, fmt.Errorf("MkdirAll err: %w, name: %s", err, bindSrcRoot)
	}
	bindSrcFile := bindSrcRoot + "/my.ipynb"
	bindDstFile := "/home/jovyan/my.ipynb"
	file, err := os.Create(bindSrcFile)
	if err != nil {
		return nil, fmt.Errorf("os.Create err: %w", err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(nb); err != nil {
		return nil, fmt.Errorf("encode err: %w", err)
	}
	return []string{bindSrcFile + ":" + bindDstFile}, nil
}

func Run(input Input) (*Output, error) {
	var nb = nbformat.Notebook{
		NBFormat:      4,
		NBFormatMinor: 4,
	}
	switch input.Lang {
	case "r":
		nb.Metadata.Kernelspec.Name = "ir"
		nb.Metadata.LanguageInfo.Name = "R"
	case "python":
		nb.Metadata.LanguageInfo.Name = "python"
	default:
		return nil, ErrInvalidLanguage
	}
	var cells = []nbformat.Cell{}
	for _, ct := range input.CellTexts {
		cells = append(cells, nbformat.Cell{
			CellType: "code",
			Source:   ct,
			Outputs:  []nbformat.Output{},
		})
	}
	nb.Cells = cells
	runID := runid.New("notebook", input.Lang)
	binds, err := writeNotebookFile(nb, runID)
	if err != nil {
		return nil, err
	}
	command := "jupyter nbconvert --execute --to notebook --allow-errors --stdout my.ipynb"
	opts := docker.Options{
		RunID:     runID,
		Image:     fmt.Sprintf("jmnote/runbox:%s-notebook", input.Lang),
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
