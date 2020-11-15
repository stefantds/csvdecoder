# csvdecoder

csvdecoder is a Go library for parsing and deserializing csv files into Go objects.
It relies on [encoding/csv](https://golang.org/pkg/encoding/csv/) for the actual parsing of the CSV file and follows a similar usage pattern as the [database/sql](https://golang.org/pkg/database/sql/) package for scanning rows.

csvdecoder allows to iterate through the CSV records (using 'Next') and scan the fields into target variables or fields of variables (using 'Scan').

The scanning method is **not thread safe**; `Next` and `Scan` are not expected to be called concurrently.

## Installation

```bash
go get github.com/stefantds/csvdecoder
```

## Supported formats

Csvdecoder supports converting columns read from the source file into the following types:
- `*string`
- `*int`, `*int8`, `*int16`, `*int32`, `*int64`
- `*uint`, `*uint8`, `*uint16`, `*uint32`, `*uint64`
- `*bool`
- `*float32`, `*float64`
- a slice of values. Note that the CSV field must be a valid JSON array. If not a JSON array, a custom decoder implementing the `csvdecoder.Interface` interface must be implemented.
- an array of values. Note that the CSV field must be a valid JSON array. If not a JSON array, a custom decoder implementing the `csvdecoder.Interface` interface must be implemented.
- a pointer to any type implementing the `csvdecoder.Interface` interface

## Usage

```golang
import (
	"fmt"
	"os"

	"github.com/stefantds/csvdecoder"
)

type User struct {
	Name   string
	Active bool
	Age    int
}

func Example_simple() {
	// the csv file contains the values:
	//john,44,true
	//lucy,48,false
	//mr hyde,34,true
	file, err := os.Open("./data/simple.csv")
	if err != nil {
		// handle error
		return
	}
	defer file.Close()

	// create a new decoder that will read from the given file
	decoder, err := csvdecoder.New(file)
	if err != nil {
		// handle error
		return
	}

	// iterate over the rows in the file
	for decoder.Next() {
		var u User

		// scan the first three values in the name, age and active fields respectively
		if err := decoder.Scan(&u.Name, &u.Age, &u.Active); err != nil {
			// handle error
			return
		}
		fmt.Println(u)
	}

	// check if the loop stopped prematurely because of an error
	if err = decoder.Err(); err != nil {
		// handle error
		return
	}

	// Output: {john true 44}
	// {lucy false 48}
	// {mr hyde true 34}
}
```

See also the example files for more usage examples.

## Configuration

The behavior of the decoder can be configured by passing one of following options when creating the decoder:
- Comma: the character that separates values. Default value is comma.
- IgnoreHeaders: if set to true, the first line will be ignored. This is useful when the CSV file contains a header line.
- IgnoreUnmatchingFields: if set to true, the number of fields and scan targets are allowed to be different. By default, if they don't match exactly it will cause an error.
- EscapeChar: the character used to escape the quote character in quoted fields. The default is the quote itself as used by the `encoding/csv` reader.

```golang
	decoder, err := csvdecoder.NewWithConfig(file, csvdecoder.Config{Comma: ';', IgnoreHeaders: true})
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as needed.

## License

[MIT](https://choosealicense.com/licenses/mit/)
