package io

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// TimestampInjector is a gRPC middleware that injects the current timestamp into the metadata of each request
func TimestampInjector(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Generate the timestamp
	timestamp := time.Now().Format(time.RFC3339)

	// Add the timestamp to the metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	md = metadata.Join(md, metadata.Pairs("timestamp", timestamp))

	// Create a new context with the updated metadata
	ctx = metadata.NewIncomingContext(ctx, md)

	// Call the handler with the new context
	return handler(ctx, req)
}
