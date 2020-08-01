package csvdecoder_test

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/stefantds/csvdecoder"
)

type Point struct {
	X int
	Y int
}

// DecodeRecord implements the csvdecoder.Interface type
func (p *Point) DecodeRecord(record string) error {
	// the decode code is specific to the way the object is serialized.
	// in this example the point is encoded as a JSON array with two values

	data := make([]int, 2)
	if err := json.NewDecoder(strings.NewReader(record)).Decode(&data); err != nil {
		return fmt.Errorf("could not parse %s as JSON array: %w", record, err)
	}

	(*p).X = data[0]
	(*p).Y = data[1]

	return nil
}

func Example_custom_decoder() {
	// the csv separator is a semicolon in this example
	exampleData := strings.NewReader(
		`[0, 0];[0, 2];[1, 2]
[-1, 2];[0, -2];[1, 0]
`)

	// create a new decoder that will read from the given file
	decoder, err := csvdecoder.NewWithConfig(exampleData, csvdecoder.Config{Comma: ';'})
	if err != nil {
		// handle error
		return
	}

	// iterate over the rows in the file
	for decoder.Next() {
		var a, b, c Point

		// scan the first values to the types
		if err := decoder.Scan(&a, &b, &c); err != nil {
			// handle error
			return
		}
		fmt.Printf("a: %v, b: %v, c: %v\n", a, b, c)
	}

	// check if the loop stopped prematurely because of an error
	if err = decoder.Err(); err != nil {
		// handle error
		return
	}

	// Output: a: {0 0}, b: {0 2}, c: {1 2}
	// a: {-1 2}, b: {0 -2}, c: {1 0}
}
