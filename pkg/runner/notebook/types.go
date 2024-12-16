package notebook

import (
	"github.com/jmnote/nbformat"
)

type Input struct {
	Lang    string   `json:"lang"`
	Sources []string `json:"sources"`
}

type Output nbformat.Output
type Outputs []Output

type Result struct {
	OutputsList []Outputs `json:"outputsList"`
	CPU         int       `json:"cpu"`
	MEM         int       `json:"mem"`
	Time        int       `json:"time"`
	Timedout    bool      `json:"timedout"`
	Stderr      string    `json:"-"`
}
