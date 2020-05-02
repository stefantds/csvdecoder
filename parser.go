package csvdecoder

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
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

func (p *Parser) Next(data interface{}) (eof bool, err error) {

	decoderType := reflect.TypeOf((*Decoder)(nil)).Elem()

	dataReflected := reflect.ValueOf(data)

	// resolve pointers and interfaces
	for {
		if dataReflected.Kind() == reflect.Interface || dataReflected.Kind() == reflect.Ptr {
			dataReflected = dataReflected.Elem()
		} else {
			break
		}
	}

	records, err := p.Reader.Read()
	if err != nil {
		if err.Error() == "EOF" {
			return true, nil
		}
		return false, err
	}

	for i, record := range records {
		field := dataReflected.Field(i)

		switch field.Kind() {
		case reflect.Ptr:
			fieldType := field.Type()
			if field.Type().Implements(decoderType) {
				if field.IsZero() {
					// create a new value
					newValue := reflect.New(fieldType.Elem())
					field.Set(newValue)
				}

				fieldDecoder, _ := field.Addr().Elem().Interface().(Decoder)
				fieldDecoder.DecodeRecord(record)
			} else {
				return false, errors.New("Unsupported pointer to struct that doesn't implement the Decoder interface")
			}
		case reflect.Interface:
			fieldType := field.Type()
			if field.Type().Implements(decoderType) {
				if field.IsZero() {
					// create a new value
					newValue := reflect.New(fieldType.Elem())
					field.Set(newValue)
				}

				fieldDecoder, _ := field.Addr().Elem().Interface().(Decoder)
				fieldDecoder.DecodeRecord(record)
			} else {
				return false, errors.New("Unsupported pointer to struct that doesn't implement the Decoder interface")
			}
		case reflect.String:
			field.SetString(record)
		case reflect.Bool:
			if record != "" {
				col, err := strconv.ParseBool(record)
				if err != nil {
					return false, err
				}
				field.SetBool(col)
			}
		case reflect.Int:
			if record != "" {
				col, err := strconv.ParseInt(record, 10, 0)
				if err != nil {
					return false, err
				}
				field.SetInt(col)
			}
		case reflect.Int32:
			if record != "" {
				col, err := strconv.ParseInt(record, 10, 32)
				if err != nil {
					return false, err
				}
				field.SetInt(col)
			}
		case reflect.Int64:
			if record != "" {
				col, err := strconv.ParseInt(record, 10, 64)
				if err != nil {
					return false, err
				}
				field.SetInt(col)
			}
		case reflect.Float32:
			if record != "" {
				col, err := strconv.ParseFloat(record, 32)
				if err != nil {
					return false, err
				}
				field.SetFloat(col)
			}
		case reflect.Float64:
			if record != "" {
				col, err := strconv.ParseFloat(record, 64)
				if err != nil {
					return false, err
				}
				field.SetFloat(col)
			}
		case reflect.Slice:
			objType := field.Type()
			obj := reflect.New(objType).Interface()

			if err := json.NewDecoder(strings.NewReader(record)).Decode(&obj); err != nil {
				return false, fmt.Errorf("could not parse %s as JSON array: %w", record, err)
			}

			field.Set(reflect.ValueOf(obj).Elem())
		default:
			return false, errors.New("Unsupported field type")
		}
	}

	return false, nil
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
func (p *Parser) Nexty() bool {
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
