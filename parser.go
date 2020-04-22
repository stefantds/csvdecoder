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
	Reader *csv.Reader
	config *ParserConfig
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
				return false, errors.New("Unsupported pointer to struct that doesn't implement deserializer")
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
