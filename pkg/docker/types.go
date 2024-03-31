package docker

type Kind int

const (
	KindDocker Kind = iota
	KindSingle
	KindMulti
)

type Options struct {
	Kind           Kind
	Name           string
	Image          string
	Shell          string
	Command        string
	Env            []string
	PidsLimit      int
	TimeoutSeconds int
	Binds          []string
	WorkingDir     string
}

type Log struct {
	Log    string
	Stream string
}

type Result struct {
	Logs     []Log
	ExitCode int
}
