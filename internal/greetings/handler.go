package greetings

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/gitops-ci-cd/greeting-service/internal/_gen/pb/v1"
)

// Handler implements the GreetingServiceServer interface
type Handler struct {
	// Embedding for forward compatibility
	pb.UnimplementedGreetingServiceServer
	Service service
}

// Register associates the handler with the given gRPC server
func (h *Handler) Register(server *grpc.Server) {
	pb.RegisterGreetingServiceServer(server, h)
}

// Fetch handles an RPC request
func (h *Handler) Fetch(ctx context.Context, req *pb.GreetingRequest) (*pb.GreetingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	language, greeting := h.Service.Lookup(req.Language)

	return &pb.GreetingResponse{
		Language:  language,
		Greeting:  greeting,
		Timestamp: timestamppb.New(time.Now()),
	}, nil
}
