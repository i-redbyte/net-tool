package netscanner

import (
	"errors"
	"fmt"
	"net"
	sp "net-tool/spinner"
	"sort"
	"sync"
	"time"
)

const (
	PortCount     = 65535
	MaxGoroutines = 1000
	Timeout       = 100 * time.Millisecond
)

type ScanResult struct {
	Protocol string
	Port     int
	State    string
}

type PortSorter []ScanResult

func (a PortSorter) Len() int           { return len(a) }
func (a PortSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PortSorter) Less(i, j int) bool { return a[i].Port < a[j].Port }

type ScanObject struct {
	Protocol string
	Hostname string
}

func StartScan(protocol, host string) error {
	startTime := time.Now()

	hostname, err := resolveHost(host)
	if err != nil {
		return err
	}
	fmt.Printf("Resolved hostname: %s\n", hostname)
	results := make(chan ScanResult, PortCount)
	done := make(chan bool)
	var scanResults []ScanResult
	var wg sync.WaitGroup

	go sp.Spinner(done, "Done")

	portsPerGoroutine := PortCount / MaxGoroutines

	for i := 0; i < MaxGoroutines; i++ {
		wg.Add(1)
		startPort := i * portsPerGoroutine
		endPort := startPort + portsPerGoroutine
		if i == MaxGoroutines-1 {
			endPort = PortCount + 1
		}
		go func(start, end int) {
			defer wg.Done()
			for port := start; port < end; port++ {
				scanPort(protocol, hostname, port, results)
			}
		}(startPort, endPort)
	}

	go func() {
		wg.Wait()
		close(results)
		done <- true
	}()

	for scanResult := range results {
		if scanResult.Port != -1 {
			scanResults = append(scanResults, scanResult)
		}
	}
	sort.Sort(PortSorter(scanResults))

	elapsedTime := time.Since(startTime).Seconds()
	printScanResult(scanResults, elapsedTime)
	return nil
}

func resolveHost(host string) (string, error) {
	ip := net.ParseIP(host)
	if ip == nil {
		ips, err := net.LookupIP(host)
		if err != nil {
			return "", err
		}
		if len(ips) == 0 {
			return "", errors.New("no IP addresses found for host")
		}
		ip = ips[0]
	}
	return ip.String() + ":", nil
}

func printScanResult(scanResults []ScanResult, elapsedTime float64) {
	fmt.Printf("\n\t\tSCAN RESULT (%.2f seconds):\n", elapsedTime)
	fmt.Printf("==================================================\n")
	if len(scanResults) == 0 {
		fmt.Printf("No open ports found.\n")
	}
	for _, port := range scanResults {
		fmt.Printf("\t%s\t%d\t%s\n", port.Protocol, port.Port, port.State)
	}
	fmt.Printf("==================================================\n")
}

func scanPort(protocol, hostname string, port int, results chan ScanResult) {
	address := fmt.Sprintf("%s%d", hostname, port)
	conn, err := net.DialTimeout(protocol, address, Timeout)
	if err != nil {
		results <- ScanResult{Protocol: protocol, Port: -1, State: "Closed"}
		return
	}
	conn.Close()
	results <- ScanResult{Protocol: protocol, Port: port, State: "Open"}
}
