package run

type Error string

const (
	NoError     Error = ""
	ErrBindJSON Error = "ErrBindJSON"
	ErrUnknown  Error = "ErrUnknown"
)
