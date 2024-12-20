package telemetry

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// loggingInterceptor logs all incoming gRPC requests
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		reqFields := make(map[string]interface{})
		reqReflect := req.(protoreflect.ProtoMessage).ProtoReflect()

		reqReflect.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
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
