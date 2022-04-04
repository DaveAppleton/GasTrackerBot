package main

import (
	"fmt"
	"testing"
)

func TestGasNowDecimal(t *testing.T) {
	initViper()
	gnd, err := getGasNowDecimal()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", gnd)
	t.Fail()
}
