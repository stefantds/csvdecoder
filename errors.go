package csvdecoder

import (
	"errors"
)

var (
	// ErrEOF is thrown if the EOF is reached by the Next method.
	ErrEOF = errors.New("end of file reached")

	ErrScanTargetsNotMatch = errors.New("the number of scan targets does not match the number of csv records")

	errNilPtr = errors.New("destination is a nil pointer")
	errNotPtr = errors.New("destination not a pointer")
)
