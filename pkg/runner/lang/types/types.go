package types

const (
	WarnTimeout            = "WarnTimeout"
	WarnOutputLimitReached = "WarnOutputLimitReached"
)

type Output struct {
	Logs     []string `json:"logs"`
	Images   []string `json:"images,omitempty"`
	Time     string   `json:"time"`
	CPU      float32  `json:"cpu"`
	MEM      float32  `json:"mem"`
	Warnings []string `json:"warnings,omitempty"`
}

type RunOpts struct {
	Command          string
	Env              []string
	Image            string
	FileName         string
	FileExt          string
	ModifySourceFunc func(string) string
	Postflight       func(*Output)
	Shell            string
	TimeoutCommand   string
	TimeoutSeconds   int
	VolSrcRoot       string
	VolSubPath       string
	WorkingDir       string
}
