package csvdecoder_test

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
