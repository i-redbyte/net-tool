package main

import (
	"fmt"
)

func main() {
	fmt.Printf("*****\tWELCOME TO RED SCANNER!\t*****\n")
	_ = StartScan("tcp", "192.168.1.1")
	//_ = StartScan("udp", "192.168.1.1")
}
