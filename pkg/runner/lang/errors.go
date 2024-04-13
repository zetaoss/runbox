package lang

type Error string

const (
	ErrNoFiles         = Error("ErrNoFiles")
	ErrDockerNew       = Error("ErrDockerNew")
	ErrDockerRun       = Error("ErrDockerRun")
	ErrInvalidLanguage = Error("ErrInvalidLanguage")
	ErrFileIO          = Error("ErrFileIO")
)

func (e Error) Error() string {
	return string(e)
}
