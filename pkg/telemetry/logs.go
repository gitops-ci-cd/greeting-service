package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
)

// loggingInterceptor logs all incoming gRPC requests
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	slog.Debug("Incoming gRPC request",  "method", info.FullMethod, "request", fmt.Sprintf("%+v", req))

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
