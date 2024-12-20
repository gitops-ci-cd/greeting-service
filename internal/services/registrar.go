package services

import (
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/gitops-ci-cd/greeting-service/internal/_gen/pb/v1"
	"github.com/gitops-ci-cd/greeting-service/internal/greetings"
)

// Register registers all gRPC service handlers and debugging capabilities with the server
func Register(server *grpc.Server) {
	// Register the GreetingService business logic
	pb.RegisterGreetingServiceServer(server, greetings.NewGreetingServiceHandler())

	// Register reflection service for debugging
	reflection.Register(server)

	for key := range server.GetServiceInfo() {
		slog.Info("Service registered", "service", key)
	}
}
