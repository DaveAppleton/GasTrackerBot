package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

type gasNowData struct {
	Fast   string
	Medium string
	Safe   string
}

type gasNowDecimalData struct {
	Fast   decimal.Decimal `json:"Top50"`
	Medium decimal.Decimal `json:"Top200"`
	Safe   decimal.Decimal `json:"Top400"`
}

func gasNowMessage(data gasNowData, reaped time.Time, advert string) (message string) {

	duration := int(math.Floor(time.Since(reaped).Minutes()))
	return fmt.Sprintf("*Gas prices*\n\n```"+`
Tx Speed Gas needed  
--------+-----------+
Fastest  %5v GWei
Fast     %5v GWei
Safelow  %5v GWei

Data received %d minutes ago
`+"```"+`
/sources to see data sources
%s
`, data.Fast, data.Medium, data.Safe, duration, advert)

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
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(spark)
	if err != nil {
		return
	}
	err = json.NewDecoder(resp.Body).Decode(&reply)
	if err != nil {
		return
	}
	gnd.Safe = fmt.Sprintf("%4d", new(big.Int).Div(reply.Data.Top400, big.NewInt(1000000000)))
	gnd.Medium = fmt.Sprintf("%4d", new(big.Int).Div(reply.Data.Top200, big.NewInt(1000000000)))
	gnd.Fast = fmt.Sprintf("%4d", new(big.Int).Div(reply.Data.Top50, big.NewInt(1000000000)))
	return gnd, nil
}

func getGasNowDecimal() (g gasNowDecimalData, err error) {
	var reply struct {
		Code int
		Data gasNowDecimalData
	}
	spark := "https://www.gasnow.org/api/v1/gas/price"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(spark)
	if err != nil {
		return
	}
	err = json.NewDecoder(resp.Body).Decode(&reply)
	if err != nil {
		return
	}
	dividend := decimal.NewFromInt(1000000000)

	reply.Data.Fast = reply.Data.Fast.Div(dividend)
	reply.Data.Medium = reply.Data.Medium.Div(dividend)
	reply.Data.Safe = reply.Data.Safe.Div(dividend)
	return reply.Data, nil
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
