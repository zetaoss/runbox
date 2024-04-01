package docker

type Error string

const (
	NoError        Error = ""
	ErrClose       Error = "ErrClose"
	ErrCollectLogs Error = "ErrCollectLogs"
	ErrOpenFile    Error = "ErrOpenFile"
	ErrScanner     Error = "ErrScanner"
	ErrWriteString Error = "ErrWriteString"
)
