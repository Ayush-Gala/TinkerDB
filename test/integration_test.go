package test

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/ayushgala/tinkerdb/internal/server"
	"github.com/ayushgala/tinkerdb/pkg/client"
	pb "github.com/ayushgala/tinkerdb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterKVStoreServer(s, server.NewKVStoreServer())
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(fmt.Sprintf("Server exited with error: %v", err))
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func getTestClient(t *testing.T) *client.Client {
	// Note: This is a placeholder function for future client-level integration tests
	// Current tests use the gRPC client directly for better control
	cfg := &client.Config{
		Address:  "bufnet",
		TenantID: "test-tenant",
	}
	c, err := client.NewClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	return c
}

func TestIntegration_BasicOperations(t *testing.T) {
	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewKVStoreClient(conn)
	ctx := context.Background()
	tenantID := "test-tenant"

	// Test Set
	setResp, err := client.Set(ctx, &pb.SetRequest{
		TenantId: tenantID,
		Key:      "name",
		Value:    []byte("TinkerDB"),
	})
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if !setResp.Success {
		t.Fatalf("Set should succeed: %s", setResp.Message)
	}

	// Test Get
	getResp, err := client.Get(ctx, &pb.GetRequest{
		TenantId: tenantID,
		Key:      "name",
	})
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !getResp.Found {
		t.Fatal("Key should be found")
	}
	if string(getResp.Value) != "TinkerDB" {
		t.Fatalf("Expected 'TinkerDB', got '%s'", getResp.Value)
	}

	// Test Exists
	existsResp, err := client.Exists(ctx, &pb.ExistsRequest{
		TenantId: tenantID,
		Key:      "name",
	})
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if !existsResp.Exists {
		t.Fatal("Key should exist")
	}

	// Test Keys
	client.Set(ctx, &pb.SetRequest{
		TenantId: tenantID,
		Key:      "version",
		Value:    []byte("1.0.0"),
	})

	keysResp, err := client.Keys(ctx, &pb.KeysRequest{
		TenantId: tenantID,
	})
	if err != nil {
		t.Fatalf("Keys failed: %v", err)
	}
	if len(keysResp.Keys) != 2 {
		t.Fatalf("Expected 2 keys, got %d", len(keysResp.Keys))
	}

	// Test Delete
	deleteResp, err := client.Delete(ctx, &pb.DeleteRequest{
		TenantId: tenantID,
		Key:      "name",
	})
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if !deleteResp.Success {
		t.Fatalf("Delete should succeed: %s", deleteResp.Message)
	}

	// Verify deletion
	existsResp, err = client.Exists(ctx, &pb.ExistsRequest{
		TenantId: tenantID,
		Key:      "name",
	})
	if err != nil {
		t.Fatalf("Exists check failed: %v", err)
	}
	if existsResp.Exists {
		t.Fatal("Key should not exist after deletion")
	}
}

func TestIntegration_MultiTenant(t *testing.T) {
	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewKVStoreClient(conn)
	ctx := context.Background()

	tenant1 := "tenant1"
	tenant2 := "tenant2"
	key := "config"

	// Set different values for same key in different tenants
	client.Set(ctx, &pb.SetRequest{
		TenantId: tenant1,
		Key:      key,
		Value:    []byte("config1"),
	})

	client.Set(ctx, &pb.SetRequest{
		TenantId: tenant2,
		Key:      key,
		Value:    []byte("config2"),
	})

	// Verify tenant isolation
	resp1, err := client.Get(ctx, &pb.GetRequest{
		TenantId: tenant1,
		Key:      key,
	})
	if err != nil {
		t.Fatalf("Get tenant1 failed: %v", err)
	}
	if string(resp1.Value) != "config1" {
		t.Fatalf("Tenant1 expected 'config1', got '%s'", resp1.Value)
	}

	resp2, err := client.Get(ctx, &pb.GetRequest{
		TenantId: tenant2,
		Key:      key,
	})
	if err != nil {
		t.Fatalf("Get tenant2 failed: %v", err)
	}
	if string(resp2.Value) != "config2" {
		t.Fatalf("Tenant2 expected 'config2', got '%s'", resp2.Value)
	}
}

func TestIntegration_ConcurrentClients(t *testing.T) {
	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewKVStoreClient(conn)
	ctx := context.Background()

	numClients := 10
	numOps := 50
	done := make(chan bool, numClients)

	// Simulate multiple concurrent clients
	for i := 0; i < numClients; i++ {
		go func(clientID int) {
			tenantID := fmt.Sprintf("tenant-%d", clientID)

			for j := 0; j < numOps; j++ {
				key := fmt.Sprintf("key-%d", j)
				value := []byte(fmt.Sprintf("value-%d-%d", clientID, j))

				// Set
				_, err := client.Set(ctx, &pb.SetRequest{
					TenantId: tenantID,
					Key:      key,
					Value:    value,
				})
				if err != nil {
					t.Errorf("Set failed: %v", err)
					done <- false
					return
				}

				// Get
				_, err = client.Get(ctx, &pb.GetRequest{
					TenantId: tenantID,
					Key:      key,
				})
				if err != nil {
					t.Errorf("Get failed: %v", err)
					done <- false
					return
				}
			}

			done <- true
		}(i)
	}

	// Wait for all goroutines
	success := 0
	for i := 0; i < numClients; i++ {
		if <-done {
			success++
		}
	}

	if success != numClients {
		t.Fatalf("Expected %d successful clients, got %d", numClients, success)
	}
}

func TestIntegration_LargeValues(t *testing.T) {
	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewKVStoreClient(conn)
	ctx := context.Background()
	tenantID := "test-tenant"

	// Create a 1MB value
	largeValue := make([]byte, 1024*1024)
	for i := range largeValue {
		largeValue[i] = byte(i % 256)
	}

	// Set large value
	setResp, err := client.Set(ctx, &pb.SetRequest{
		TenantId: tenantID,
		Key:      "large-data",
		Value:    largeValue,
	})
	if err != nil {
		t.Fatalf("Set large value failed: %v", err)
	}
	if !setResp.Success {
		t.Fatal("Set large value should succeed")
	}

	// Get large value
	getResp, err := client.Get(ctx, &pb.GetRequest{
		TenantId: tenantID,
		Key:      "large-data",
	})
	if err != nil {
		t.Fatalf("Get large value failed: %v", err)
	}
	if !getResp.Found {
		t.Fatal("Large value should be found")
	}

	// Verify data integrity
	if len(getResp.Value) != len(largeValue) {
		t.Fatalf("Expected length %d, got %d", len(largeValue), len(getResp.Value))
	}

	for i := range largeValue {
		if getResp.Value[i] != largeValue[i] {
			t.Fatalf("Data mismatch at index %d", i)
		}
	}
}

func BenchmarkIntegration_Set(b *testing.B) {
	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		b.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewKVStoreClient(conn)
	ctx := context.Background()
	tenantID := "bench-tenant"
	value := []byte("benchmark-value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		client.Set(ctx, &pb.SetRequest{
			TenantId: tenantID,
			Key:      key,
			Value:    value,
		})
	}
}

func BenchmarkIntegration_Get(b *testing.B) {
	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		b.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewKVStoreClient(conn)
	ctx := context.Background()
	tenantID := "bench-tenant"

	// Prepare data
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key-%d", i)
		client.Set(ctx, &pb.SetRequest{
			TenantId: tenantID,
			Key:      key,
			Value:    []byte("value"),
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i%1000)
		client.Get(ctx, &pb.GetRequest{
			TenantId: tenantID,
			Key:      key,
		})
	}
}
