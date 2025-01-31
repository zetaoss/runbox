package apperror

import "errors"

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrNoFiles         Error = "no files"
	ErrNoSources       Error = "no sources"
	ErrInvalidLanguage Error = "invalid language"
)

func IsAppError(err error) bool {
	var appErr Error
	return errors.As(err, &appErr)
}
