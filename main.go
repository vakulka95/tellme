package main

import (
	"log"

	"gitlab.com/tellmecomua/tellme.api/pkg/postgres"
)

func main() {
	sql, args := postgres.NewQueryBuilder().
		Select("id", "id_2").
		From("users").
		Where(
			postgres.NewNotEqual("status", "active"),
		).
		Build()

	log.Printf("sql: [%s]", sql)
	log.Printf("args: %+v", args)
}
