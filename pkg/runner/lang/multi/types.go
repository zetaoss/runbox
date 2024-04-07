package multi

type File struct {
	Name string `json:"name"`
	Text string `json:"text"`
	Main bool   `json:"main,omitempty"`
}

type Input struct {
	RunID string `json:"-"`
	Lang  string `json:"lang"`
	Files []File `json:"files"`
}
