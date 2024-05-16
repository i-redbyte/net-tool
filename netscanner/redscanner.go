package netscanner

import (
	"errors"
	"fmt"
	"net"
	sp "net-tool/spinner"
	"sort"
	"time"
)

const PortCount = 65535

type ScanResult struct {
	Protocol string
	Port     int
	State    string
}

type PortSorter []ScanResult

func (a PortSorter) Len() int      { return len(a) }
func (a PortSorter) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a PortSorter) Less(i, j int) bool {
	return a[i].Port < a[j].Port
}

type ScanObject struct {
	Protocol string
	Hostname string
}

func StartScan(protocol, host string) error {
	hostname, err := hostnameValidator(host)
	if err != nil {
		return err
	}
	ports := make(chan int, 100)
	results := make(chan ScanResult)
	done := make(chan bool)
	scan := ScanObject{Hostname: hostname, Protocol: protocol}
	var scanResults []ScanResult

	go sp.Spinner(done, "")

	for i := 0; i < cap(ports); i++ {
		go worker(scan, ports, results)
	}

	go func() {
		for i := 1; i <= PortCount; i++ {
			ports <- i
		}
		close(ports)
	}()

	for i := 1; i <= PortCount; i++ {
		scanResult := <-results
		if scanResult.Port != -1 {
			scanResults = append(scanResults, scanResult)
		}
	}
	sort.Sort(PortSorter(scanResults))
	done <- true

	printScanResult(scanResults)
	return nil
}

func printScanResult(scanResults []ScanResult) {
	fmt.Printf("\n\t\tSCAN RESULT:\n")
	fmt.Printf("==================================================\n")
	for _, port := range scanResults {
		fmt.Printf("\t%s\t%d\t%s\n", port.Protocol, port.Port, port.State)
	}
	fmt.Printf("==================================================\n")
}

func hostnameValidator(hostname string) (string, error) {
	n := len(hostname)
	address := hostname
	if n > 0 && hostname[n-1] != ':' {
		address = address + ":"
	} else if n == 0 {
		return "", errors.New("missing address to scan")
	}
	return address, nil
}

func worker(scan ScanObject, ports chan int, results chan ScanResult) {
	for port := range ports {
		address := fmt.Sprintf("%s%d", scan.Hostname, port)
		conn, err := net.DialTimeout(scan.Protocol, address, 1*time.Second)
		if err != nil {
			results <- ScanResult{Protocol: scan.Protocol, Port: -1, State: "Closed"}
			continue
		}
		conn.Close()
		results <- ScanResult{Protocol: scan.Protocol, Port: port, State: "Open"}
	}
}
