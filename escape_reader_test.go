package csvdecoder

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestEscapeReader(t *testing.T) {
	for _, tc := range []struct {
		name           string
		input          string
		escapeChar     rune
		expectedResult string
	}{
		{
			name:           "should work without anything to escape",
			input:          "my example string",
			escapeChar:     '_',
			expectedResult: "my example string",
		},
		{
			name:           "should replace escaping quotes",
			input:          `my _"example_" string`,
			escapeChar:     '_',
			expectedResult: `my ""example"" string`,
		},
		{
			name:           "should not replace escaping chars without quotes",
			input:          "my _example_ string",
			escapeChar:     '_',
			expectedResult: "my _example_ string",
		},
		{
			name:           "should ignore escaped escaped chars",
			input:          `my example string__"`,
			escapeChar:     '_',
			expectedResult: `my example string__"`,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r, err := NewReaderWithCustomEscape(strings.NewReader(tc.input), tc.escapeChar)
			if err != nil {
				t.Fatal(err)
			}

			result, err := ioutil.ReadAll(r)
			if err != nil {
				t.Fatal(err)
			}

			if string(result) != tc.expectedResult {
				t.Errorf("expected value '%s' got '%s'", tc.expectedResult, result)
			}
		})
	}
}
