package io

import (
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Registerable is an interface for handlers to implement so they can be registered with a server
type Registerable interface {
	Register(server *grpc.Server)
}

// Registry is a collection of services that can be registered with a server
type Registry struct {
	Services []Registerable
}

// Add appends a service to the registry
func (r *Registry) Add(service Registerable) {
	r.Services = append(r.Services, service)
}

// RegisterAll registers all services in the registry with the given server
func (r *Registry) RegisterAll(server *grpc.Server) {
	for _, service := range r.Services {
		service.Register(server)
	}

	// Register reflection service for debugging
	reflection.Register(server)

	for key, value := range server.GetServiceInfo() {
		slog.Info("Service registered", "service", key, "methods", value.Methods, "metadata", value.Metadata)
	}
}
