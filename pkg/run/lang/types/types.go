package types

type File struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	IsMain  bool   `json:"isMain,omitempty"`
}

type SingleInput struct {
	RunID  string `json:"-"`
	Lang   string `json:"lang"`
	Source string `json:"source"`
	Hash   string `json:"hash"`
}

type MultiInput struct {
	RunID string `json:"-"`
	Lang  string `json:"lang"`
	Files []File `json:"files"`
	Hash  string `json:"hash"`
}

type Output struct {
	Logs    []string `json:"logs"`
	Images  []string `json:"images,omitempty"`
	Timeout bool     `json:"timeout,omitempty"`
	Time    string   `json:"time"`
	CPU     float32  `json:"cpu"`
	MEM     float32  `json:"mem"`
}

type RunOpts struct {
	Command          string
	Env              []string
	FileName         string
	FileExt          string
	PidsLimit        int
	ModifySourceFunc func(string) string
	Shell            string
	TimeoutCommand   string
	TimeoutSeconds   int
	VolSrcRoot       string
	VolSubPath       string
	WorkingDir       string
}
