package types

type Error string

const (
	ErrNoSource        = Error("no source")
	ErrDockerNew       = Error("docker new error")
	ErrDockerRun       = Error("docker run error")
	ErrInvalidLanguage = Error("invalid language")
	ErrFileIO          = Error("file IO error")
)

func (e Error) Error() string {
	return string(e)
}
