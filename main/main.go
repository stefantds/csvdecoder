package main

import (
	"fmt"
	"os"
	"strings"

	csvdecoder "github.com/stefantds/csvdecoder"
)

type TestExampleRow struct {
	Name            string
	Age             int
	Gender          string
	FavoriteNumbers []int
	FavoriteHeroes  []string
}

type ExampleStringArrayRow struct {
	Heroes []string
}

type ExampleSimpleRow struct {
	Name   string
	Age    int
	Gender string
	Real   bool
}

type ExampleIntArrayRow struct {
	Heroes []int
}

type ExampleFloatArrayRow struct {
	Heroes []float64
}

type ExampleIntArrayArrayRow struct {
	Heroes [][]int
}

type MyDeserializer1 struct {
	value string
}

func (d *MyDeserializer1) DecodeRecord(value string) error {
	d.value = strings.Join(strings.Split(value, "e"), ".")
	return nil
}

type MyDeserializer2 struct {
	value string
}

func (d *MyDeserializer2) DecodeRecord(value string) error {
	d.value = strings.Join(strings.Split(value, "e"), ".")
	return nil
}

type TestExampleDeserializerRow struct {
	Name            *MyDeserializer1
	Age             *MyDeserializer2
	Gender          csvdecoder.Decoder
	FavoriteNumbers []int
	FavoriteHeroes  []string
}

func main() {
	fmt.Println("running")
	file, err := os.Open("/Users/stefan.tudose/private/gitrepos/csvdecoder/example_simple.csvdecoder")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	parser, err := csvdecoder.NewParser(file)
	if err != nil {
		panic(err)
	}

	for {
		data := ExampleSimpleRow{}
		eof, err := parser.Next(&data)
		if eof {
			return
		}
		if err != nil {
			panic(err)
		}
		fmt.Println(data)
	}
}
