package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// loggingInterceptor logs all incoming gRPC requests
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	// Marshal the request to JSON
	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		reqProto, ok := req.(proto.Message)
		if !ok {
			slog.Warn("Request is not a proto.Message", "type", fmt.Sprintf("%T", req))
		} else {
			reqJSON, err := protojson.Marshal(reqProto)
			if err != nil {
				slog.Error("Failed to marshal request", "error", err)
			} else {
				slog.Debug("Incoming gRPC request", "method", info.FullMethod, "request", string(reqJSON))
			}
		}
	}

	// Process the request
	res, err := handler(ctx, req)
	duration := time.Since(start)

	fields := []any{
		"method", info.FullMethod,
		"duration", duration.String(),
	}

	if err != nil {
		fields = append(fields, "error", err)
	}

	slog.Info("Handled gRPC request", fields...)

	return res, err
}
