package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/DaveAppleton/gasBot/timedhttp"
	gecko "github.com/superoo7/go-gecko/v3"
)

type oneInchData struct {
	FromToken struct {
		Decimals int64
		Symbol   string
		Name     string
	}
	ToToken struct {
		Decimals int64
		Symbol   string
		Name     string
	}
	FromTokenAmount string `json:"fromTokenAmount"`
	ToTokenAmount   string `json:"toTokenAmount"`
	Message         string
	OK              bool
}

func getTokenPrice(token string) (reply oneInchData) {
	cg := gecko.NewClient(nil)
	price, err := cg.SimpleSinglePrice(token, "eth")
	if err != nil {
		log.Println(err)
		return
	}
	reply.Message = fmt.Sprintf("1 ETH buys %4.2f %s", 1/price.MarketPrice, token)
	reply.OK = true
	reply.ToTokenAmount = fmt.Sprintf("%4.2f", 1/price.MarketPrice)

	return
}

func getOneInchPrice(token string) (reply oneInchData) {
	call := "https://api.1inch.exchange/v1.1/quote?fromTokenSymbol=ETH&toTokenSymbol=%s&amount=1000000000000000000"
	ethURL := fmt.Sprintf(call, token)
	resp, err := timedhttp.Get(ethURL)
	if err == nil {
		err = json.NewDecoder(resp.Body).Decode(&reply)
		if err == nil {
			reply.OK = true

			decimals := int(reply.ToToken.Decimals)
			if err != nil {
				log.Println(err)
				return oneInchData{OK: false}
			}
			if decimals != 0 {
				if len(reply.ToTokenAmount) > decimals {
					fmt.Println("case 1")
					pos := len(reply.ToTokenAmount) - decimals
					reply.ToTokenAmount = reply.ToTokenAmount[0:pos] + "." + reply.ToTokenAmount[pos:]
				} else if len(reply.ToTokenAmount) == decimals {
					fmt.Println("case 2")
					reply.ToTokenAmount = "0." + reply.ToTokenAmount
				} else {
					fmt.Println("case 3")
					pos := len(reply.ToTokenAmount) - decimals
					add := "000000000000000000000000000000000000"
					add = add[len(reply.ToTokenAmount)+pos:]
					reply.ToTokenAmount = "0." + add + reply.ToTokenAmount
				}
			}
			reply.Message = fmt.Sprintf("1 ETH buys %s `%s`", reply.ToTokenAmount, reply.ToToken.Name)
		} else {
			fmt.Println("failed to decode", err)
		}
	} else {
		fmt.Println("failed to receive reply", err)
	}
	return
}
