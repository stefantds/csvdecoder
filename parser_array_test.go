package csvdecoder

import (
	"reflect"
	"strings"
	"testing"
)

func TestIntArray(t *testing.T) {
	type TestRow struct {
		Field []int
	}

	for _, tc := range []struct {
		name      string
		RowStruct TestRow
		data      string
		expected  []int
	}{
		{
			name:      "should work for an empty array",
			RowStruct: TestRow{},
			data:      "[]\n",
			expected:  []int{},
		},
		{
			name:      "should work for a single-valued array",
			RowStruct: TestRow{},
			data:      "[1]\n",
			expected:  []int{1},
		},
		{
			name:      "should work for an array with positive values",
			RowStruct: TestRow{},
			data:      "[1, 2, 3]\n",
			expected:  []int{1, 2, 3},
		},
		{
			name:      "should work for an array with negative values",
			RowStruct: TestRow{},
			data:      "[-1, -2, -3]\n",
			expected:  []int{-1, -2, -3},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			parser, err := NewParser(strings.NewReader(tc.data), &ParserConfig{IgnoreHeaders: false, Comma: '\t'})
			if err != nil {
				t.Fatalf("could not create parser: %w", err)
			}

			for {
				eof, err := parser.Next(&tc.RowStruct)
				if err != nil {
					t.Error(err)
				}

				if !reflect.DeepEqual(tc.RowStruct.Field, tc.expected) {
					t.Errorf("expected value '%v' got '%v'", tc.expected, tc.RowStruct.Field)
				}

				if eof {
					break
				}
			}
		})
	}
}

func TestMultiLevelIntArray(t *testing.T) {
	type TestRow struct {
		Field [][][]int
	}

	for _, tc := range []struct {
		name      string
		RowStruct TestRow
		data      string
		expected  [][][]int
	}{
		{
			name:      "should work for an empty array",
			RowStruct: TestRow{},
			data:      "[]\n",
			expected:  [][][]int{},
		},
		{
			name:      "should work for an empty array in an array",
			RowStruct: TestRow{},
			data:      "[[]]\n",
			expected:  [][][]int{{}},
		},
		{
			name:      "should work for items with same level",
			RowStruct: TestRow{},
			data:      "[[], [[1, 2],[3, 4]]]\n",
			expected:  [][][]int{{}, {{1, 2}, {3, 4}}},
		},
		{
			name:      "should work for items with different levels",
			RowStruct: TestRow{},
			data:      "[[], [[1, 2]]]\n",
			expected:  [][][]int{{}, {{1, 2}}},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			parser, err := NewParser(strings.NewReader(tc.data), &ParserConfig{IgnoreHeaders: false, Comma: '\t'})
			if err != nil {
				t.Fatalf("could not create parser: %w", err)
			}

			for {
				eof, err := parser.Next(&tc.RowStruct)
				if err != nil {
					t.Error(err)
				}

				if !reflect.DeepEqual(tc.RowStruct.Field, tc.expected) {
					t.Errorf("expected value '%v' got '%v'", tc.expected, tc.RowStruct.Field)
				}

				if eof {
					break
				}
			}
		})
	}
}
