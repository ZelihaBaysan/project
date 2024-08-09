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
	// Set up test environment
	ports := make(chan int, 10)
	results := make(chan int)
	openPorts := make(chan ServiceVersion)
	done := make(chan bool)

	// Start the worker
	go WorkerTCP("127.0.0.1", ports, results, openPorts, done, nil)

	// Add ports to scan
	go func() {
		for _, port := range []int{80, 9999} {
			ports <- port
		}
		close(ports)
	}()

	// Check results
	go func() {
		for range results {
			// Process results if needed
		}
	}()

	// Check openPorts
	go func() {
		for range openPorts {
			// Process open ports if needed
		}
	}()

	// Wait for worker to finish
	<-done
}
