package notebook

import (
	"encoding/json"
	"fmt"

	"github.com/jmnote/nbformat"
	"github.com/zetaoss/runbox/pkg/errors"
	"github.com/zetaoss/runbox/pkg/runner/box"
	"k8s.io/utils/ptr"
)

type Notebook struct {
	box *box.Box
}

func New(box *box.Box) *Notebook {
	return &Notebook{box}
}

func (n *Notebook) Run(input Input) (*Result, error) {
	fileBody, err := toFileBody(input)
	if err != nil {
		return nil, fmt.Errorf("toFileBody err: %w", err)
	}
	opts := &box.Opts{
		CollectStats:  ptr.To(true),
		CollectImages: false,
		Command:       "jupyter nbconvert --execute --to notebook --allow-errors --stdout /tmp/runbox.ipynb",
		Files:         []box.File{{Name: "/tmp/runbox.ipynb", Body: fileBody}},
		Image:         fmt.Sprintf("jmnote/runbox:%s-notebook", input.Lang),
		WorkingDir:    "/tmp",
	}
	boxResult, err := n.box.Run(opts)
	if err != nil {
		return nil, fmt.Errorf("run err: %w", err)
	}
	result, err := toResult(boxResult)
	if err != nil {
		return nil, fmt.Errorf("toResult err: %w", err)
	}
	return result, nil
}

func toFileBody(input Input) (string, error) {
	var nb = nbformat.Notebook{
		Metadata: nbformat.Metadata{
			Kernelspec:   &nbformat.Kernelspec{},
			LanguageInfo: map[string]any{},
		},
		NbformatMinor: 4,
		Nbformat:      4,
		Cells:         []nbformat.Cell{},
	}
	switch input.Lang {
	case "r":
		nb.Metadata.Kernelspec.Name = "ir"
		nb.Metadata.LanguageInfo["name"] = "R"
	case "python":
		nb.Metadata.Kernelspec.Name = "python3"
		nb.Metadata.LanguageInfo["name"] = "python"
	default:
		return "", errors.ErrInvalidLanguage
	}
	for _, source := range input.Sources {
		cell := nbformat.Cell{
			CellType: "code",
			Metadata: map[string]any{},
			Source:   []string{source},
			Outputs:  []nbformat.Output{},
		}
		nb.Cells = append(nb.Cells, cell)
	}
	jsonBytes, err := json.Marshal(nb)
	if err != nil {
		return "", fmt.Errorf("json.Marshal err: %w", err)
	}
	return string(jsonBytes), nil
}

func toResult(boxResult *box.Result) (*Result, error) {
	outString, errString := boxResult.StreamStrings()
	nb, err := toNotebook(outString)
	if err != nil {
		return nil, fmt.Errorf("toNotebook err: %w", err)
	}
	outputsList := make([]Outputs, len(nb.Cells))
	for i, cell := range nb.Cells {
		outputs := Outputs{}
		for _, output := range cell.Outputs {
			outputs = append(outputs, Output(output))
		}
		outputsList[i] = outputs
	}
	result := &Result{
		OutputsList: outputsList,
		CPU:         boxResult.CPU,
		MEM:         boxResult.MEM,
		Time:        boxResult.Time,
		Timedout:    boxResult.Timedout,
		Stderr:      errString,
	}
	return result, nil
}

func toNotebook(str string) (*nbformat.Notebook, error) {
	var nb nbformat.Notebook
	if err := json.Unmarshal([]byte(str), &nb); err != nil {
		return nil, fmt.Errorf("json.Unmarshal err: %w", err)
	}
	return &nb, nil
}
