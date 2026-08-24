package pkg

import (
	"log"
	"strings"

	tb "github.com/tigerbeetle/tigerbeetle-go"
)

var (
	clusterId        = 0
	addresses string = "0.0.0.0:3000"
	Client    *tb.Client
)

func CreateTigerbeetleClient() *tb.Client {
	if Client != nil {
		return Client
	}
	splitAddresses := strings.Split(addresses, ",")

	newClient, err := tb.NewClient(tb.ToUint128(0), splitAddresses)
	if err != nil {
		log.Panicf("Failed to spin up a tigerbeetle client - %v\n", err)
	}
	Client = &newClient

	log.Println("Successfully created tigerbeetle client")
	return Client
}
