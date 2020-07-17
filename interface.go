package csvdecoder

// The Interface type describes the requirements
// for a type that can be decoded into a Go value
// by the csvdecoder.
// Any type that implements it may be used as a
// target in the Scan method.
//
// The Decode method allows to implement a custom
// decoding logic. If it returns an error, the
// parsing and decoding is stopped and the error
// is returned to the caller of Scan.
type Interface interface {
	DecodeRecord(s string) error
}
