package lib

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/edsrzf/mmap-go"
	"io"
	"log"
	"os"
	"time"
)

type linereader func([]string)

/*
  -
  - nr lines in the file
  - divide into chunks
*/
func Readfile(s string, parser linereader) {
	readingTime := time.Now()
	file, err := os.Open(s)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// mmap file
	mmap, _ := mmap.Map(file, mmap.RDONLY, 0)
	defer mmap.Unmap()
	fileReader := bytes.NewReader(mmap)

	// read csv
	csvReader := csv.NewReader(fileReader)
	csvReader.Comma = '|'
	csvReader.ReuseRecord = true
	// data, err := csvReader.ReadAll()

	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		parser(rec)
	}

	fmt.Printf("time for reading file %s: %fs\n", s, time.Since(readingTime).Seconds())
}
