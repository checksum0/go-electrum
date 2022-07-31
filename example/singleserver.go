package main

import (
	"log"
	"time"

	"github.com/checksum0/go-electrum/electrum"
)

func main() {
	client := electrum.NewClient(&electrum.ClientOptions{
		ConnTimeout: time.Second * 10,
		ReqTimeout:  time.Second * 10,
	})

	if err := client.ConnectTCP("bch.imaginary.cash:50001"); err != nil {
		log.Fatal(err)
	}

	serverVer, protocolVer, err := client.ServerVersion()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Server version: %s [Protocol %s]", serverVer, protocolVer)

	go func() {
		for {
			if err := client.Ping(); err != nil {
				log.Fatal(err)
			}
			time.Sleep(60 * time.Second)
		}
	}()
}
