package main

import (
	"fmt"

	goqu "github.com/doug-martin/goqu/v9"

	// import the dialect
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

func main() {

	goqu.Dialect("postgres")

	// use dialect.From to get a dataset to build your SQL
	ds := goqu.From("test").Where(goqu.Ex{"id": 10})
	sql, args, err := ds.ToSQL()
	if err != nil {
		fmt.Println("An error occurred while generating the SQL", err.Error())
	} else {
		fmt.Println(sql, args)
	}
}
