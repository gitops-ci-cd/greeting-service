package io

import (
	"testing"

	"google.golang.org/grpc"
)

// MockRegisterable is a mock implementation of the Registerable interface
type MockRegisterable struct {
	registered bool
}

// Register satisfies the Registerable interface
func (m *MockRegisterable) Register(server *grpc.Server) {
	m.registered = true
}

func TestAddService(t *testing.T) {
	tests := []struct {
		name            string
		servicesToAdd   []*MockRegisterable
		expectedService int
	}{
		{
			name:            "Add one service",
			servicesToAdd:   []*MockRegisterable{{}},
			expectedService: 1,
		},
		{
			name:            "Add multiple services",
			servicesToAdd:   []*MockRegisterable{{}, {}},
			expectedService: 2,
		},
		{
			name:            "Add no service",
			servicesToAdd:   []*MockRegisterable{},
			expectedService: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := &Registry{}

			// Add the services to the registry
			for _, service := range tt.servicesToAdd {
				registry.Add(service)
			}

			// Validate the number of services added
			if len(registry.Services) != tt.expectedService {
				t.Errorf("expected %d services, got %d", tt.expectedService, len(registry.Services))
			}
		})
	}
}

func TestRegisterAllServices(t *testing.T) {
	tests := []struct {
		name          string
		servicesToAdd []*MockRegisterable
	}{
		{
			name:          "Register one service",
			servicesToAdd: []*MockRegisterable{{}},
		},
		{
			name:          "Register multiple services",
			servicesToAdd: []*MockRegisterable{{}, {}},
		},
		{
			name:          "Register no services",
			servicesToAdd: []*MockRegisterable{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := &Registry{}
			server := grpc.NewServer()

			// Add services to the registry
			for _, service := range tt.servicesToAdd {
				registry.Add(service)
			}

			// Register all services
			registry.RegisterAll(server)

			// Verify all services in the test case were registered
			for i, service := range tt.servicesToAdd {
				if !service.registered {
					t.Errorf("service %d in test case %q was not registered", i, tt.name)
				}
			}
		})
	}
}
