package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// loggingInterceptor logs all incoming gRPC requests
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		// Log the incoming request
		if protoMsg, ok := req.(proto.Message); ok {
			marshaler := protojson.MarshalOptions{
				AllowPartial:    true,
				EmitUnpopulated: true,
				UseProtoNames:   true,
			}
			fields := []any{
				"type", fmt.Sprintf("%T", req),
			}
			if bytes, err := marshaler.Marshal(protoMsg); err == nil {
				fields = append(fields, "request", string(bytes))
				slog.Debug("Incoming gRPC request", fields...)
			} else {
				fields = append(fields, "error", err)
				slog.Error("Failed to marshal Protobuf message", fields...)
			}
		}

		// Log the incoming context
		md, _ := metadata.FromIncomingContext(ctx)
		slog.Debug("Incoming gRPC context", "metadata", md)
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
