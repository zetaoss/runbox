package docker

type Error string

const (
	NoError        Error = ""
	ErrCollectLogs Error = "ErrCollectLogs"
	ErrScanner     Error = "ErrScanner"
)
