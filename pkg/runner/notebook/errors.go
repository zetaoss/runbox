package notebook

type Error string

const (
	ErrInvalidLanguage = Error("ErrInvalidLanguage")
)

func (e Error) Error() string {
	return string(e)
}
