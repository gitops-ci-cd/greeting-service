package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/gitops-ci-cd/greeter/internal/gen/pb/v1"
	"github.com/gitops-ci-cd/greeter/internal/greetings"
)

func main() {
	// Define the port
	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port = ":50051"
	}

	// Create a TCP listener
	listener, err := net.Listen("tcp", port)
	if err != nil {
		slog.Error("Unable to listen on port", "port", port, "error", err)
		os.Exit(1)
	}

	// Run the server. Separated into a function to better facilitate testing
	if err := run(port, listener); err != nil {
		slog.Error("Server terminated with error", "error", err)
		os.Exit(1)
	} else {
		slog.Warn("Server stopped")
	}
}

// run sets up and starts the gRPC server
func run(port string, listener net.Listener) error {
	// Create a new gRPC server
	server := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)

	// Register the Greeter service
	pb.RegisterGreeterServer(server, greetings.NewGreeterHandler())

	// Register reflection service for debugging
	reflection.Register(server)

	// Handle graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Run the server in a goroutine
	go func() {
		slog.Info("Server listening...", "port", port)
		if err := server.Serve(listener); err != nil {
			slog.Error("Failed to serve", "error", err)
			stop() // Trigger graceful shutdown
		}
	}()

	// Wait for termination signal
	<-ctx.Done()
	slog.Warn("Server shutting down gracefully...")
	server.GracefulStop()

	return nil
}

func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	res, err := handler(ctx, req)
	duration := time.Since(start)

	// Add more structured logging
	slog.Info("Handled gRPC request",
		"method", info.FullMethod,
		"duration", duration.String(),
		"error", err,
	)

	return res, err
}
