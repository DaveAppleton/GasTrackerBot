package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/DaveAppleton/etherUtils"
	"github.com/DaveAppleton/gasBot/timedhttp"
)

/*
{
	"code": 200,
	"data": {
		"top50": 132109663366,
		"top200": 117000000000,
		"top400": 106405024781,
		"timestamp": 1597515192378
	}
}
*/

func TestSpark(t *testing.T) {
	initViper()
	var reply struct {
		Code int
		Data struct {
			Top50     *big.Int
			Top200    *big.Int
			Top400    *big.Int
			timestamp int64
		}
	}
	spark := "https://www.gasnow.org/api/v1/gas/price"
	resp, err := timedhttp.Get(spark)
	if err != nil {
		t.Fatal(err)
	}
	err = json.NewDecoder(resp.Body).Decode(&reply)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(
		etherUtils.CoinToStr(reply.Data.Top400, 9),
		etherUtils.CoinToStr(reply.Data.Top200, 9),
		etherUtils.CoinToStr(reply.Data.Top50, 9),
	)
	t.Fail()
}

func TestGetGasNow(t *testing.T) {
	initViper()
	gnd, err := getGasNow()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(gnd.Fast)
	t.Log(gnd.Medium)
	t.Log(gnd.Safe)
	t.Fail()
}

func testLoopData(t *testing.T, g chan gasNowData, b chan bool, c chan bool) {
	for {
		select {
		case gd := <-g:
			t.Log(gd.Fast, gd.Medium, gd.Safe)
			continue
		case <-time.After(30 * time.Second):
			b <- true
			c <- true
			return
		}
	}
}

func TestGasNowLoop(t *testing.T) {
	gndChan := make(chan gasNowData, 10)
	forceExitChan := make(chan bool, 2)
	cChan := make(chan bool, 2)
	go gasNowLoop(gndChan, forceExitChan)
	go testLoopData(t, gndChan, forceExitChan, cChan)
	<-cChan

}
