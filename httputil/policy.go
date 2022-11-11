package httputil

//go:generate stringer -type=Policy
type Policy int

const (
	PolicyFirst  Policy = 0
	PolicyRandom Policy = 1
)
