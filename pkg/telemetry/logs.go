package telemetry

import (
	"context"
	"log"
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
		log.Printf("INTERCEPTOR: Request: %v", req)
		reqJSON, err := protojson.Marshal(req.(proto.Message))
		if err != nil {
			slog.Error("Failed to marshal request", "error", err)
			reqJSON = []byte("{}")
		}
		slog.Debug("Incoming gRPC request", "method", info.FullMethod, "request", string(reqJSON))
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
