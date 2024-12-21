package io

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// Server is the expected interface for servers to implement
type Server interface {
	Serve(net.Listener) error
	GracefulStop()
}

// Run sets up and starts a gRPC server
func Run(port string, server Server) error {
	// Create a TCP listener
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("could not create tcp listener on port %s: %w", port, err)
	}
	defer listener.Close()

	// Run the server in a goroutine to allow for graceful shutdown
	ctx := setupSignalHandler()
	go func() {
		slog.Info("Server listening...", "port", port)
		if err := server.Serve(listener); err != nil {
			slog.Error("gRPC server failed", "error", err)
		}
	}()

	// Wait for termination signal
	<-ctx.Done()
	slog.Warn("Server shutting down gracefully...")
	server.GracefulStop()

	return nil
}

// setupSignalHandler creates a cancellable context for signal handling
func setupSignalHandler() context.Context {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ctx.Done()
		slog.Warn("Received termination signal")
		stop()
	}()
	return ctx
}
