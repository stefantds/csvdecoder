package csvdecoder_test

import (
	"fmt"
	"strings"

	"github.com/stefantds/csvdecoder"
)

type MyStringCollection []string

// DecodeRecord implements the csvdecoder.Interface type
func (c *MyStringCollection) DecodeRecord(record string) error {
	// the decode code is specific to the way the value is serialized.
	// in this example the array is represented as int values separated by space
	*c = MyStringCollection(strings.Split(record, " "))

	return nil
}

func Example_slices() {
	// the csv separator is a semicolon in this example
	// the values are arrays serialized in two different ways.
	exampleData := strings.NewReader(
		`jon;elvis boris ahmed jane;["jo", "j"]
jane;lucas george;["j", "jay"]
`)

	// create a new decoder that will read from the given file
	decoder, err := csvdecoder.NewWithConfig(exampleData, csvdecoder.Config{Comma: ';'})
	if err != nil {
		// handle error
		return
	}

	type Person struct {
		Name      string
		Friends   MyStringCollection
		Nicknames []string
	}

	// iterate over the rows in the file
	for decoder.Next() {
		var p Person

		// scan the first values to the types
		if err := decoder.Scan(&p.Name, &p.Friends, &p.Nicknames); err != nil {
			// handle error
			return
		}
		fmt.Printf("%v\n", p)
	}

	// check if the loop stopped prematurely because of an error
	if err = decoder.Err(); err != nil {
		// handle error
		return
	}

	// Output: {jon [elvis boris ahmed jane] [jo j]}
	// {jane [lucas george] [j jay]}
}
