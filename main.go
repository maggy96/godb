package main

import (
	"encoding/csv"
	"fmt"
	util "godb/lib"
	"log"
	"os"
	"strconv"
	"strings"
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

func processSuppliers(data [][]string) []Supplier {
	suppliers := make([]Supplier, len(data))
	for i, line := range data {
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
	}

	return suppliers
}

func processLineItems(data [][]string) []LineItem {
	lineItems := make([]LineItem, len(data))
	for i, line := range data {
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
	}

	return lineItems
}

func main() {
	lineItemRaw := util.Readfile("./assets/lineitem.tbl")
	supplierRaw := util.Readfile("./assets/supplier.tbl")

	lineItems := processLineItems(lineItemRaw)
	fmt.Printf("%+v\n", lineItems[0])

	suppliers := processSuppliers(supplierRaw)
	fmt.Printf("%+v\n", suppliers[0])

	supplierMap := make(map[int64]int)

	for i, supplier := range suppliers {
		supplierMap[supplier.suppkey] = i
	}

	join := make([]JointType, len(lineItems))

	pointer := 0
	for _, lineItem := range lineItems {
		var newRow JointType
		supplier := suppliers[supplierMap[lineItem.suppkey]]
		if supplier.acctbal < 0 {
			newRow.comment = lineItem.comment
			newRow.shipdate = lineItem.shipdate
			newRow.name = supplier.name
			newRow.acctbal = fmt.Sprintf("%.2f", supplier.acctbal)
			join[pointer] = newRow
			pointer++
		}
	}
	fmt.Printf("%d records joined\n", len(suppliers)+len(lineItems))
	WriteFile(join[:pointer], "./assets/out.tbl")
	fmt.Printf("%d lines written\n", pointer)
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
