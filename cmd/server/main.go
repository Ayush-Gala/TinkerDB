package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/ayushgala/tinkerdb/internal/server"
	pb "github.com/ayushgala/tinkerdb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	defaultPort = "8080"
)

func main() {
	port := os.Getenv("TINKERDB_PORT")
	if port == "" {
		port = defaultPort
	}

	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register KVStore service
	kvStoreServer := server.NewKVStoreServer()
	pb.RegisterKVStoreServer(grpcServer, kvStoreServer)

	// Register reflection service for debugging with tools like grpcurl
	reflection.Register(grpcServer)

	// Set up signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf("TinkerDB server starting on port %s...", port)
		log.Printf("Server is ready to accept connections")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for termination signal
	<-sigCh
	log.Println("\nShutting down server gracefully...")
	grpcServer.GracefulStop()
	log.Println("Server stopped")
}
