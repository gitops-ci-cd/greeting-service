package greetings

import (
	"context"
	"testing"
	"time"

	pb "github.com/gitops-ci-cd/greeting-service/internal/_gen/pb/v1"
)

func TestFetch_ValidRequest(t *testing.T) {
	handler := &GreetingServiceHandler{}

	// Define test cases
	tests := []struct {
		name             string
		language         pb.Language
		expectedLanguage pb.Language
	}{
		{
			name:             "English",
			language:         pb.Language_EN,
			expectedLanguage: pb.Language_EN,
		},
		{
			name:             "British English",
			language:         pb.Language_EN_GB,
			expectedLanguage: pb.Language_EN_GB,
		},
		{
			name:             "Unknown language",
			language:         pb.Language_UNKNOWN,
			expectedLanguage: pb.Language_EN,
		},
		{
			name:             "Unrecognized language",
			language:         999,
			expectedLanguage: pb.Language_EN,
		},
		// Add other cases...
	}

	// Iterate over test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &pb.GreetingRequest{Language: tc.language}

			// Call the method under test
			resp, err := handler.Fetch(context.Background(), req)

			// Validate the response
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.Language != tc.expectedLanguage {
				t.Errorf("got language %v, want %v", resp.Language, tc.expectedLanguage)
			}
			if time.Since(resp.Timestamp.AsTime()) > time.Second {
				t.Errorf("timestamp is too old: %v", resp.Timestamp.AsTime())
			}
		})
	}
}
