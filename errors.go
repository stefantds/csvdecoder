package csvdecoder

import (
	"errors"
)

var (
	ErrEOF                 = errors.New("end of file reached") // ErrEOF is thrown if the EOF is reached by the Next method.
	ErrScanTargetsNotMatch = errors.New("the number of scan targets does not match the number of csv records")
	ErrReadingOccurred     = errors.New("can't continue after a reading error")
	ErrNextNotCalled       = errors.New("scan called without calling Next")

	errNilPtr = errors.New("destination is a nil pointer")
	errNotPtr = errors.New("destination not a pointer")
)
