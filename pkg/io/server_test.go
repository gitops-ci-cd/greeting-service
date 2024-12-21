package io

import (
	"net"
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

	mockServer := new(Server)

	// Call the function under test
	err = Run(port, *mockServer)
	if err == nil {
		t.Errorf("Expected run() to fail when the port is already in use")
	}
}
