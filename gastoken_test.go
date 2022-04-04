package main

import (
	"fmt"
	"testing"
)

func TestCHI(t *testing.T) {
	initViper()
	eth := getTokenPrice("chi-gastoken")
	if eth.OK {
		fmt.Println(" 1 ETH buys", eth.ToTokenAmount, "CHI")
		fmt.Println(eth.Message)
	}
	t.Fail()
}

func TestGAS(t *testing.T) {
	initViper()
	eth := getTokenPrice("gastoken")
	if eth.OK {
		fmt.Println(" 1 ETH buys", eth.ToTokenAmount, "GST2")
		fmt.Println(eth.Message)
	}
	t.Fail()
}

func Test1InchGST(t *testing.T) {
	initViper()
	eth := getOneInchPrice("GST2")
	if eth.OK {
		fmt.Printf("1 ETH buys %s GST2\n", eth.ToTokenAmount)
		fmt.Println(eth.Message)
	}
	t.Fail()
}

func Test1InchCHI(t *testing.T) {
	initViper()
	eth := getOneInchPrice("CHI")
	if eth.OK {
		fmt.Printf(" 1 ETH buys %s CHI\n", eth.ToTokenAmount)
		fmt.Println(eth.Message)
	}
	t.Fail()
}

// 1 000 000 000 000 000 000
