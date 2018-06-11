package main

import (
	"log"
	"time"

	"github.com/meeDamian/satoshis-bch/invoice"
)

func main() {
	time.Sleep(6 * time.Second)

	pic := invoice.Picture{
		invoice.Pixel{Coords: [2]string{"1", "1"}, Color: "#ffffff"},
		invoice.Pixel{Coords: [2]string{"0", "0"}, Color: "#ffffff"},
		invoice.Pixel{Coords: [2]string{"1", "0"}, Color: "#222222"},
		invoice.Pixel{Coords: [2]string{"0", "1"}, Color: "#222222"},
	}

	inv, err := invoice.GetInvoice(pic)
	if err !=nil {
		log.Println("unable to get an invoice")
		return
	}

	log.Println("invoice:", inv)

	time.Sleep(10 * time.Second)
}
