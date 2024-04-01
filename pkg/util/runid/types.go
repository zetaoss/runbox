package runid

type Kind string

const (
	KindDocker Kind = "D"
	KindSingle Kind = "S"
	KindMulti  Kind = "M"
)
