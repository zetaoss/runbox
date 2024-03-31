package types

type Error string

const (
	ErrDockerNew       = Error("docker new error")
	ErrDockerRun       = Error("docker run error")
	ErrInvalidLanguage = Error("invalid language")
	ErrFileIO          = Error("file IO error")
)

func (e Error) Error() string {
	return string(e)
}
