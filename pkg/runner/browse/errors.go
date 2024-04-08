package browse

type Error string

const (
	ErrNoURL = Error("ErrNoURL")
)

func (e Error) Error() string {
	return string(e)
}
