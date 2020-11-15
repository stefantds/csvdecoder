package csvdecoder

import (
	"io"
	"io/ioutil"
	"strings"
	"unicode"
)

type readerCustomEscape struct {
	reader io.Reader
}

const (
	// defaultEscapeChar is the character used by the encoding/csv package to escape a quote
	defaultEscapeChar = '"'
	quote             = '"'
)

// NewReaderWithCustomEscape creates a reader that uses a custom character as escape character
// instead of the quote used by the encoding/csv Reader.
func NewReaderWithCustomEscape(r io.Reader, escapeChar rune) (*readerCustomEscape, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	tmpEscape := unicode.ReplacementChar // assuming this character doesn't appear in the string

	// replace the escaped escape character as it should not influence any quote
	// for simplicity we temporarily replace the escaped escape chars with a special character
	s := strings.ReplaceAll(
		string(b),
		string([]rune{escapeChar, escapeChar}),
		string(tmpEscape),
	)

	// replace the escaped quotes with the standard encoding/csv escape sequence
	s = strings.ReplaceAll(
		s,
		string([]rune{escapeChar, quote}),
		string([]rune{defaultEscapeChar, quote}),
	)

	// replace the back the escaped escape character
	s = strings.ReplaceAll(
		s,
		string(tmpEscape),
		string([]rune{escapeChar, escapeChar}),
	)

	return &readerCustomEscape{
		reader: strings.NewReader(s),
	}, nil
}

func (r readerCustomEscape) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}
