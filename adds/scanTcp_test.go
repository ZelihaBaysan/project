package adds

import (
	"testing"
	"time"
)

// TestScanPortTCP verifies that the ScanPortTCP function correctly identifies the state of a TCP port.
//
// This test case uses a loopback address (localhost) and commonly open and closed ports to ensure
// that the function behaves as expected. It also includes tests with an invalid port to simulate
// network errors and timeouts.
func TestScanPortTCP(t *testing.T) {
	// Test cases
	tests := []struct {
		ip       string
		port     int
		expected string
	}{
		{"127.0.0.1", 80, Open},     // Commonly open port
		{"127.0.0.1", 9999, Closed}, // Commonly closed port
		{"127.0.0.1", 0, Filtered},  // Invalid port number
	}

	for _, tt := range tests {
		t.Run(tt.ip, func(t *testing.T) {
			// Allow some time for the scan to complete
			time.Sleep(1 * time.Second)

			result := ScanPortTCP(tt.ip, tt.port)
			if result != tt.expected {
				t.Errorf("ScanPortTCP(%s, %d) = %v; want %v", tt.ip, tt.port, result, tt.expected)
			}
		})
	}
}

// TestWorkerTCP verifies that the WorkerTCP function correctly processes and sends port scanning results.
//
// This test case simulates a scenario with multiple ports and checks whether the results and open
// port details are correctly sent to the appropriate channels.
func TestWorkerTCP(t *testing.T) {
	ip := "127.0.0.1"
	ports := make(chan int, 3)
	results := make(chan int, 3)
	openPorts := make(chan ServiceVersion, 3)
	done := make(chan bool)
	services := map[int]string{
		80: "HTTP",
	}

	go WorkerTCP(ip, ports, results, openPorts, done, services)

	// Send ports to scan
	ports <- 80
	ports <- 9999
	close(ports)

	// Wait for worker to finish
	<-done

	// Check results channel
	expectedResults := []int{80, 9999}
	for i := 0; i < len(expectedResults); i++ {
		result := <-results
		if result != expectedResults[i] {
			t.Errorf("Expected port %d, but got %d", expectedResults[i], result)
		}
	}

	// Check openPorts channel
	expectedOpenPorts := []ServiceVersion{
		{Port: 80, Protocol: "Unknown", Service: "HTTP", Response: "Service Detected"},
	}
	for i := 0; i < len(expectedOpenPorts); i++ {
		openPort := <-openPorts
		if openPort != expectedOpenPorts[i] {
			t.Errorf("Expected open port %+v, but got %+v", expectedOpenPorts[i], openPort)
		}
	}
}
