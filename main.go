package main

import (
	"context"
	"fmt"
	"github.com/phillipahereza/momoapi-go/momo"
	"log"
)

func main() {
	collectionPK := "0d31d966e5674a999c82772aa95f2cca"
	userID := "356c0142-21a5-4ff3-9f00-87c268a60378"
	apiKey := "c0d857dba3944ce3b6d436c04963e1ea"

	ctx := context.Background()

	client := momo.NewClient(collectionPK, "sandbox", "https://sandbox.momodeveloper.mtn.com/")
	_, err := client.Collection.GetToken(ctx, apiKey, userID)
	if err != nil {
		log.Fatal(err)
	}

	transactionID, err := client.Collection.RequestToPay(ctx, "46733123453", 500, "2323", "", "", "EUR")
	if err != nil {
		log.Fatal(err)
	}

	status, err := client.Collection.GetTransaction(ctx, transactionID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", status)

	balance, err := client.Collection.GetBalance(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", balance)

	active, err := client.Collection.IsPayeeActive(ctx, "256789997290")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(active)
}