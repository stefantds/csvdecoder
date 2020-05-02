package csvdecoder

import (
	"errors"
)

var (
	// ErrEOF is thrown if the EOF is reached by the Next method.
	ErrEOF = errors.New("end of file reached")

	errNilPtr = errors.New("destination is a nil pointer")
	errNotPtr = errors.New("destination not a pointer")
)
