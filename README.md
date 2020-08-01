# csvdecoder

csvdecoder is a Go library for parsing and decoding csv files into Go objects. It uses [encoding/csv](https://golang.org/pkg/encoding/csv/) for the parsing and follows a usage pattern inspired by the [database/sql](https://golang.org/pkg/database/sql/) package.

The scanning method is *not multi-thread safe*; `Next` and `Scan` are not expected to be called concurrently.
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
- a slice of values. Note that the CSV record must be a valid JSON array. If not a JSON array, a custom decoder implementing the `csvdecoder.Interface` interface must be implemented.
- an array of values. Note that the CSV record must be a valid JSON array. If not a JSON array, a custom decoder implementing the `csvdecoder.Interface` interface must be implemented.
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

The behaviour of the decoder can be configured by passing one of following options when creating the decoder:
- Comma: the character that separates values. Default value is comma.
- IgnoreHeaders: if set to true, the first line will be ignored
- IgnoreUnmatchingFields: if set to true, the number of records and scan targets are allowed to be different. By default, if they don't match exactly it will cause an error.

```golang
	decoder, err := csvdecoder.NewWithConfig(file, csvdecoder.Config{Comma: ';', IgnoreHeaders: true})
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as needed.

## License
[MIT](https://choosealicense.com/licenses/mit/)
