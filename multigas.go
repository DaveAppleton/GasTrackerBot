package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
)

type multiData struct {
	EGS       GasData
	GN        gasNowDecimalData
	ES        etherscanData
	Timestamp string
	Errors    string
	Advert    string
	Version   string
}

type sortableDecimals []decimal.Decimal

func (a sortableDecimals) Len() int           { return len(a) }
func (a sortableDecimals) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortableDecimals) Less(i, j int) bool { return a[i].LessThan(a[j]) }

func multiGas() (output string, median decimal.Decimal, data GasData, err error) {
	safeLowz := sortableDecimals{}
	fastz := sortableDecimals{}
	fastestz := sortableDecimals{}
	Errors := ""
	egsData, err := getGasData()
	blocknum := uint64(0)
	if err == nil {
		ten := decimal.NewFromFloat(10.0)
		egsData.Fastest = egsData.Fastest.Div(ten)
		egsData.Fast = egsData.Fast.Div(ten)
		egsData.SafeLow = egsData.SafeLow.Div(ten)
		egsData.Average = egsData.Average.Div(ten)
		safeLowz = append(safeLowz, egsData.SafeLow)
		fastz = append(fastz, egsData.Fast)
		fastestz = append(fastestz, egsData.Fastest)
		blocknum = egsData.BlockNum
	} else {
		Errors = "error getting results from ETHGasStation\n"
	}
	gasNowData, err := getGasNowDecimal()
	if err != nil {
		Errors += "Error getting results from GasNow\n"
	} else {
		safeLowz = append(safeLowz, gasNowData.Safe)
		fastz = append(fastz, gasNowData.Medium)
		fastestz = append(fastestz, gasNowData.Fast)
	}

	esData, err := etherscanGas()
	if err != nil {
		Errors += "Error getting results from Etherscan"
	} else {
		safeLowz = append(safeLowz, esData.SafeGasPrice)
		fastz = append(fastz, esData.ProposeGasPrice)
		fastestz = append(fastestz, esData.FastGasPrice)
	}
	switch len(safeLowz) {
	case 0:
		err = errors.New(Errors)
		return
	case 1:
		median = safeLowz[0]
		data = GasData{SafeLow: safeLowz[0], Fast: fastz[0], Fastest: fastestz[0], BlockNum: blocknum}
	case 2:
		avSL := safeLowz[0].Add(safeLowz[1]).Div(decimal.NewFromInt(2))
		avF := fastz[0].Add(fastz[1]).Div(decimal.NewFromInt(2))
		avFF := fastestz[0].Add(fastestz[1]).Div(decimal.NewFromInt(2))
		median = (avSL)
		data = GasData{SafeLow: avSL, Fast: avF, Fastest: avFF, BlockNum: blocknum}
	case 3:
		sort.Sort(safeLowz)
		sort.Sort(fastz)
		sort.Sort(fastestz)
		median = safeLowz[1]
		data = GasData{SafeLow: safeLowz[1], Fast: fastz[1], Fastest: fastestz[1], BlockNum: blocknum}
	}
	advert := viper.GetString("ADVERT")
	advertBytes, err := ioutil.ReadFile("files/advert.md")
	if err == nil {
		advert = string(advertBytes)
	}
	index, err := template.ParseFiles("files/index.md")
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	versionStr := fmt.Sprintf("%4.2f", viper.GetFloat64("Version"))
	versionStr = strings.Replace(versionStr, ".", "\\.", 1)
	if err = index.Execute(buf, multiData{EGS: egsData, GN: gasNowData, ES: esData, Timestamp: time.Now().String(), Errors: Errors, Advert: advert, Version: versionStr}); err != nil {
		return
	}
	output = buf.String()
	return
}
