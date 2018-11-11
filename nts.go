package main

import (
	"log"
	"net/http"
	"os"

	"github.com/krqa/nts/db"
	"github.com/krqa/nts/routing"
	stripe "github.com/stripe/stripe-go"
)

func main() {
	stripePublicKey := os.Getenv("NTS_STRIPE_PUBLISHABLE_KEY")
	stripe.Key = os.Getenv("NTS_STRIPE_SECRET_KEY")
	if stripePublicKey == "" || stripe.Key == "" {
		log.Fatal("Stripe Key(s) not provided!")
	}

	address := os.Getenv("NTS_ADDRESS")
	if address == "" {
		address = ":9000"
	}

	router := routing.New(db.NewMemTicketsDB(5), stripePublicKey)

	log.Fatal(http.ListenAndServe(address, router))
}
