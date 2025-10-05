package server

import (
	"context"
	"testing"

	pb "github.com/ayushgala/tinkerdb/proto"
)

func TestKVStoreServer_Set(t *testing.T) {
	server := NewKVStoreServer()
	ctx := context.Background()

	// Test successful set
	resp, err := server.Set(ctx, &pb.SetRequest{
		TenantId: "test-tenant",
		Key:      "test-key",
		Value:    []byte("test-value"),
	})

	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	if !resp.Success {
		t.Fatalf("Expected success, got: %s", resp.Message)
	}

	// Test with empty tenant ID
	resp, err = server.Set(ctx, &pb.SetRequest{
		TenantId: "",
		Key:      "key",
		Value:    []byte("value"),
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.Success {
		t.Fatal("Expected failure for empty tenant ID")
	}

	// Test with empty key
	resp, err = server.Set(ctx, &pb.SetRequest{
		TenantId: "tenant",
		Key:      "",
		Value:    []byte("value"),
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.Success {
		t.Fatal("Expected failure for empty key")
	}
}

func TestKVStoreServer_Get(t *testing.T) {
	server := NewKVStoreServer()
	ctx := context.Background()

	tenantID := "test-tenant"
	key := "test-key"
	value := []byte("test-value")

	// Set a value first
	server.Set(ctx, &pb.SetRequest{
		TenantId: tenantID,
		Key:      key,
		Value:    value,
	})

	// Test successful get
	resp, err := server.Get(ctx, &pb.GetRequest{
		TenantId: tenantID,
		Key:      key,
	})

	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if !resp.Found {
		t.Fatal("Expected to find key")
	}

	if string(resp.Value) != string(value) {
		t.Fatalf("Expected value %s, got %s", value, resp.Value)
	}

	// Test get nonexistent key
	resp, err = server.Get(ctx, &pb.GetRequest{
		TenantId: tenantID,
		Key:      "nonexistent",
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.Found {
		t.Fatal("Should not find nonexistent key")
	}

	// Test with empty tenant ID
	resp, err = server.Get(ctx, &pb.GetRequest{
		TenantId: "",
		Key:      key,
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.Found {
		t.Fatal("Should not find key with empty tenant ID")
	}
}

func TestKVStoreServer_Delete(t *testing.T) {
	server := NewKVStoreServer()
	ctx := context.Background()

	tenantID := "test-tenant"
	key := "test-key"

	// Set a value first
	server.Set(ctx, &pb.SetRequest{
		TenantId: tenantID,
		Key:      key,
		Value:    []byte("value"),
	})

	// Test successful delete
	resp, err := server.Delete(ctx, &pb.DeleteRequest{
		TenantId: tenantID,
		Key:      key,
	})

	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if !resp.Success {
		t.Fatalf("Expected success, got: %s", resp.Message)
	}

	// Verify key is deleted
	getResp, _ := server.Get(ctx, &pb.GetRequest{
		TenantId: tenantID,
		Key:      key,
	})

	if getResp.Found {
		t.Fatal("Key should not exist after deletion")
	}

	// Test delete nonexistent key
	resp, err = server.Delete(ctx, &pb.DeleteRequest{
		TenantId: tenantID,
		Key:      "nonexistent",
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.Success {
		t.Fatal("Should not successfully delete nonexistent key")
	}

	// Test with empty tenant ID
	resp, err = server.Delete(ctx, &pb.DeleteRequest{
		TenantId: "",
		Key:      key,
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.Success {
		t.Fatal("Should not delete with empty tenant ID")
	}
}

func TestKVStoreServer_Exists(t *testing.T) {
	server := NewKVStoreServer()
	ctx := context.Background()

	tenantID := "test-tenant"
	key := "test-key"

	// Check nonexistent key
	resp, err := server.Exists(ctx, &pb.ExistsRequest{
		TenantId: tenantID,
		Key:      key,
	})

	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}

	if resp.Exists {
		t.Fatal("Key should not exist initially")
	}

	// Set the key
	server.Set(ctx, &pb.SetRequest{
		TenantId: tenantID,
		Key:      key,
		Value:    []byte("value"),
	})

	// Check existing key
	resp, err = server.Exists(ctx, &pb.ExistsRequest{
		TenantId: tenantID,
		Key:      key,
	})

	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}

	if !resp.Exists {
		t.Fatal("Key should exist after Set")
	}

	// Test with empty tenant ID
	resp, err = server.Exists(ctx, &pb.ExistsRequest{
		TenantId: "",
		Key:      key,
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.Exists {
		t.Fatal("Should return false for empty tenant ID")
	}
}

func TestKVStoreServer_Keys(t *testing.T) {
	server := NewKVStoreServer()
	ctx := context.Background()

	tenantID := "test-tenant"

	// Test empty tenant
	resp, err := server.Keys(ctx, &pb.KeysRequest{
		TenantId: tenantID,
	})

	if err != nil {
		t.Fatalf("Keys failed: %v", err)
	}

	if len(resp.Keys) != 0 {
		t.Fatal("Expected no keys for new tenant")
	}

	// Add some keys
	keys := []string{"key1", "key2", "key3"}
	for _, key := range keys {
		server.Set(ctx, &pb.SetRequest{
			TenantId: tenantID,
			Key:      key,
			Value:    []byte("value"),
		})
	}

	// Get all keys
	resp, err = server.Keys(ctx, &pb.KeysRequest{
		TenantId: tenantID,
	})

	if err != nil {
		t.Fatalf("Keys failed: %v", err)
	}

	if len(resp.Keys) != len(keys) {
		t.Fatalf("Expected %d keys, got %d", len(keys), len(resp.Keys))
	}

	// Test with empty tenant ID
	resp, err = server.Keys(ctx, &pb.KeysRequest{
		TenantId: "",
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(resp.Keys) != 0 {
		t.Fatal("Should return empty slice for empty tenant ID")
	}
}

func TestKVStoreServer_TenantIsolation(t *testing.T) {
	server := NewKVStoreServer()
	ctx := context.Background()

	tenant1 := "tenant1"
	tenant2 := "tenant2"
	key := "shared-key"

	// Set same key for different tenants
	server.Set(ctx, &pb.SetRequest{
		TenantId: tenant1,
		Key:      key,
		Value:    []byte("tenant1-value"),
	})

	server.Set(ctx, &pb.SetRequest{
		TenantId: tenant2,
		Key:      key,
		Value:    []byte("tenant2-value"),
	})

	// Verify isolation
	resp1, _ := server.Get(ctx, &pb.GetRequest{
		TenantId: tenant1,
		Key:      key,
	})

	resp2, _ := server.Get(ctx, &pb.GetRequest{
		TenantId: tenant2,
		Key:      key,
	})

	if string(resp1.Value) != "tenant1-value" {
		t.Fatal("Tenant1 value incorrect")
	}

	if string(resp2.Value) != "tenant2-value" {
		t.Fatal("Tenant2 value incorrect")
	}

	// Delete from tenant1 shouldn't affect tenant2
	server.Delete(ctx, &pb.DeleteRequest{
		TenantId: tenant1,
		Key:      key,
	})

	exists1, _ := server.Exists(ctx, &pb.ExistsRequest{
		TenantId: tenant1,
		Key:      key,
	})

	exists2, _ := server.Exists(ctx, &pb.ExistsRequest{
		TenantId: tenant2,
		Key:      key,
	})

	if exists1.Exists {
		t.Fatal("Key should not exist in tenant1 after deletion")
	}

	if !exists2.Exists {
		t.Fatal("Key should still exist in tenant2")
	}
}

func TestKVStoreServer_BinaryData(t *testing.T) {
	server := NewKVStoreServer()
	ctx := context.Background()

	tenantID := "test-tenant"
	key := "binary-key"
	binaryValue := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE}

	// Set binary data
	server.Set(ctx, &pb.SetRequest{
		TenantId: tenantID,
		Key:      key,
		Value:    binaryValue,
	})

	// Get and verify
	resp, err := server.Get(ctx, &pb.GetRequest{
		TenantId: tenantID,
		Key:      key,
	})

	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if !resp.Found {
		t.Fatal("Key should exist")
	}

	if len(resp.Value) != len(binaryValue) {
		t.Fatalf("Expected length %d, got %d", len(binaryValue), len(resp.Value))
	}

	for i, b := range binaryValue {
		if resp.Value[i] != b {
			t.Fatalf("Binary data mismatch at index %d: expected %x, got %x", i, b, resp.Value[i])
		}
	}
}
