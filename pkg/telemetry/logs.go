package telemetry

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// loggingInterceptor logs all incoming gRPC requests
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		msg := req.(protoreflect.ProtoMessage)
		fields := make(map[string]interface{})
		msg.ProtoReflect().Range(func(fd protoreflect.FieldDescriptor, value protoreflect.Value) bool {
			// Only include fields that are set
			if value.IsValid() {
				fields[string(fd.Name())] = value.Interface()
			}
			return true
		})

		fieldsJSON, err := json.Marshal(fields)
		if err != nil {
			slog.Error("Failed to marshal request fields", "error", err)
		} else {
			slog.Debug("Incoming gRPC request", "method", info.FullMethod, "type", fmt.Sprintf("%T", req), "request", string(fieldsJSON))
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
