package notebook

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/zetaoss/runbox/pkg/docker"
	"github.com/zetaoss/runbox/pkg/runner/notebook/nbformat"
	"github.com/zetaoss/runbox/pkg/util/runid"
)

type Lang string

const (
	LangR      = Lang("r")
	LangPython = Lang("python")
)

func (l Lang) Language() string {
	switch l {
	case LangR:
		return "R"
	case LangPython:
		return "Python"
	default:
		return ""
	}
}

type Input struct {
	RunID string
	Lang  Lang
	Texts []string
}

type Output struct {
	Metadata nbformat.Metadata
	Cells    []nbformat.Cell
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
	fmt.Println(nb)
	return &Output{}, nil
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
		Metadata: nbformat.Metadata{
			LanguageInfo: nbformat.LanguageInfo{Name: input.Lang.Language()},
		},
		NBFormat:      4,
		NBFormatMinor: 4,
		Cells:         []nbformat.Cell{},
	}
	switch input.Lang {
	case LangR:
		nb.Metadata.Kernelspec = nbformat.Kernelspec{Name: "ir", DisplayName: "R", Language: "R"}
	case LangPython:
	default:
		return nil, ErrInvalidLanguage
	}
	runID := runid.New("notebook", string(input.Lang))
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
