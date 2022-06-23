package main

import (
	"encoding/csv"
	"fmt"
	util "godb/lib"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type LineItem struct {
	comment  string
	shipdate string
	suppkey  int64
}

type Supplier struct {
	acctbal float64
	name    string
	suppkey int64
}

type JointType struct {
	acctbal  string
	name     string
	comment  string
	shipdate string
}

func readSuppliers(filename string) []Supplier {
	suppliers := make([]Supplier, 59986052)
	i := 0
	util.Readfile(filename, func(line []string) {
		var rec Supplier
		var err error
		for j, field := range line {
			if j == 0 {
				rec.suppkey, err = strconv.ParseInt(field, 10, 32)
				if err != nil {
					log.Fatal(err)
				}
			} else if j == 1 {
				rec.name = strings.TrimSpace(field)
			} else if j == 5 {
				rec.acctbal, err = strconv.ParseFloat(field, 32)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		suppliers[i] = rec
		i = i + 1
	})
	return suppliers
}

func readLineItems(filename string) []LineItem {
	lineItems := make([]LineItem, 59986052)
	i := 0
	util.Readfile(filename, func(line []string) {
		var rec LineItem
		var err error
		for j, field := range line {
			if j == 2 {
				rec.suppkey, err = strconv.ParseInt(field, 10, 32)
				if err != nil {
					log.Fatal(err)
				}
			} else if j == 10 {
				rec.shipdate = strings.TrimSpace(field)
			} else if j == 15 {
				rec.comment = strings.TrimSpace((field))
			}
		}
		lineItems[i] = rec
		i = i + 1
	})
	return lineItems
}

func main() {

	parsingTime := time.Now()

	lineItems := readLineItems("./assets/lineitem.tbl")
	// fmt.Printf("%+v\n", lineItems[0])

	suppliers := readSuppliers("./assets/supplier.tbl")
	// fmt.Printf("%+v\n", suppliers[0])

	duration := time.Since(parsingTime)
	fmt.Printf("time for reading: %fs\n", duration.Seconds())

	join := make([]JointType, len(lineItems))

	joiningTime := time.Now()
	pointer := 0
	for _, lineItem := range lineItems {
		var newRow JointType
		supplier := suppliers[lineItem.suppkey-1]
		if supplier.acctbal >= 0 {
			continue
		}
		newRow.comment = lineItem.comment
		newRow.shipdate = lineItem.shipdate
		newRow.name = supplier.name
		newRow.acctbal = fmt.Sprintf("%.2f", supplier.acctbal)
		join[pointer] = newRow
		pointer++
	}
	duration = time.Since(joiningTime)
	fmt.Printf("time for joining: %fs\n", duration.Seconds())
	fmt.Printf("%d records joined\n", len(suppliers)+len(lineItems))
	writingTime := time.Now()
	WriteFile(join[:pointer], "./assets/out.tbl")
	duration = time.Since(writingTime)
	fmt.Printf("%d lines written\n", pointer)
	fmt.Printf("time for writing: %fs\n", duration.Seconds())
}

func WriteFile(data []JointType, filename string) {
	f, err := os.Create(filename)

	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Comma = '|'
	defer w.Flush()

	for _, record := range data {
		rec := []string{record.name, record.shipdate, record.acctbal, record.comment}
		if err := w.Write(rec); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
}
