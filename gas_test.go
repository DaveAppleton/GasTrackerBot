package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
)

func TestSaveUser(t *testing.T) {
	initViper()
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	exists, err := saveUserToDatabase("dave", "appleton", "DaveAppleton", 42)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(exists)
	t.Fail()
}

func TestSaveLevel(t *testing.T) {
	initViper()
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	five := decimal.NewFromFloat(5.0)
	err = saveLevelToDatabase(42, five)
	if err != nil {
		log.Fatal(err)
	}
	fivepointfive := decimal.NewFromFloat(5.5)
	err = saveLevelToDatabase(42, fivepointfive)
	if err != nil {
		log.Fatal(err)
	}

}

func TestStatus(t *testing.T) {
	initViper()
	levelState, userState, levelValue, err := getUserInfo(37262138)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(levelState, userState, levelValue)
	t.Fail()
}

func TestVersion(t *testing.T) {
	initViper()
	fmt.Println(viper.GetFloat64("Version"))
	t.Fail()
}

func TestSetShow(t *testing.T) {
	initViper()
	err := setOne(37262138, "show_level", true)
	if err != nil {
		t.Fatal(err)
	}
	err = setOne(37262138, "show_info", true)
	if err != nil {
		t.Fatal(err)
	}
	err = setOne(37262138, "version", 0.05)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetBelow(t *testing.T) {
	initViper()
	list, err := getAllUsersBelowVersion(0.04)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 0 {
		for _, num := range list {
			t.Log(num)
		}
		t.Fatal(len(list), "below 0.04")
	}
	list, err = getAllUsersBelowVersion(0.06)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		for _, num := range list {
			t.Log(num)
		}
		t.Fatal(len(list), "below 0.06")
	}
}

func TestNewVersion(t *testing.T) {
	initViper()
	version := viper.GetFloat64("Version")
	t.Log(version)
	v := fmt.Sprintf("NewIn%v", version)
	t.Fatal(v)
}

func TestGasLevel(t *testing.T) {
	initViper()
	data, err := getGasData()
	if err != nil {
		t.Fatal("error getting gas prices" + err.Error())
	}
	ten := decimal.NewFromFloat(10.0)
	data.Fastest = data.Fastest.Div(ten)
	data.Fast = data.Fast.Div(ten)
	data.SafeLow = data.SafeLow.Div(ten)
	data.Average = data.Average.Div(ten)
	t.Log(data)
	t.Fail()
}

func TestSaveGasLevel(t *testing.T) {
	initViper()
	data, err := getGasData()
	if err != nil {
		t.Fatal("error getting gas prices" + err.Error())
	}
	ten := decimal.NewFromFloat(10.0)
	data.Fastest = data.Fastest.Div(ten)
	data.Fast = data.Fast.Div(ten)
	data.SafeLow = data.SafeLow.Div(ten)
	data.Average = data.Average.Div(ten)
	t.Log(data)
	err = data.save()
	if err != nil {
		t.Log(err.Error())
	}
	t.Fail()
}

func TestConsolidatedGasData(t *testing.T) {
	initViper()
	egsData, err := getGasData()
	if err != nil {
		t.Fatal("error getting gas prices" + err.Error())
	}
	ten := decimal.NewFromFloat(10.0)
	egsData.Fastest = egsData.Fastest.Div(ten)
	egsData.Fast = egsData.Fast.Div(ten)
	egsData.SafeLow = egsData.SafeLow.Div(ten)
	egsData.Average = egsData.Average.Div(ten)

	gasNowData, err := getGasNowDecimal()
	if err != nil {
		t.Fatal(err)
	}

	esData, err := etherscanGas()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(egsData.Fastest, egsData.Fast, egsData.SafeLow)
	fmt.Println(gasNowData.Fast, gasNowData.Medium, gasNowData.Safe)
	fmt.Println(esData.FastGasPrice, esData.ProposeGasPrice, esData.SafeGasPrice)

	index, err := template.ParseFiles("files/index.html")
	if err != nil {
		t.Fatal(err)
	}
	if err = index.Execute(os.Stdout, multiData{EGS: egsData, GN: gasNowData, ES: esData, Timestamp: time.Now().String()}); err != nil {
		t.Fatal(err)
	}
	t.Fail()
}

func TestMultiGas(t *testing.T) {
	initViper()
	gasHTML, median, data, err := multiGas()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(gasHTML)
	fmt.Println("median ", median)
	fmt.Printf("%+v\n", data)
	ioutil.WriteFile("files/filled3.html", []byte(gasHTML), 0777)
	t.Fail()
}
