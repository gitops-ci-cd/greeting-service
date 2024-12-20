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
		// Check the type of req
		slog.Debug("Inspecting req type", "type", fmt.Sprintf("%T", req))

		// Log raw req value
		slog.Debug("Raw req value", "value", req)

		// Attempt reflection if req is a ProtoMessage
		if msg, ok := req.(protoreflect.ProtoMessage); ok {
			reqFields := make(map[string]interface{})
			msg.ProtoReflect().Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
				slog.Debug("Inspecting field", "name", fd.Name(), "value", v.Interface())
				reqFields[string(fd.Name())] = v.Interface()
				return true
			})

			reqDetailsBytes, err := json.Marshal(reqFields)
			if err != nil {
				slog.Error("Error marshaling request fields", "method", info.FullMethod, "error", err)
			} else if slog.Default().Enabled(ctx, slog.LevelDebug) {
				slog.Debug("Incoming gRPC request", "method", info.FullMethod, "request", string(reqDetailsBytes))
			}
		} else {
			slog.Warn("Req does not implement protoreflect.ProtoMessage", "type", fmt.Sprintf("%T", req))
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
