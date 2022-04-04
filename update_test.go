package main

import (
	"io/ioutil"
	"testing"

	"github.com/spf13/viper"
)

func TestUpdate(t *testing.T) {
	initViper()
	version := viper.GetFloat64("Version")
	html, err := newInVersion(version)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile("files/updateFilled.html", []byte(html), 0777)
	if err != nil {
		t.Fatal(err)
	}
}
