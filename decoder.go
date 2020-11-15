package csvdecoder

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
)

type Decoder struct {
	reader           *csv.Reader
	config           Config
	currentRowValues []string
	lastErr          error
}

// Config is a type that can be used to configure a decoder.
type Config struct {
	Comma                  rune // the character that separates values. Default value is comma.
	IgnoreHeaders          bool // if set to true, the first line will be ignored
	IgnoreUnmatchingFields bool // if set to true, the number of fields and scan targets are allowed to be different
	EscapeChar             rune // the character used to escape the quote character in quoted fields. The default is the quote itself.
}

// New returns a new CSV decoder that reads from r.
// The decoder can be given a custom configuration.
func NewWithConfig(r io.Reader, config Config) (*Decoder, error) {
	return newDecoder(r, config)
}

// New returns a new CSV decoder that reads from r
func New(r io.Reader) (*Decoder, error) {
	return newDecoder(r, Config{
		EscapeChar: defaultEscapeChar,
	})
}

func newDecoder(reader io.Reader, config Config) (*Decoder, error) {
	if config.EscapeChar != defaultEscapeChar {
		var err error
		reader, err = NewReaderWithCustomEscape(reader, config.EscapeChar)
		if err != nil {
			return nil, err
		}
	}

	p := &Decoder{
		reader: csv.NewReader(reader),
		config: config,
	}

	if config.Comma != 0 {
		p.reader.Comma = config.Comma
	}

	p.reader.LazyQuotes = true
	p.reader.FieldsPerRecord = -1

	if config.IgnoreHeaders {
		// consume the first line
		_, _ = p.reader.Read()
	}

	return p, nil
}

// Scan copies the values in the current row into the values pointed
// at by dest.
// With the default behavior, it will throw an error if the number of values in dest
// is different from the number of values. If the `IgnoreUnmatchingFields` flag is
// set, it will ignore the fields and the arguments that have no match.
//
// Scan converts columns read from the source into the following
// types:
//    *string
//    *int, *int8, *int16, *int32, *int64
//    *uint, *uint8, *uint16, *uint32, *uint64
//    *bool
//    *float32, *float64
//    a pointer to any type implementing Decoder interface
//    a slice of values that can be decoded from a JSON array by the JSON Decoder
//    an array of values that can be decoded from a JSON array by the JSON Decoder
//
// Scan must not be called concurrently.
func (p *Decoder) Scan(dest ...interface{}) error {
	switch {
	case errors.Is(p.lastErr, ErrEOF):
		return ErrEOF
	case p.lastErr != nil:
		return ErrReadingOccurred
	case p.currentRowValues == nil:
		return ErrNextNotCalled
	case !p.config.IgnoreUnmatchingFields && len(p.currentRowValues) != len(dest):
		return fmt.Errorf("%w: got %d scan targets and %d fields",
			ErrScanTargetsNotMatch,
			len(dest),
			len(p.currentRowValues),
		)
	}
	for i, val := range p.currentRowValues {
		if i >= len(dest) {
			// ignore the remaining fields as they have no scan target
			break
		}
		err := convertAssignValue(dest[i], val)
		if err != nil {
			return fmt.Errorf("scan error on value index %d: %w", i, err)
		}
	}
	return nil
}

// Next prepares the next result row for reading with the Scan method. It
// returns nil on success, or false if there is no next result row or an error
// happened while preparing it. Err should be consulted to distinguish between
// the two cases.
//
// Every call to Scan, even the first one, must be preceded by a call to Next.
// Next must not be called concurrently.
func (p *Decoder) Next() bool {
	var err error
	p.currentRowValues, err = p.reader.Read()
	if err != nil {
		if err.Error() == "EOF" {
			p.lastErr = ErrEOF
			return false
		}
		p.lastErr = fmt.Errorf("error while reading: %w", err)
		return false
	}
	return true
}

// Err returns the reading error, if any, that was encountered during iteration.
func (p *Decoder) Err() error {
	if p.lastErr != nil && p.lastErr != ErrEOF {
		return p.lastErr
	}
	return nil
}
