package main

import (
	"fmt"
	"testing"
)

func TestEscan(t *testing.T) {
	initViper()
	etherscanData, err := etherscanGas()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(etherscanData.FastGasPrice, etherscanData.ProposeGasPrice, etherscanData.SafeGasPrice)
	t.Fail()
}
