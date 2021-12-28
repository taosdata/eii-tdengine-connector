package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		fmt.Println("TDengineConnector is running...")
		time.Sleep(time.Duration(2) * time.Second)
	}
}
