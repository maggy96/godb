package lib

import (
	"encoding/csv"
	"log"
	"os"
)

func Readfile(s string) [][]string {
	file, err := os.Open(s)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	fileReader := csv.NewReader(file)
	fileReader.Comma = '|'
	data, err := fileReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return data
}
