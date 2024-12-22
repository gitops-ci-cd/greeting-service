package greetings

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/gitops-ci-cd/greeting-service/internal/_gen/pb/v1"
)

type mockService struct{}

func (m *mockService) Lookup(language pb.Language) (pb.Language, string) {
	// Mock response based on input language
	switch language {
	case pb.Language_EN:
		return pb.Language_EN, "Hello"
	case pb.Language_EN_GB:
		return pb.Language_EN_GB, "Hiya"
	default:
		return pb.Language_EN, "Hello"
	}
}

func TestHandlerFetch(t *testing.T) {
	mockSvc := &mockService{}
	handler := &Handler{Service: mockSvc}

	tests := []struct {
		name     string
		req      *pb.GreetingRequest
		wantResp *pb.GreetingResponse
		wantErr  codes.Code
	}{
		{
			name: "Valid English request",
			req: &pb.GreetingRequest{
				Language: pb.Language_EN,
			},
			wantResp: &pb.GreetingResponse{
				Language:  pb.Language_EN,
				Greeting:  "Hello",
			},
			wantErr: codes.OK,
		},
		{
			name: "Valid British English request",
			req: &pb.GreetingRequest{
				Language: pb.Language_EN_GB,
			},
			wantResp: &pb.GreetingResponse{
				Language:  pb.Language_EN_GB,
				Greeting:  "Hiya",
			},
			wantErr: codes.OK,
		},
		{
			name:    "Nil request",
			req:     nil,
			wantErr: codes.InvalidArgument,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := handler.Fetch(context.Background(), tc.req)

			// Validate error code
			if status.Code(err) != tc.wantErr {
				t.Errorf("got error code %v, want %v", status.Code(err), tc.wantErr)
			}

			// Validate response if no error
			if tc.wantErr == codes.OK {
				if resp.Language != tc.wantResp.Language {
					t.Errorf("got language %v, want %v", resp.Language, tc.wantResp.Language)
				}
				if resp.Greeting != tc.wantResp.Greeting {
					t.Errorf("got greeting %v, want %v", resp.Greeting, tc.wantResp.Greeting)
				}
			}
		})
	}
}
