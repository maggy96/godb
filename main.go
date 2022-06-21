package main

import (
	"encoding/csv"
	"fmt"
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
	f_lineitem, err := os.Open("./assets/lineitem.tbl")
	if err != nil {
		log.Fatal(err)
	}

	f_supplier, err := os.Open("./assets/supplier.tbl")
	if err != nil {
		log.Fatal(err)
	}

	defer f_lineitem.Close()
	defer f_supplier.Close()

	lineItemReader := csv.NewReader(f_lineitem)
	supplierReader := csv.NewReader(f_supplier)
	lineItemReader.Comma = '|'
	supplierReader.Comma = '|'

	data, err := lineItemReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	lineItems := processLineItems(data)
	fmt.Printf("%+v\n", lineItems[0])

	data, errS := supplierReader.ReadAll()
	if errS != nil {
		log.Fatal(err)
	}

	suppliers := processSuppliers(data)

	fmt.Printf("%+v\n", suppliers[0])
}
