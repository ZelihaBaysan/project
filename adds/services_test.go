package adds

import (
	"testing"
)

// TestServicess verifies that the Servicess map correctly associates port numbers with service names.
//
// This test case ensures that the port numbers in the Servicess map are associated with the expected service names.
// It checks both known entries and some that are not present in the map to verify proper handling.
func TestServicess(t *testing.T) {
	// Define expected results
	expected := map[int]string{
		1:    "tcpmux",
		7:    "echo",
		13:   "daytime",
		19:   "chargen",
		9999: "Unknown", // Port not in the map
	}

	// Iterate over expected results and check if the Servicess map has correct values
	for port, service := range expected {
		if result, ok := Servicess[port]; ok {
			if result != service {
				t.Errorf("Servicess[%d] = %v; want %v", port, result, service)
			}
		} else if service != "Unknown" {
			t.Errorf("Servicess[%d] not found; expected %v", port, service)
		}
	}
}
