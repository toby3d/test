//go:generate mockgen -package=model -source=./../../internal/model/model.pb.go -destination=./../../internal/model/model_mock.go CartShopClient
package main

import (
	"flag"
	"log"

	"gitlab.com/toby3d/test/internal/client"
	"gitlab.com/toby3d/test/internal/model"
	"golang.org/x/net/context"
)

var flagAddr = flag.String("addr", ":2368", "set specific address and port for client instance")

func main() {
	flag.Parse()

	c, err := client.NewClient(*flagAddr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer c.Close()

	resp, err := c.Add(context.TODO(), &model.AddRequest{ProductId: 5, Quanity: 42})
	if err != nil {
		log.Fatalln(err.Error())
	}
	if !resp.GetOk() {
		log.Printf("Get error on adding product: %s", resp.GetDescription())
		return
	}
	log.Printf(
		"Product %d has been added to cart, current quanity of this product is %d",
		resp.GetItem().GetProductId(), resp.GetItem().GetQuanity(),
	)

	if resp, err = c.Get(context.TODO(), &model.GetRequest{}); err != nil {
		log.Fatalln(err.Error())
	}
	if !resp.GetOk() {
		log.Printf("Get error on getting cart: %s", resp.GetDescription())
		return
	}
	log.Printf(
		"Cart contains %d unique products (in %d quanity) with total price %g",
		resp.GetCart().GetItemsCount(),
		resp.GetCart().GetQuanityCount(),
		resp.GetCart().GetTotalPrice(),
	)
}
