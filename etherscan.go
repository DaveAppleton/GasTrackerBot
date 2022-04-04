package main

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/DaveAppleton/gasBot/timedhttp"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
)

type etherscanData struct {
	LastBlock       string
	SafeGasPrice    decimal.Decimal
	ProposeGasPrice decimal.Decimal
	FastGasPrice    decimal.Decimal
}

func etherscanGas() (data etherscanData, err error) {
	var reply struct {
		Status  string
		Message string
		Result  struct {
			LastBlock       string
			SafeGasPrice    string
			ProposeGasPrice string
			FastGasPrice    string
		}
	}
	resp, err := timedhttp.Get("https://api.etherscan.io/api?module=gastracker&action=gasoracle&apikey=" + viper.GetString("ETHERSCAN_KEY"))
	if err != nil {
		log.Println("etherscan", err)
		return
	}
	if err = json.NewDecoder(resp.Body).Decode(&reply); err != nil {
		log.Println("etherscan", err)
		return
	}
	if reply.Status != "1" {
		err = errors.New(reply.Message)
		return
	}
	data.LastBlock = reply.Result.LastBlock
	data.SafeGasPrice, err = decimal.NewFromString(reply.Result.SafeGasPrice)
	if err != nil {
		return
	}
	data.ProposeGasPrice, err = decimal.NewFromString(reply.Result.ProposeGasPrice)
	if err != nil {
		return
	}
	data.FastGasPrice, err = decimal.NewFromString(reply.Result.FastGasPrice)
	return
}
