package main

import (
	"fmt"
	"os"
	"testing"
	"text/template"

	"github.com/spf13/viper"
)

func TestTextTemplate(t *testing.T) {
	initViper()
	var myData multiData
	myData.Version = fmt.Sprintf("%4.2f", viper.GetFloat64("Version"))
	myData.Advert = viper.GetString("ADVERT")
	index, err := template.ParseFiles("files/index.md")
	if err != nil {
		t.Fatal(err)
	}
	err = index.Execute(os.Stdout, myData)
	if err != nil {
		t.Fatal(err)
	}
}
