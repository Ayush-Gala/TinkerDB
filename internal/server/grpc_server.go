package server

import (
	"context"
	"fmt"
	"log"

	"github.com/ayushgala/tinkerdb/internal/storage"
	pb "github.com/ayushgala/tinkerdb/proto"
)

// KVStoreServer implements the gRPC KVStore service
type KVStoreServer struct {
	pb.UnimplementedKVStoreServer
	store *storage.Store
}

// NewKVStoreServer creates a new gRPC server instance
func NewKVStoreServer() *KVStoreServer {
	return &KVStoreServer{
		store: storage.NewStore(),
	}
}

// Set implements the Set RPC method
func (s *KVStoreServer) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	log.Printf("Set: tenant=%s, key=%s, value_size=%d bytes", req.TenantId, req.Key, len(req.Value))

	if req.TenantId == "" {
		return &pb.SetResponse{
			Success: false,
			Message: "tenant ID cannot be empty",
		}, nil
	}

	if req.Key == "" {
		return &pb.SetResponse{
			Success: false,
			Message: "key cannot be empty",
		}, nil
	}

	err := s.store.Set(req.TenantId, req.Key, req.Value)
	if err != nil {
		return &pb.SetResponse{
			Success: false,
			Message: fmt.Sprintf("failed to set key: %v", err),
		}, nil
	}

	return &pb.SetResponse{
		Success: true,
		Message: "key set successfully",
	}, nil
}

// Get implements the Get RPC method
func (s *KVStoreServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log.Printf("Get: tenant=%s, key=%s", req.TenantId, req.Key)

	if req.TenantId == "" {
		return &pb.GetResponse{
			Found:   false,
			Message: "tenant ID cannot be empty",
		}, nil
	}

	if req.Key == "" {
		return &pb.GetResponse{
			Found:   false,
			Message: "key cannot be empty",
		}, nil
	}

	value, found := s.store.Get(req.TenantId, req.Key)
	if !found {
		return &pb.GetResponse{
			Found:   false,
			Message: "key not found",
		}, nil
	}

	return &pb.GetResponse{
		Found:   true,
		Value:   value,
		Message: "key found",
	}, nil
}

// Delete implements the Delete RPC method
func (s *KVStoreServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	log.Printf("Delete: tenant=%s, key=%s", req.TenantId, req.Key)

	if req.TenantId == "" {
		return &pb.DeleteResponse{
			Success: false,
			Message: "tenant ID cannot be empty",
		}, nil
	}

	if req.Key == "" {
		return &pb.DeleteResponse{
			Success: false,
			Message: "key cannot be empty",
		}, nil
	}

	deleted := s.store.Delete(req.TenantId, req.Key)
	if !deleted {
		return &pb.DeleteResponse{
			Success: false,
			Message: "key not found",
		}, nil
	}

	return &pb.DeleteResponse{
		Success: true,
		Message: "key deleted successfully",
	}, nil
}

// Exists implements the Exists RPC method
func (s *KVStoreServer) Exists(ctx context.Context, req *pb.ExistsRequest) (*pb.ExistsResponse, error) {
	log.Printf("Exists: tenant=%s, key=%s", req.TenantId, req.Key)

	if req.TenantId == "" || req.Key == "" {
		return &pb.ExistsResponse{
			Exists: false,
		}, nil
	}

	exists := s.store.Exists(req.TenantId, req.Key)
	return &pb.ExistsResponse{
		Exists: exists,
	}, nil
}

// Keys implements the Keys RPC method
func (s *KVStoreServer) Keys(ctx context.Context, req *pb.KeysRequest) (*pb.KeysResponse, error) {
	log.Printf("Keys: tenant=%s", req.TenantId)

	if req.TenantId == "" {
		return &pb.KeysResponse{
			Keys: []string{},
		}, nil
	}

	keys := s.store.Keys(req.TenantId)
	return &pb.KeysResponse{
		Keys: keys,
	}, nil
}
