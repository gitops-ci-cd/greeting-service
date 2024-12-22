package main

import (
	"log/slog"
	"os"

	"google.golang.org/grpc"

	"github.com/gitops-ci-cd/greeting-service/internal/greetings"
	"github.com/gitops-ci-cd/greeting-service/pkg/io"
	"github.com/gitops-ci-cd/greeting-service/pkg/telemetry"
)

const defaultPort = "50051"

// Configure the logger
func init() {
	level := func() slog.Level {
		switch os.Getenv("LOG_LEVEL") {
		case "ERROR":
			return slog.LevelError
		case "WARN":
			return slog.LevelWarn
		case "INFO":
			return slog.LevelInfo
		case "DEBUG":
			return slog.LevelDebug
		default:
			return slog.LevelInfo
		}
	}()
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})))
}

// main is the entry point for the server
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(telemetry.LoggingInterceptor),
	}
	server := grpc.NewServer(opts...)

	// Create and populate the registry with gRPC services that satisfy the io.Registerable interface
	registry := &io.Registry{}
	registry.Add(&greetings.Handler{Service: &greetings.Service{}})
	registry.RegisterAll(server)

	if err := io.Run(":"+port, server); err != nil {
		slog.Error("Server terminated", "error", err)
		os.Exit(1)
	} else {
		slog.Warn("Server stopped")
	}
}
