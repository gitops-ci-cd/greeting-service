package greetings

import (
	"context"
	"math/rand"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/gitops-ci-cd/greeting-service/internal/_gen/pb/v1"
)

// greetingServiceHandler implements the GreetingServiceServer interface.
type greetingServiceHandler struct {
	pb.UnimplementedGreetingServiceServer // Embedding for forward compatibility
}

// NewGreetingServiceHandler creates a new instance of greetingServiceHandler.
func NewGreetingServiceHandler() pb.GreetingServiceServer {
	return &greetingServiceHandler{}
}

// Fetch handles an RPC request for a greeting.
func (h *greetingServiceHandler) Fetch(ctx context.Context, req *pb.GreetingRequest) (*pb.GreetingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	// Default to English if language is explicitly UNKNOWN
	language := req.Language
	if language == pb.Language_UNKNOWN {
		language = pb.Language_EN
	}

	// Default to English if language is not recognized
	_, exists := greetingData[language]
	if !exists {
		language = pb.Language_EN
	}

	return &pb.GreetingResponse{
		Language:  language,
		Greeting:  getRandomGreeting(language),
		Timestamp: timestamppb.New(time.Now()),
	}, nil
}

// Define greetings per language
var greetingData = map[pb.Language][]string{
	pb.Language_EN:    {"Hello", "Hi", "Hey", "Greetings"},
	pb.Language_EN_GB: {"Hello", "Hiya", "Cheers", "Greetings"},
	pb.Language_ES:    {"Hola", "Qué tal", "Buenos días", "Saludos"},
	pb.Language_FR:    {"Bonjour", "Salut", "Coucou", "Bienvenue"},
}

// Randomly select one greeting for the given language
func getRandomGreeting(language pb.Language) string {
	greetings := greetingData[language]

	return greetings[rand.Intn(len(greetings))]
}
