package csvdecoder

import (
	"errors"
	"strings"
	"testing"
)

func TestIgnoreUnmatchingFields(t *testing.T) {
	var strVal string
	var intVal int
	var anotherIntVal int

	for _, tc := range []struct {
		name          string
		config        Config
		data          string
		scanTargets   []interface{}
		expectedError error
	}{
		{
			name: "should work when numbers match and flag is false",
			config: Config{
				IgnoreUnmatchingFields: false,
			},
			data:          "rec,2\n",
			scanTargets:   []interface{}{&strVal, &intVal},
			expectedError: nil,
		},
		{
			name: "should work when numbers match and flag is true",
			config: Config{
				IgnoreUnmatchingFields: true,
			},
			data:          "rec,2\n",
			scanTargets:   []interface{}{&strVal, &intVal},
			expectedError: nil,
		},
		{
			name:          "should work when numbers match with default config",
			config:        Config{},
			data:          "rec,2\n",
			scanTargets:   []interface{}{&strVal, &intVal},
			expectedError: nil,
		},
		{
			name: "should work with more targets when the flag is true",
			config: Config{
				IgnoreUnmatchingFields: true,
			},
			data:          "rec,2\n",
			scanTargets:   []interface{}{&strVal, &intVal, &anotherIntVal},
			expectedError: nil,
		},
		{
			name: "should work with more records when the flag is true",
			config: Config{
				IgnoreUnmatchingFields: true,
			},
			data:          "rec,2\n",
			scanTargets:   []interface{}{&strVal},
			expectedError: nil,
		},
		{
			name: "should fail with more targets when the flag is false",
			config: Config{
				IgnoreUnmatchingFields: false,
			},
			data:          "rec,2\n",
			scanTargets:   []interface{}{&strVal, &intVal, &anotherIntVal},
			expectedError: ErrScanTargetsNotMatch,
		},
		{
			name: "should fail with more records when the flag is false",
			config: Config{
				IgnoreUnmatchingFields: false,
			},
			data:          "rec,2\n",
			scanTargets:   []interface{}{&strVal},
			expectedError: ErrScanTargetsNotMatch,
		},
		{
			name:          "should fail with more targets with the default config",
			config:        Config{},
			data:          "rec,2\n",
			scanTargets:   []interface{}{&strVal, &intVal, &anotherIntVal},
			expectedError: ErrScanTargetsNotMatch,
		},
		{
			name:          "should fail with more records with the default config",
			config:        Config{},
			data:          "rec,2\n",
			scanTargets:   []interface{}{&strVal},
			expectedError: ErrScanTargetsNotMatch,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			d, err := NewWithConfig(strings.NewReader(tc.data), tc.config)
			if err != nil {
				t.Fatalf("could not create d: %w", err)
			}

			for d.Next() {
				if err := d.Scan(tc.scanTargets...); !errors.Is(err, tc.expectedError) {
					t.Errorf("expected '%s', got '%v'", tc.expectedError, err)
				}
			}
			if d.Err() != nil {
				t.Error(err)
			}
		})
	}
}
