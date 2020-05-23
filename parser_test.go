package csvdecoder

import (
	"errors"
	"strings"
	"testing"
)

func TestParserIgnoreUnmatchingFields(t *testing.T) {
	var strVal string
	var intVal int
	var anotherIntVal int

	for _, tc := range []struct {
		name          string
		config        ParserConfig
		data          string
		scanTargets   []interface{}
		expectedError error
	}{
		{
			name: "should work when numbers match and flag is false",
			config: ParserConfig{
				IgnoreUnmatchingFields: false,
			},
			data:          "rec,2\n",
			scanTargets:   []interface{}{&strVal, &intVal},
			expectedError: nil,
		},
		{
			name: "should work when numbers match and flag is true",
			config: ParserConfig{
				IgnoreUnmatchingFields: true,
			},
			data:          "rec,2\n",
			scanTargets:   []interface{}{&strVal, &intVal},
			expectedError: nil,
		},
		{
			name:          "should work when numbers match with default config",
			config:        ParserConfig{},
			data:          "rec,2\n",
			scanTargets:   []interface{}{&strVal, &intVal},
			expectedError: nil,
		},
		{
			name: "should work with more targets when the flag is true",
			config: ParserConfig{
				IgnoreUnmatchingFields: true,
			},
			data:          "rec,2\n",
			scanTargets:   []interface{}{&strVal, &intVal, &anotherIntVal},
			expectedError: nil,
		},
		{
			name: "should work with more records when the flag is true",
			config: ParserConfig{
				IgnoreUnmatchingFields: true,
			},
			data:          "rec,2\n",
			scanTargets:   []interface{}{&strVal},
			expectedError: nil,
		},
		{
			name: "should fail with more targets when the flag is false",
			config: ParserConfig{
				IgnoreUnmatchingFields: false,
			},
			data:          "rec,2\n",
			scanTargets:   []interface{}{&strVal, &intVal, &anotherIntVal},
			expectedError: ErrScanTargetsNotMatch,
		},
		{
			name: "should fail with more records when the flag is false",
			config: ParserConfig{
				IgnoreUnmatchingFields: false,
			},
			data:          "rec,2\n",
			scanTargets:   []interface{}{&strVal},
			expectedError: ErrScanTargetsNotMatch,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			parser, err := NewParserWithConfig(strings.NewReader(tc.data), tc.config)
			if err != nil {
				t.Fatalf("could not create parser: %w", err)
			}

			for parser.Next() {
				if err := parser.Scan(tc.scanTargets...); !errors.Is(err, tc.expectedError) {
					t.Errorf("expected '%s', got '%v'", tc.expectedError, err)
				}
			}
			if parser.Err() != nil {
				t.Error(err)
			}
		})
	}
}
