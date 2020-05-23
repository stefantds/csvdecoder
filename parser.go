package csvdecoder

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
)

type Parser struct {
	Reader           *csv.Reader
	config           ParserConfig
	currentRowValues []string
	lastErr          error
}

// ParserConfig allows to configure the parser.
type ParserConfig struct {
	Comma                  rune
	IgnoreHeaders          bool
	IgnoreUnmatchingFields bool
}

type Decoder interface {
	DecodeRecord(s string) error
}

func NewParserWithConfig(reader io.Reader, config ParserConfig) (*Parser, error) {
	return newParser(reader, config)
}

func NewParser(reader io.Reader) (*Parser, error) {
	return newParser(reader, ParserConfig{})
}

func newParser(reader io.Reader, config ParserConfig) (*Parser, error) {
	p := &Parser{
		Reader: csv.NewReader(reader),
		config: config,
	}

	if config.Comma != 0 {
		p.Reader.Comma = config.Comma
	}

	p.Reader.LazyQuotes = true
	p.Reader.FieldsPerRecord = -1

	if config.IgnoreHeaders {
		// consume the first line
		_, _ = p.Reader.Read()
	}

	return p, nil
}

// Scan copies the values in the current row into the values pointed
// at by dest.
// With the defult behaviour, it will throw an error if the number of values in dest
// is different from the number of values.
// If the `IgnoreUnmatchingFields` flag is set, it will ignore the records and the
// arguments that have no match.
//
// Scan converts columns read from the database into the following
// common types:
//
//    *string
//    *[]byte
//    *int, *int8, *int16, *int32, *int64
//    *uint, *uint8, *uint16, *uint32, *uint64
//    *bool
//    *float32, *float64
//    *interface{}
//    any type implementing Decoder
//
func (p *Parser) Scan(dest ...interface{}) error {
	if p.currentRowValues == nil {
		return errors.New("scan called without calling Next")
	}
	if !p.config.IgnoreUnmatchingFields && len(p.currentRowValues) != len(dest) {
		return fmt.Errorf("%w: got %d scan targets and %d records",
			ErrScanTargetsNotMatch,
			len(dest),
			len(p.currentRowValues),
		)
	}
	for i, val := range p.currentRowValues {
		if i >= len(dest) {
			// ignore the remaining records as they have no scan target
			break
		}
		err := convertAssignValue(dest[i], val)
		if err != nil {
			return fmt.Errorf("scan error on value index %d: %v", i, err)
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
func (p *Parser) Next() bool {
	var err error
	p.currentRowValues, err = p.Reader.Read()
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

// Err returns the error, if any, that was encountered during iteration.
func (p *Parser) Err() error {
	if p.lastErr != nil && p.lastErr != ErrEOF {
		return p.lastErr
	}
	return nil
}
