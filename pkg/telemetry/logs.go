package telemetry

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// loggingInterceptor logs all incoming gRPC requests
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		if msg, ok := req.(proto.Message); ok {
			// Use ProtoReflect for detailed introspection
			reqReflect := msg.ProtoReflect()
			fields := make(map[string]interface{})
			reqReflect.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
				fields[string(fd.Name())] = v.Interface()
				return true
			})
			reqDetailsBytes, err := json.Marshal(fields)
			if err != nil {
				slog.Error("Error marshaling request fields", "method", info.FullMethod, "error", err)
			} else {
				slog.Debug("Incoming gRPC request", "method", info.FullMethod, "request", string(reqDetailsBytes))
			}
		} else {
			// Fallback for non-Protobuf types
			slog.Debug("Incoming gRPC request", "method", info.FullMethod, "request", fmt.Sprintf("%+v", req))
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
