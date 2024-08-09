package main

import (
	"fmt"
	"net"
	"os"
	"port/adds"
)

// PortScanner defines a structure for scanning ports on a domain.
//
// Fields:
// - Domain: The domain to be scanned.
// - IPs: A slice of IP addresses associated with the domain.
// - Ports: A slice of ports to be scanned.
// - NumWorkers: The number of worker goroutines to use for scanning.
// - PortChannel: A channel used to distribute ports to be scanned by workers.
// - ResultChannel: A channel used to collect scan results from workers.
// - OpenPorts: A channel to collect open TCP ports and their detected services.
// - OpenPortsUDP: A channel to collect open UDP ports and their detected services.
// - Done: A channel to signal the completion of worker tasks.
// - Services: A map of known services with port numbers as keys and service names as values.
type PortScanner struct {
	Domain        string
	IPs           []string
	Ports         []int
	NumWorkers    int
	PortChannel   chan int
	ResultChannel chan int
	OpenPorts     chan adds.ServiceVersion
	OpenPortsUDP  chan adds.ServiceVersion
	Done          chan bool
	Services      map[int]string
}

// NewTarget creates a new PortScanner instance for a given domain.
//
// Parameters:
// - domain: The domain to be scanned.
// - numWorkers: The number of worker goroutines to use for scanning.
//
// Returns:
// - *PortScanner: A pointer to the created PortScanner instance.
// - error: An error if domain resolution fails.
//
// This function resolves the given domain to its associated IP addresses and
// prepares a list of ports from 1 to 65535 to be scanned.
func NewTarget(domain string, numWorkers int) (*PortScanner, error) {
	ips, err := net.LookupHost(domain)
	if err != nil {
		return nil, err
	}

	ports := make([]int, 0, 65535)
	for port := 1; port <= 65535; port++ {
		ports = append(ports, port)
	}

	return &PortScanner{
		Domain:        domain,
		IPs:           ips,
		Ports:         ports,
		NumWorkers:    numWorkers,
		PortChannel:   make(chan int, len(ports)),
		ResultChannel: make(chan int, len(ports)),
		OpenPorts:     make(chan adds.ServiceVersion, len(ports)),
		OpenPortsUDP:  make(chan adds.ServiceVersion, len(ports)),
		Done:          make(chan bool, numWorkers*2),
		Services:      make(map[int]string),
	}, nil
}

// Scan performs the port scanning operation on the resolved IP addresses.
//
// The function starts the specified number of worker goroutines for both TCP and UDP scanning,
// distributes the ports to be scanned through channels, and writes the results to an output file.
//
// The results include the status of each port (open or closed) and the detected services running on open ports.
func (t *PortScanner) Scan() {
	for i := 0; i < t.NumWorkers; i++ {
		go adds.WorkerTCP("", t.PortChannel, t.ResultChannel, t.OpenPorts, t.Done, t.Services)
		go adds.WorkerUDP("", t.PortChannel, t.ResultChannel, t.OpenPortsUDP, t.Done, t.Services)
	}

	for _, ip := range t.IPs {
		fmt.Printf("Scanning IP: %s\n", ip)

		for _, port := range t.Ports {
			fmt.Printf("Enqueueing port %d\n", port)
			t.PortChannel <- port
		}
		close(t.PortChannel)

		for range t.Ports {
			<-t.ResultChannel
		}

		for i := 0; i < t.NumWorkers*2; i++ {
			<-t.Done
		}
		close(t.OpenPorts)
		close(t.OpenPortsUDP)

		file, err := os.Create("output.txt")
		if err != nil {
			fmt.Printf("Error creating file: %s\n", err)
			return
		}
		defer file.Close()

		_, err = file.WriteString("Open TCP Ports with Services:\n")
		if err != nil {
			fmt.Printf("Error writing to file: %s\n", err)
			return
		}
		for service := range t.OpenPorts {
			_, err = file.WriteString(fmt.Sprintf("Port %d (TCP) is Open, Service: %s\n", service.Port, service.Service))
			if err != nil {
				fmt.Printf("Error writing to file: %s\n", err)
				return
			}
		}

		_, err = file.WriteString("Open UDP Ports with Services:\n")
		if err != nil {
			fmt.Printf("Error writing to file: %s\n", err)
			return
		}
		for service := range t.OpenPortsUDP {
			_, err = file.WriteString(fmt.Sprintf("Port %d (UDP) is Open, Service: %s\n", service.Port, service.Service))
			if err != nil {
				fmt.Printf("Error writing to file: %s\n", err)
				return
			}
		}
	}
}

// main is the entry point of the application.
//
// It prompts the user for a domain name, creates a PortScanner instance,
// and initiates the port scanning process.
func main() {
	var domain string
	fmt.Print("Enter domain: ")
	fmt.Scanln(&domain)

	target, err := NewTarget(domain, 100)
	if err != nil {
		fmt.Printf("Error resolving domain: %s\n", err)
		return
	}

	target.Scan()
}
