package main

import (
	"fmt"
	"net"
	"port/scan"
	"time"
)

const (
	Open     = "Open"
	Closed   = "Closed"
	Filtered = "Filtered"
)

type ServiceVersion struct {
	Port     int
	Protocol string
	Response string
}

func scanPort(ip string, port int) string {
	address := fmt.Sprintf("%s:%d", ip, port)

	timeout := 5 * time.Second
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		if _, ok := err.(*net.OpError); ok {
			return Filtered
		}
		return Closed
	}
	conn.Close()
	return Open
}

func detectService(port int) ServiceVersion {
	if svc, ok := scan.Servicess[port]; ok {
		return ServiceVersion{Port: port, Protocol: "tcp", Response: svc}
	}
	return ServiceVersion{Port: port, Protocol: "Unknown", Response: "Unknown"}
}

func worker(ip string, ports, results chan int, openPorts chan int, done chan bool) {
	for port := range ports {
		state := scanPort(ip, port)
		service := detectService(port)
		results <- port
		if state == Open {
			openPorts <- port
		}
		fmt.Printf("Port %d: %s, Service: %s, Response: %s\n", port, state, service.Protocol, service.Response)
	}
	done <- true
}

func main() {
	var domain string
	fmt.Print("Enter domain: ")
	fmt.Scanln(&domain)

	ips, err := net.LookupHost(domain)
	if err != nil {
		fmt.Printf("Error resolving domain: %s\n", err)
		return
	}

	ports := make([]int, 0, 65535)
	for port := 1; port <= 65535; port++ {
		ports = append(ports, port)
	}

	portChannel := make(chan int, len(ports))
	resultChannel := make(chan int, len(ports))
	openPorts := make(chan int, len(ports))
	done := make(chan bool, 100)

	numWorkers := 100
	for i := 0; i < numWorkers; i++ {
		go worker("", portChannel, resultChannel, openPorts, done)
	}

	for _, ip := range ips {
		fmt.Printf("Scanning IP: %s\n", ip)

		for _, port := range ports {
			portChannel <- port
		}
		close(portChannel)

		for range ports {
			<-resultChannel
		}

		for i := 0; i < numWorkers; i++ {
			<-done
		}
		close(openPorts)

		fmt.Println("Open Ports:")
		for port := range openPorts {
			fmt.Printf("Port %d is Open\n", port)
		}
	}

	close(resultChannel)
}
