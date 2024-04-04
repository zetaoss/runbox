package docker

type Options struct {
	RunID          string
	Image          string
	Shell          string
	Command        string
	Env            []string
	PidsLimit      int
	TimeoutSeconds int
	OutputLimit    int
	Binds          []string
	WorkingDir     string
}

type Log struct {
	Log    string
	Stream string
}

type Result struct {
	ExitCode           int
	Logs               []Log
	OutputLimitReached bool
}
