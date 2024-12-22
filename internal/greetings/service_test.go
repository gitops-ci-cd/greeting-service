package greetings

import (
	"testing"

	pb "github.com/gitops-ci-cd/greeting-service/internal/_gen/pb/v1"
)

func TestLookup(t *testing.T) {
	service := &Service{}

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
	}

	// Iterate over test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Call the method under test
			actualLanguage, _ := service.Lookup(tc.language)

			// Validate the response
			if actualLanguage != tc.expectedLanguage {
				t.Errorf("got language %v, want %v", actualLanguage, tc.expectedLanguage)
			}
		})
	}
}
