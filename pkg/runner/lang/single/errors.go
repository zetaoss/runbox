package single

type Error string

const (
	ErrNoSource        = Error("ErrNoSource")
	ErrDockerNew       = Error("ErrDockerNew")
	ErrDockerRun       = Error("ErrDockerRun")
	ErrInvalidLanguage = Error("ErrInvalidLanguage")
	ErrFileIO          = Error("ErrFileIO")
)

func (e Error) Error() string {
	return string(e)
}
