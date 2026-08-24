package services

import (
	"log"
	"main/pkg"

	tb "github.com/tigerbeetle/tigerbeetle-go"
)

func FetchAllAccounts() *[]tb.Account {
	client := pkg.CreateTigerbeetleClient()

	res, err := (*client).QueryAccounts(tb.QueryFilter{})

	if err != nil {
		log.Fatalf("Failed to fetch accounts - %v\n", err)
	}

	return &res
}
