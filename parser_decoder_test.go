package csvdecoder

import (
	"strings"
	"testing"
)

type MyDecoder struct {
	value string
}

func (d *MyDecoder) DecodeRecord(value string) error {
	d.value = value + "!"
	return nil
}

func (d *MyDecoder) DecodedValue() string {
	return d.value
}

func TestDecoderStruct(t *testing.T) {
	type TestRowWithInterface struct {
		Field      Decoder
		OtherField int
	}

	for _, tc := range []struct {
		name     string
		dest     MyDecoder
		data     string
		expected string
	}{
		{
			name:     "should work for a struct using the Decoder interface",
			dest:     MyDecoder{},
			data:     "record1\n",
			expected: "record1!",
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			parser, err := NewParser(strings.NewReader(tc.data), &ParserConfig{IgnoreHeaders: false, Comma: '\t'})
			if err != nil {
				t.Fatalf("could not create parser: %w", err)
			}

			for parser.Next() {
				if err := parser.Scan(&tc.dest); err != nil {
					t.Error(err)
				}
				if tc.dest.DecodedValue() != tc.expected {
					t.Errorf("expected value '%s' got '%s'", tc.expected, tc.dest.DecodedValue())
				}
			}
			if parser.Err() != nil {
				t.Error(err)
			}
		})
	}
}

func TestDecoderPointer(t *testing.T) {
	for _, tc := range []struct {
		name     string
		dest     *MyDecoder
		data     string
		expected string
	}{
		{
			name:     "should work for a struct holding a pointer to a decoder",
			dest:     &MyDecoder{},
			data:     "record1\n",
			expected: "record1!",
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			parser, err := NewParser(strings.NewReader(tc.data), &ParserConfig{IgnoreHeaders: false, Comma: '\t'})
			if err != nil {
				t.Fatalf("could not create parser: %w", err)
			}

			for parser.Next() {
				if err := parser.Scan(&tc.dest); err != nil {
					t.Error(err)
				}
				if tc.dest.DecodedValue() != tc.expected {
					t.Errorf("expected value '%s' got '%s'", tc.expected, tc.dest.DecodedValue())
				}
			}
			if parser.Err() != nil {
				t.Error(err)
			}
		})
	}
}

func TestDecoderDoublePointer(t *testing.T) {
	myDec := &MyDecoder{}
	for _, tc := range []struct {
		name     string
		dest     **MyDecoder
		data     string
		expected string
	}{
		{
			name:     "should work for a struct holding a double pointer to a decoder",
			dest:     &myDec,
			data:     "record1\n",
			expected: "record1!",
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			parser, err := NewParser(strings.NewReader(tc.data), &ParserConfig{IgnoreHeaders: false, Comma: '\t'})
			if err != nil {
				t.Fatalf("could not create parser: %w", err)
			}

			for parser.Next() {
				if err := parser.Scan(&tc.dest); err != nil {
					t.Error(err)
				}
				if (*tc.dest).DecodedValue() != tc.expected {
					t.Errorf("expected value '%s' got '%s'", tc.expected, (*tc.dest).DecodedValue())
				}
			}
			if parser.Err() != nil {
				t.Error(err)
			}
		})
	}
}

func TestDecoderInterface(t *testing.T) {
	for _, tc := range []struct {
		name     string
		dest     Decoder
		data     string
		expected string
	}{
		{
			name:     "should work for a struct using the Decoder interface",
			dest:     &MyDecoder{},
			data:     "record1\n",
			expected: "record1!",
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			parser, err := NewParser(strings.NewReader(tc.data), &ParserConfig{IgnoreHeaders: false})
			if err != nil {
				t.Fatalf("could not create parser: %w", err)
			}

			for parser.Next() {
				if err := parser.Scan(tc.dest); err != nil {
					t.Error(err)
				}
				if tc.dest.(*MyDecoder).DecodedValue() != tc.expected {
					t.Errorf("expected value '%s' got '%s'", tc.expected, tc.dest.(*MyDecoder).DecodedValue())
				}
			}
			if parser.Err() != nil {
				t.Error(err)
			}
		})
	}
}
