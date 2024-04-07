package single

type Input struct {
	RunID  string `json:"-"`
	Lang   string `json:"lang"`
	Source string `json:"source"`
}
