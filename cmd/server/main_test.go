package main

import (
	"net"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port = ":50051"
	}

	// Bind to the same port as the gRPC server to simulate a failure
	listener, err := net.Listen("tcp", port)
	if err != nil {
		t.Fatalf("Failed to bind to port 50051 for testing: %v", err)
	}
	defer listener.Close()

	// Run the server and expect it to fail (port already in use)
	err = run(port, listener)
	if err == nil {
		t.Fatal("Expected run() to fail when the port is already in use")
	}
}
