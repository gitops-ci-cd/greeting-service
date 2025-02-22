package io

import (
	"context"
	"net/http"
	"time"
)

// HTTPServer is a wrapper around http.Server that implements the Server interface
type HTTPServer struct {
	*http.Server
}

// GracefulStop shuts down the server gracefully
func (h *HTTPServer) GracefulStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	h.Shutdown(ctx)
}

// NewHTTPServer creates a new HTTPServer instance
func NewHTTPServer(handler http.Handler) Server {
	return &HTTPServer{
		Server: &http.Server{
			Handler: handler,
		},
	}
}
