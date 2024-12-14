package main

import (
	"net"
	"strings"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	// Start a dummy server to occupy the port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to allocate a dynamic port: %v", err)
	}
	port := listener.Addr().String()

	defer listener.Close()

	// Give the dummy server time to fully bind
	time.Sleep(100 * time.Millisecond)

	// Call the function under test
	err = run(port)
	if err == nil {
		t.Errorf("Expected run() to fail when the port is already in use")
	} else if !strings.Contains(err.Error(), "address already in use") {
		// Check if the error message contains "address already in use"
		t.Errorf("Unexpected error when the port is already in use: %v", err)
	} else {
		t.Logf("run() failed as expected with error: %v", err)
	}
}
