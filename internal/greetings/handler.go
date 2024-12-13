package handlers

// import (
// 	"context"
// 	"log"
// 	"time"

// 	"google.golang.org/protobuf/types/known/timestamppb"

// 	pb "github.com/gitops-ci-cd/schema/v1"
// )

// // greeterHandler implements the GreeterServer interface.
// type greeterHandler struct {
// 	pb.UnimplementedGreeterServer // Embedding for forward compatibility
// }

// // NewGreeterHandler creates a new instance of greeterHandler.
// func NewGreeterHandler() pb.GreeterServer {
// 	return &greeterHandler{}
// }

// // Greeting handles the GetGreeting RPC.
// func (h *greeterHandler) Greeting(ctx context.Context, req *pb.GreetingRequest) (*pb.GreetingResponse, error) {
// 	log.Printf("Received request for language: %v", req.Language)

// 	// Map language enum to greeting
// 	var greeting string
// 	switch req.Language {
// 	case pb.Language_EN:
// 		greeting = "Hello"
// 	case pb.Language_EN_GB:
// 		greeting = "Hello (British)"
// 	case pb.Language_ES:
// 		greeting = "Hola"
// 	case pb.Language_FR:
// 		greeting = "Bonjour"
// 	default:
// 		greeting = "Hello (Default)"
// 	}

// 	// Return the response
// 	return &pb.GreetingResponse{
// 		Language:  req.Language,
// 		Greeting:  greeting,
// 		Timestamp: timestamppb.New(time.Now()),
// 	}, nil
// }
