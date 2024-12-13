package main

// import (
// 	"log"
// 	"net"

// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/reflection"

// 	pb "github.com/gitops-ci-cd/schema/v1"
// )

// func main() {
// 	// Call run and handle errors
// 	if err := run(); err != nil {
// 		log.Fatalf("Failed to run the server: %v", err)
// 	}
// }

// // run sets up and starts the gRPC server
// func run() error {
// 	// Define the port
// 	const port = ":50051"

// 	// Create a TCP listener
// 	listener, err := net.Listen("tcp", port)
// 	if err != nil {
// 		return err
// 	}
// 	log.Printf("gRPC server is listening on %s", port)

// 	// Create a new gRPC server
// 	server := grpc.NewServer()

// 	// Register the Greeter service
// 	pb.RegisterGreeterServer(server, handlers.NewGreeterHandler())

// 	// Register reflection service for debugging
// 	reflection.Register(server)

// 	// Start the server
// 	return server.Serve(listener)
// }
