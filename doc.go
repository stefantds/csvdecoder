// csvdecoder is a tool for parsing and deserializing CSV values into Go objects.
// It follows the same usage pattern as the Rows scanning using database/sql package.
// It relies on encoding/csv for the actual csv parsing.
//
// csvdecoder allows to iterate through the CSV records (using 'Next')
// and scan the fields into target variables or fields of variables (using 'Scan').
// The methods 'Next' and 'Scan' are not thread-safe and are not expected to be called concurrently.
//
// csvdecoder supports converting CSV fields into any of the following types:
//	*string
//	*int, *int8, *int16, *int32, *int64
//	*uint, *uint8, *uint16, *uint32, *uint64
//	*bool
//	*float32, *float64
//	a slice of values. Note that the CSV field must be a valid JSON array. If not a JSON array, a custom decoder implementing the csvdecoder.Interface interface must be implemented.
//	an array of values. Note that the CSV field must be a valid JSON array. If not a JSON array, a custom decoder implementing the csvdecoder.Interface interface must be implemented.
//	a pointer to any type implementing the csvdecoder.Interface interface
//
// csvdecoder uses the same terminology as package encoding/csv:
// A csv file contains zero or more records. Each record contains one or more
// fields separated by the fields separator (the "comma"). The fields separator character
// can be configured to be another character than comma.
// Each record is separated by the newline character. The final record may
// optionally be followed by a newline character.
//
// The behavior of the decoder can be configured by passing one of following options when creating the decoder:
//	Comma: the character that separates values. The default value is comma.
//	IgnoreHeaders: if set to true, the first line will be ignored. This is useful when the CSV file contains a header line.
//	IgnoreUnmatchingFields: if set to true, the number of fields and scan targets are allowed to be different. By default, if they don't match exactly it will cause an error.
//
// See README.md for more info.
package csvdecoder
