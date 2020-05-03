package csvdecoder

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
)

type Parser struct {
	Reader           *csv.Reader
	config           *ParserConfig
	currentRowValues []string
	lastErr          error
}

// ParserConfig has information for parser
type ParserConfig struct {
	Comma         rune
	IgnoreHeaders bool
}

type Decoder interface {
	DecodeRecord(s string) error
}

func NewParser(reader io.Reader, config *ParserConfig) (*Parser, error) {
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
// at by dest. The number of values in dest can be different than the
// number of values. In this case Scan will only fill the available
// values.
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
		return errors.New("Scan called without calling Next")
	}
	for i, val := range p.currentRowValues {
		err := convertAssignValue(dest[i], val)
		if err != nil {
			return fmt.Errorf(`Scan error on value index %d: %v`, i, err)
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
