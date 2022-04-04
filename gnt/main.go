package main

import (
	"fmt"
	"os"
	"time"
)

func testLoopData(g chan gasNowData, b chan bool) {
	for {
		select {
		case gd := <-g:
			fmt.Println(gd.Fast, gd.Medium, gd.Safe)
			continue
		case <-b:
			fmt.Println("leaving now")
			return
		case <-time.After(1 * time.Minute):
			fmt.Println("one minute gone")
		}
	}
}

func main() {
	gndChan := make(chan gasNowData, 10)
	forceExitChan := make(chan bool, 2)
	go testLoopData(gndChan, forceExitChan)
	go gasNowLoop(gndChan, forceExitChan)
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		forceExitChan <- true
		forceExitChan <- true
		done <- true
	}()
	fmt.Println("Waiting for an exit")
	<-done
	fmt.Println("out of here")
}
