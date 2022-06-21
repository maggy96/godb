# Readme

## what is this?
this is meant to implement and benchmark a JOIN operation on CSV files

```SQL
SELECT s_name, l_shipdate, s_acctbal, l_comment
FROM lineitem l
JOIN supplier s ON l_suppkey = s_suppkey
WHERE s_acctbal < 0
```

## how to run
generate files with [this repo][https://github.com/electrum/tpch-dbgen] and move them to `./assets` then run `go run main.go`
