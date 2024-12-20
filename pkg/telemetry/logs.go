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

	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		reqJSON := "{}" // Default empty JSON
		if protoMsg, ok := req.(proto.Message); ok {
			if bytes, err := protojson.Marshal(protoMsg); err == nil {
				reqJSON = string(bytes)
			} else {
				slog.Error("Failed to marshal request to JSON", "method", info.FullMethod, "type", fmt.Sprintf("%T", req), "error", err)
			}
		}
		slog.Debug("Incoming gRPC request", "method", info.FullMethod, "type", fmt.Sprintf("%T", req), "request", reqJSON)
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
