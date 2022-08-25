# go-electrum [![GoDoc](https://godoc.org/github.com/checksum0/go-electrum?status.svg)](https://godoc.org/github.com/checksum0/go-electrum)
A pure Go [Electrum](https://electrum.org/) bitcoin library supporting the latest [ElectrumX](https://github.com/kyuupichan/electrumx) protocol versions.  
This makes it easy to write cryptocurrencies based services in a trustless fashion using Go without having to run a full node.

![go-electrum](https://raw.githubusercontent.com/checksum0/go-electrum/master/media/logo.png)

## Usage
See [example/](https://github.com/checksum0/go-electrum/tree/master/example) for more.

### electrum [![GoDoc](https://godoc.org/github.com/checksum0/go-electrum/electrum?status.svg)](https://godoc.org/github.com/checksum0/go-electrum/electrum)
```bash
$ go get github.com/checksum0/go-electrum
```

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/checksum0/go-electrum/electrum"
)

func main() {
	// Establishing a new SSL connection to an ElectrumX server
	client := electrum.NewClient()
	if err := client.ConnectTCP(context.TODO(), "bch.imaginary.cash:50001"); err != nil {
		log.Fatal(err)
	}
    ctx := context.TODO()
	// Making sure connection is not closed with timed "client.ping" call
	go func() {
		for {
			if err := client.Ping(ctx); err != nil {
				log.Fatal(err)
			}
			time.Sleep(60 * time.Second)
		}
	}()

	// Making sure we declare to the server what protocol we want to use
	if _, _, err := client.ServerVersion(ctx); err != nil {
		log.Fatal(err)
	}

	// Asking the server for the balance of address 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa
	// 8b01df4e368ea28f8dc0423bcf7a4923e3a12d307c875e47a0cfbf90b5c39161
	// We must use scripthash of the address now as explained in ElectrumX docs
	scripthash, _ := electrum.AddressToElectrumScriptHash("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa")
	balance, err := client.GetBalance(ctx, scripthash)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Address confirmed balance:   %+v", balance.Confirmed)
	log.Printf("Address unconfirmed balance: %+v", balance.Unconfirmed)
}
```

# License
go-electrum is licensed under the MIT license. See LICENSE file for more details.

Copyright (c) 2022 Roman Maklakov  
Copyright (c) 2019 Ian Descôteaux  
Copyright (c) 2015 Tristan Rice

Based on Tristan Rice [go-electrum](https://github.com/d4l3k/go-electrum) unmaintained library.
