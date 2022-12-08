package main

import (
	"fmt"
	"net"
	"sort"
)

const PortCount = 1024

func worker(ports, results chan int) {
	for port := range ports {
		address := fmt.Sprintf("192.168.1.1:%d", port)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- port
	}
}
func main() {
	fmt.Printf("Welcome to red scanner!\n")
	ports := make(chan int, 100)
	results := make(chan int)
	defer close(ports)
	defer close(results)
	var openports []int
	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	go func() {
		for i := 1; i <= PortCount; i++ {
			ports <- i
		}
	}()

	for i := 0; i < PortCount-1; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}
	sort.Ints(openports)
	fmt.Printf("==================================================\n")
	fmt.Printf("\t\tScan result:\n")
	fmt.Printf("==================================================\n")
	for _, port := range openports {
		fmt.Printf("\t\t\t%d\topen\n", port)
	}
	fmt.Printf("==================================================\n")
}
