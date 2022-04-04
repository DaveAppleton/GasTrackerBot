package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/DaveAppleton/etherUtils"
	"github.com/DaveAppleton/gasBot/timedhttp"
)

type gasNowData struct {
	Fast   string
	Medium string
	Safe   string
}

func getGasNow() (g gasNowData, err error) {
	var reply struct {
		Code int
		Data struct {
			Top50     *big.Int
			Top200    *big.Int
			Top400    *big.Int
			timestamp int64
		}
	}
	var gnd gasNowData
	spark := "https://www.gasnow.org/api/v1/gas/price"
	resp, err := timedhttp.Get(spark)
	if err != nil {
		return
	}
	err = json.NewDecoder(resp.Body).Decode(&reply)
	if err != nil {
		return
	}
	gnd.Safe = etherUtils.CoinToStr(reply.Data.Top400, 9)
	gnd.Medium = etherUtils.CoinToStr(reply.Data.Top200, 9)
	gnd.Fast = etherUtils.CoinToStr(reply.Data.Top50, 9)
	return gnd, nil
}

func gasNowLoop(ch chan gasNowData, exit chan bool) {
	gnd, err := getGasNow()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("returning data")
	ch <- gnd
	for {
		select {
		case <-exit:
			return
		case <-time.After(5 * time.Minute):
			gnd, err = getGasNow()
			if err != nil {
				log.Println("getGasNow", err)
			} else {
				fmt.Println("returning data")
				ch <- gnd
			}
		}
	}
}
