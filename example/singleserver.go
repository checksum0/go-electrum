package main

import (
	"context"
	"log"
	"time"

	"github.com/checksum0/go-electrum/electrum"
)

func main() {
	client := electrum.NewClient()

	if err := client.ConnectTCP(context.Background(), "bch.imaginary.cash:50001"); err != nil {
		log.Fatal(err)
	}

	serverVer, protocolVer, err := client.ServerVersion(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Server version: %s [Protocol %s]", serverVer, protocolVer)

	go func() {
		for {
			if err := client.Ping(context.Background()); err != nil {
				log.Fatal(err)
			}
			time.Sleep(60 * time.Second)
		}
	}()
}
