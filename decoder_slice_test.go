package csvdecoder

import (
	"reflect"
	"strings"
	"testing"
)

func TestIntSlice(t *testing.T) {
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
			d, err := NewWithConfig(strings.NewReader(tc.data), Config{IgnoreHeaders: false, Comma: '\t'})
			if err != nil {
				t.Fatalf("could not create d: %w", err)
			}

			for d.Next() {
				err := d.Scan(&tc.RowStruct.Field)
				if err != nil {
					t.Error(err)
				}

				if !reflect.DeepEqual(tc.RowStruct.Field, tc.expected) {
					t.Errorf("expected value '%v' got '%v'", tc.expected, tc.RowStruct.Field)
				}
			}
			if d.Err() != nil {
				t.Errorf("d error: %w", err)
			}
		})
	}
}

func TestMultiLevelIntSlice(t *testing.T) {
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
			d, err := NewWithConfig(strings.NewReader(tc.data), Config{IgnoreHeaders: false, Comma: '\t'})
			if err != nil {
				t.Fatalf("could not create d: %w", err)
			}

			for d.Next() {
				err := d.Scan(&tc.RowStruct.Field)
				if err != nil {
					t.Error(err)
				}

				if !reflect.DeepEqual(tc.RowStruct.Field, tc.expected) {
					t.Errorf("expected value '%v' got '%v'", tc.expected, tc.RowStruct.Field)
				}
			}
			if d.Err() != nil {
				t.Errorf("d error: %w", err)
			}
		})
	}
}

func TestStructSlice(t *testing.T) {
	type MyStruct struct {
		A int    `json:"a"`
		B int32  `json:"b"`
		C int64  `json:"c"`
		D string `json:"d"`
		E bool   `json:"e"`
	}

	type TestRow struct {
		Field []MyStruct
	}

	for _, tc := range []struct {
		name      string
		RowStruct TestRow
		data      string
		expected  []MyStruct
	}{
		{
			name:      "should work for an empty array",
			RowStruct: TestRow{},
			data:      "[]\n",
			expected:  []MyStruct{},
		},
		{
			name:      "should work for an array with one value",
			RowStruct: TestRow{},
			data:      `[{"a":1, "b":2, "c": 3, "d":"value1", "e": true}]`,
			expected: []MyStruct{
				{
					A: 1,
					B: 2,
					C: 3,
					D: "value1",
					E: true,
				},
			},
		},
		{
			name:      "should work for an array with multiple values",
			RowStruct: TestRow{},
			data:      `[{"a":1, "b":2, "c": 3, "d":"value1", "e": true}, {"a":4, "b":5, "c": 6, "d":"value2", "e": false}]`,
			expected: []MyStruct{
				{
					A: 1,
					B: 2,
					C: 3,
					D: "value1",
					E: true,
				},
				{
					A: 4,
					B: 5,
					C: 6,
					D: "value2",
					E: false,
				},
			},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			d, err := NewWithConfig(strings.NewReader(tc.data), Config{IgnoreHeaders: false, Comma: '\t'})
			if err != nil {
				t.Fatalf("could not create d: %w", err)
			}

			for d.Next() {
				err := d.Scan(&tc.RowStruct.Field)
				if err != nil {
					t.Error(err)
				}

				if !reflect.DeepEqual(tc.RowStruct.Field, tc.expected) {
					t.Errorf("expected value '%v' got '%v'", tc.expected, tc.RowStruct.Field)
				}
			}
			if d.Err() != nil {
				t.Errorf("d error: %w", err)
			}
		})
	}
}
