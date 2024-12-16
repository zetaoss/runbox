package errors

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrNoFiles         Error = "no files"
	ErrInvalidLanguage Error = "invalid language"
)
