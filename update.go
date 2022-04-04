package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"text/template"

	"github.com/shopspring/decimal"
)

func readMD(filename string) (result string) {
	result = "sorry cannot process that right now..."
	messageBytes, err := ioutil.ReadFile("files/" + filename)
	if err != nil {
		log.Println("reading sources.md ", err)
	} else {
		result = string(messageBytes)
	}
	return
}

func newInVersion(version float64) (output string, err error) {
	mfName := fmt.Sprintf("NewIn%4.2f.md", version)
	output = readMD(mfName)
	return
}

func levelResponse(target decimal.Decimal, median string) (output string, err error) {
	lrData := struct {
		Target  decimal.Decimal
		Current string
		Advert  string
	}{
		Target:  target,
		Current: median,
	}
	lrData.Advert = readMD("advert.md")
	index, err := template.ParseFiles("files/level_response.md")
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	if err = index.Execute(buf, lrData); err != nil {
		return
	}
	output = buf.String()
	return
}

func alertResponse(wait int, median string) (output string, err error) {
	lrData := struct {
		Target  int
		Current string
		Advert  string
	}{
		Target:  wait,
		Current: median,
	}
	lrData.Advert = readMD("advert.md")
	index, err := template.ParseFiles("files/alert_response.md")
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	if err = index.Execute(buf, lrData); err != nil {
		return
	}
	output = buf.String()
	return
}

func welcomeMessage() (output string, err error) {
	var lrData struct {
		Advert string
	}
	lrData.Advert = readMD("advert.md")
	index, err := template.ParseFiles("files/welcome.md")
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	if err = index.Execute(buf, lrData); err != nil {
		return
	}
	output = buf.String()
	return
}

func levelError(data string) (output string) {
	data = strings.Replace(data, "\\", "", 20)
	lrData := struct {
		Data   string
		Advert string
	}{
		Data: data,
	}
	lrData.Advert = readMD("advert.md")
	index, err := template.ParseFiles("files/level_error.md")
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	if err = index.Execute(buf, lrData); err != nil {
		return
	}
	output = buf.String()
	return
}
