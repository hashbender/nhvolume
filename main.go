package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Algo struct {
	Index int
	Name  string
	Norm  float64
}

type OrderGetResp struct {
	Result ResultResp `json:"result"`
	Method string     `json:"method"`
}

type ResultResp struct {
	Orders []OrderRep `json:"orders"`
}

type OrderRep struct {
	Alive bool    `json:"alive"`
	Price float64 `json:"price,string"`
	Speed float64 `json:"accepted_speed,string"`
}

func main() {
	var algos = []*Algo{
		&Algo{
			Index: 0,
			Name:  "Scrypt",
			Norm:  .001,
		},
		&Algo{
			Index: 1,
			Name:  "SHA256",
			Norm:  .000001,
		},
		&Algo{
			Index: 2,
			Name:  "ScryptNf",
			Norm:  0,
		},
		&Algo{
			Index: 3,
			Name:  "X11",
			Norm:  .001,
		},
		&Algo{
			Index: 4,
			Name:  "X13",
			Norm:  .001,
		},
		&Algo{
			Index: 5,
			Name:  "Keccak",
			Norm:  .001,
		},
		&Algo{
			Index: 6,
			Name:  "X15",
			Norm:  1,
		},
		&Algo{
			Index: 7,
			Name:  "Nist5",
			Norm:  .001,
		},
		&Algo{
			Index: 8,
			Name:  "NeoScrypt",
			Norm:  1,
		},
		&Algo{
			Index: 9,
			Name:  "Lyra2RE",
			Norm:  1,
		},
		&Algo{
			Index: 10,
			Name:  "WhirlpoolX",
			Norm:  0,
		},
		&Algo{
			Index: 11,
			Name:  "Qubit",
			Norm:  .001,
		},
		&Algo{
			Index: 12,
			Name:  "Quark",
			Norm:  .001,
		},
		&Algo{
			Index: 13,
			Name:  "Axiom",
			Norm:  0,
		},
		&Algo{
			Index: 14,
			Name:  "Lyra2REv2",
			Norm:  .001,
		},
		&Algo{
			Index: 15,
			Name:  "ScryptJaneNf16",
			Norm:  0,
		},
		&Algo{
			Index: 16,
			Name:  "Blake256r8",
			Norm:  0,
		},
		&Algo{
			Index: 17,
			Name:  "Blake256r14",
			Norm:  0,
		},
		&Algo{
			Index: 18,
			Name:  "Blake256r8vnl",
			Norm:  0,
		},
		&Algo{
			Index: 19,
			Name:  "Hodl",
			Norm:  0,
		},
		&Algo{
			Index: 20,
			Name:  "DaggerHashimoto",
			Norm:  1,
		},
		&Algo{
			Index: 21,
			Name:  "Decred",
			Norm:  .001,
		},
		&Algo{
			Index: 22,
			Name:  "CryptoNight",
			Norm:  1000,
		},
		&Algo{
			Index: 23,
			Name:  "Lbry",
			Norm:  .001,
		},
		&Algo{
			Index: 24,
			Name:  "Equihash",
			Norm:  1000,
		},
		&Algo{
			Index: 25,
			Name:  "Pascal",
			Norm:  .001,
		},
		&Algo{
			Index: 26,
			Name:  "X11Gost",
			Norm:  0,
		},
		&Algo{
			Index: 27,
			Name:  "Sia",
			Norm:  .001,
		},
		&Algo{
			Index: 28,
			Name:  "Blake2s",
			Norm:  .001,
		},
		&Algo{
			Index: 29,
			Name:  "Skunk",
			Norm:  0,
		},
		&Algo{
			Index: 30,
			Name:  "CryptoNightV7",
			Norm:  1000,
		},
		&Algo{
			Index: 31,
			Name:  "CryptoNightHeavy",
			Norm:  1000,
		},
		&Algo{
			Index: 32,
			Name:  "Lyra2Z",
			Norm:  1,
		},
		&Algo{
			Index: 33,
			Name:  "X16R",
			Norm:  1,
		},
	}

	total := 0.0
	for _, algo := range algos {
		resp, err := http.Get(fmt.Sprintf("https://api.nicehash.com/api?method=orders.get&location=0&algo=%d", algo.Index))
		if err != nil {
			log.Printf("Error getting %s.  Error: %s", algo.Name, err)
			continue
		}
		var resultRep OrderGetResp
		err = json.NewDecoder(resp.Body).Decode(&resultRep)
		if err != nil {
			log.Printf("Error parsing: %s", err)
			continue
		}
		sum := 0.0
		for _, order := range resultRep.Result.Orders {
			if order.Alive {
				sum += order.Speed * algo.Norm * order.Price
			}
		}
		total += sum
		log.Printf("EU Algo: %s. Total BTC: %f.  Num Orders: %d", algo.Name, sum, len(resultRep.Result.Orders))
	}

	for _, algo := range algos {
		resp, err := http.Get(fmt.Sprintf("https://api.nicehash.com/api?method=orders.get&location=1&algo=%d", algo.Index))
		if err != nil {
			log.Printf("Error getting %s.  Error: %s", algo.Name, err)
			continue
		}
		var resultRep OrderGetResp
		err = json.NewDecoder(resp.Body).Decode(&resultRep)
		if err != nil {
			log.Printf("Error parsing: %s", err)
			continue
		}
		sum := 0.0
		for _, order := range resultRep.Result.Orders {
			if order.Alive {
				sum += order.Speed * algo.Norm * order.Price
			}
		}
		total += sum
		log.Printf("US Algo: %s. Total BTC: %f.  Num Orders: %d", algo.Name, sum, len(resultRep.Result.Orders))
	}

	fmt.Printf("Total BTC: %f", total)
}
