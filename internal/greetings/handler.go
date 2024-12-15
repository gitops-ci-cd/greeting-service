package greetings

import (
	"context"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/gitops-ci-cd/greeting-service/internal/gen/pb/v1"
)

// greetingServiceHandler implements the GreetingServiceServer interface.
type greetingServiceHandler struct {
	pb.UnimplementedGreetingServiceServer // Embedding for forward compatibility
}

// NewGreeterHandler creates a new instance of greetingServiceHandler.
func NewGreeterHandler() pb.GreetingServiceServer {
	return &greetingServiceHandler{}
}

// Greeting handles an RPC request for a greeting.
func (h *greetingServiceHandler) Fetch(ctx context.Context, req *pb.GreetingRequest) (*pb.GreetingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	log.Printf("Received request for language: %v", req.Language)

	// Default to English if language is not specified
	language := req.Language
	if language == pb.Language_UNKNOWN {
		language = pb.Language_EN
	}

	return &pb.GreetingResponse{
		Language:  language,
		Greeting:  getRandomGreeting(language),
		Timestamp: timestamppb.New(time.Now()),
	}, nil
}

func getRandomGreeting(language pb.Language) string {
	// Define greetings per language
	greetings := map[pb.Language][]string{
		pb.Language_EN:    {"Hello", "Hi", "Hey", "Greetings"},
		pb.Language_EN_GB: {"Hello", "Hiya", "Cheers", "Greetings"},
		pb.Language_ES:    {"Hola", "Qué tal", "Buenos días", "Saludos"},
		pb.Language_FR:    {"Bonjour", "Salut", "Coucou", "Bienvenue"},
	}

	// Randomly select one greeting for the given language
	return greetings[language][rand.Intn(len(greetings[language]))]
}
