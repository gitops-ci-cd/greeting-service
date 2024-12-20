package telemetry

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// mockUnaryHandler simulates a gRPC unary handler.
func mockUnaryHandler(ctx context.Context, req interface{}) (interface{}, error) {
	// Simulate processing time
	time.Sleep(50 * time.Millisecond)

	// Return a mocked response
	return &timestamppb.Timestamp{Seconds: 12345, Nanos: 67890}, nil
}

// mockProtoRequest creates a mock proto.Message for testing.
type mockProtoRequest struct {
	Message string
}

func (m *mockProtoRequest) Reset()         {}
func (m *mockProtoRequest) String() string { return m.Message }
func (m *mockProtoRequest) ProtoMessage()  {}

func TestLoggingInterceptor(t *testing.T) {
	interceptor := LoggingInterceptor

	// Create a test request
	req := &mockProtoRequest{Message: "test request"}
	ctx := context.Background()
	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/Method",
	}

	// Invoke the interceptor with the mock handler
	resp, err := interceptor(ctx, req, info, mockUnaryHandler)

	// Check if there were errors
	if err != nil {
		t.Fatalf("Interceptor returned an error: %v", err)
	}

	// Check the response type
	tsResp, ok := resp.(*timestamppb.Timestamp)
	if !ok {
		t.Fatalf("Response type mismatch. Expected *timestamppb.Timestamp, got %T", resp)
	}

	// Validate response content
	if tsResp.Seconds != 12345 || tsResp.Nanos != 67890 {
		t.Errorf("Unexpected response content. Got: %v", tsResp)
	}
}
