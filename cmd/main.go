package main

import (
	"fmt"
	ns "net-tool/netscanner"
)

func main() {
	fmt.Printf("*****\tWELCOME TO RED SCANNER!\t*****\n")
	_ = ns.StartScan("tcp", "192.168.1.1")
	//_ = StartScan("udp", "192.168.1.1")
}
