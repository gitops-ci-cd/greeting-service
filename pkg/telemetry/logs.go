package telemetry

import (
	"context"
	"encoding/json"
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
		if protoMsg, ok := req.(proto.Message); ok {
			marshaler := protojson.MarshalOptions{
				AllowPartial:    true,
				EmitUnpopulated: true,
				UseProtoNames:   true,
			}
			if bytes, err := marshaler.Marshal(protoMsg); err == nil {
				// Parse the JSON to avoid escaping
				var parsedJSON map[string]interface{}
				if err := json.Unmarshal(bytes, &parsedJSON); err == nil {
					slog.Debug("Incoming gRPC request", "method", info.FullMethod, "type", fmt.Sprintf("%T", req), "request", parsedJSON)
				} else {
					slog.Error("Failed to unmarshal JSON", "error", err)
				}
			} else {
				slog.Error("Failed to marshal Protobuf message", "method", info.FullMethod, "type", fmt.Sprintf("%T", req), "error", err)
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
