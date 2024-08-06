package main

import (
	"fmt"
	"net"
	"port/adds"
)

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
	openPorts := make(chan adds.ServiceVersion, len(ports))
	done := make(chan bool, 100)

	numWorkers := 100
	for i := 0; i < numWorkers; i++ {
		go adds.Worker("", portChannel, resultChannel, openPorts, done, adds.Servicess)
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

		fmt.Println("Open Ports with Services:")
		for service := range openPorts {
			fmt.Printf("Port %d is Open, Service: %s\n", service.Port, service.Service)
		}
	}

	close(resultChannel)
}
