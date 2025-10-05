package storage

import (
	"sync"
	"testing"
)

func TestTenantStore_SetAndGet(t *testing.T) {
	ts := NewTenantStore()

	// Test basic set and get
	key := "test-key"
	value := []byte("test-value")
	err := ts.Set(key, value)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	retrievedValue, exists := ts.Get(key)
	if !exists {
		t.Fatal("Key should exist after Set")
	}

	if string(retrievedValue) != string(value) {
		t.Fatalf("Expected value %s, got %s", value, retrievedValue)
	}
}

func TestTenantStore_SetEmptyKey(t *testing.T) {
	ts := NewTenantStore()

	err := ts.Set("", []byte("value"))
	if err == nil {
		t.Fatal("Expected error when setting empty key")
	}
}

func TestTenantStore_GetNonexistentKey(t *testing.T) {
	ts := NewTenantStore()

	_, exists := ts.Get("nonexistent")
	if exists {
		t.Fatal("Nonexistent key should not exist")
	}
}

func TestTenantStore_Delete(t *testing.T) {
	ts := NewTenantStore()

	// Set a key
	key := "test-key"
	ts.Set(key, []byte("value"))

	// Delete the key
	deleted := ts.Delete(key)
	if !deleted {
		t.Fatal("Delete should return true for existing key")
	}

	// Verify it's deleted
	_, exists := ts.Get(key)
	if exists {
		t.Fatal("Key should not exist after deletion")
	}

	// Try deleting again
	deleted = ts.Delete(key)
	if deleted {
		t.Fatal("Delete should return false for nonexistent key")
	}
}

func TestTenantStore_Exists(t *testing.T) {
	ts := NewTenantStore()

	key := "test-key"

	// Should not exist initially
	if ts.Exists(key) {
		t.Fatal("Key should not exist initially")
	}

	// Set the key
	ts.Set(key, []byte("value"))

	// Should exist now
	if !ts.Exists(key) {
		t.Fatal("Key should exist after Set")
	}

	// Delete the key
	ts.Delete(key)

	// Should not exist after deletion
	if ts.Exists(key) {
		t.Fatal("Key should not exist after deletion")
	}
}

func TestTenantStore_Keys(t *testing.T) {
	ts := NewTenantStore()

	// Empty store should have no keys
	keys := ts.Keys()
	if len(keys) != 0 {
		t.Fatalf("Expected 0 keys, got %d", len(keys))
	}

	// Add some keys
	expectedKeys := map[string]bool{
		"key1": true,
		"key2": true,
		"key3": true,
	}

	for key := range expectedKeys {
		ts.Set(key, []byte("value"))
	}

	// Check all keys are returned
	keys = ts.Keys()
	if len(keys) != len(expectedKeys) {
		t.Fatalf("Expected %d keys, got %d", len(expectedKeys), len(keys))
	}

	for _, key := range keys {
		if !expectedKeys[key] {
			t.Fatalf("Unexpected key: %s", key)
		}
	}
}

func TestTenantStore_Size(t *testing.T) {
	ts := NewTenantStore()

	if ts.Size() != 0 {
		t.Fatalf("Expected size 0, got %d", ts.Size())
	}

	ts.Set("key1", []byte("value1"))
	if ts.Size() != 1 {
		t.Fatalf("Expected size 1, got %d", ts.Size())
	}

	ts.Set("key2", []byte("value2"))
	if ts.Size() != 2 {
		t.Fatalf("Expected size 2, got %d", ts.Size())
	}

	ts.Delete("key1")
	if ts.Size() != 1 {
		t.Fatalf("Expected size 1 after deletion, got %d", ts.Size())
	}
}

func TestTenantStore_ConcurrentAccess(t *testing.T) {
	ts := NewTenantStore()

	var wg sync.WaitGroup
	numGoroutines := 100
	numOperations := 100

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := "key"
				value := []byte("value")
				ts.Set(key, value)
			}
		}(i)
	}

	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				ts.Get("key")
			}
		}()
	}

	wg.Wait()
}

func TestStore_MultiTenant(t *testing.T) {
	store := NewStore()

	// Set values for different tenants
	tenant1 := "tenant1"
	tenant2 := "tenant2"
	key := "shared-key"

	err := store.Set(tenant1, key, []byte("tenant1-value"))
	if err != nil {
		t.Fatalf("Set failed for tenant1: %v", err)
	}

	err = store.Set(tenant2, key, []byte("tenant2-value"))
	if err != nil {
		t.Fatalf("Set failed for tenant2: %v", err)
	}

	// Get values and verify isolation
	value1, exists := store.Get(tenant1, key)
	if !exists || string(value1) != "tenant1-value" {
		t.Fatal("Tenant1 value incorrect")
	}

	value2, exists := store.Get(tenant2, key)
	if !exists || string(value2) != "tenant2-value" {
		t.Fatal("Tenant2 value incorrect")
	}
}

func TestStore_EmptyTenantID(t *testing.T) {
	store := NewStore()

	// Test Set with empty tenant ID
	err := store.Set("", "key", []byte("value"))
	if err == nil {
		t.Fatal("Expected error for empty tenant ID")
	}

	// Test Get with empty tenant ID
	_, exists := store.Get("", "key")
	if exists {
		t.Fatal("Should not find key with empty tenant ID")
	}

	// Test Delete with empty tenant ID
	deleted := store.Delete("", "key")
	if deleted {
		t.Fatal("Should not delete with empty tenant ID")
	}

	// Test Exists with empty tenant ID
	if store.Exists("", "key") {
		t.Fatal("Should return false for empty tenant ID")
	}

	// Test Keys with empty tenant ID
	keys := store.Keys("")
	if len(keys) != 0 {
		t.Fatal("Should return empty slice for empty tenant ID")
	}
}

func TestStore_TenantCount(t *testing.T) {
	store := NewStore()

	if store.TenantCount() != 0 {
		t.Fatalf("Expected 0 tenants, got %d", store.TenantCount())
	}

	store.Set("tenant1", "key", []byte("value"))
	if store.TenantCount() != 1 {
		t.Fatalf("Expected 1 tenant, got %d", store.TenantCount())
	}

	store.Set("tenant2", "key", []byte("value"))
	if store.TenantCount() != 2 {
		t.Fatalf("Expected 2 tenants, got %d", store.TenantCount())
	}

	// Setting another key in existing tenant shouldn't increase count
	store.Set("tenant1", "key2", []byte("value2"))
	if store.TenantCount() != 2 {
		t.Fatalf("Expected 2 tenants, got %d", store.TenantCount())
	}
}

func TestStore_DeleteTenant(t *testing.T) {
	store := NewStore()

	tenant := "test-tenant"
	store.Set(tenant, "key1", []byte("value1"))
	store.Set(tenant, "key2", []byte("value2"))

	// Verify tenant exists
	if store.TenantCount() != 1 {
		t.Fatal("Tenant should exist")
	}

	// Delete tenant
	deleted := store.DeleteTenant(tenant)
	if !deleted {
		t.Fatal("DeleteTenant should return true")
	}

	// Verify tenant is deleted
	if store.TenantCount() != 0 {
		t.Fatal("Tenant should be deleted")
	}

	// Try deleting again
	deleted = store.DeleteTenant(tenant)
	if deleted {
		t.Fatal("DeleteTenant should return false for nonexistent tenant")
	}
}

func TestStore_ConcurrentMultiTenant(t *testing.T) {
	store := NewStore()

	var wg sync.WaitGroup
	numTenants := 10
	numOpsPerTenant := 100

	// Concurrent operations across multiple tenants
	for i := 0; i < numTenants; i++ {
		wg.Add(1)
		go func(tenantID int) {
			defer wg.Done()
			tenant := string(rune('A' + tenantID))

			for j := 0; j < numOpsPerTenant; j++ {
				key := "key"
				value := []byte("value")

				store.Set(tenant, key, value)
				store.Get(tenant, key)
				store.Exists(tenant, key)
				store.Keys(tenant)
			}
		}(i)
	}

	wg.Wait()

	// Verify tenant count
	if store.TenantCount() != numTenants {
		t.Fatalf("Expected %d tenants, got %d", numTenants, store.TenantCount())
	}
}

func TestStore_ValueIsolation(t *testing.T) {
	store := NewStore()

	// Set a value
	originalValue := []byte("original")
	store.Set("tenant", "key", originalValue)

	// Modify the original slice
	originalValue[0] = 'X'

	// Get the value and verify it wasn't affected
	retrievedValue, _ := store.Get("tenant", "key")
	if string(retrievedValue) != "original" {
		t.Fatal("Value should be isolated from external modifications")
	}

	// Modify the retrieved slice
	retrievedValue[0] = 'Y'

	// Get again and verify it wasn't affected
	retrievedValue2, _ := store.Get("tenant", "key")
	if string(retrievedValue2) != "original" {
		t.Fatal("Value should be isolated from modifications to retrieved values")
	}
}

func BenchmarkTenantStore_Set(b *testing.B) {
	ts := NewTenantStore()
	value := []byte("benchmark-value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ts.Set("key", value)
	}
}

func BenchmarkTenantStore_Get(b *testing.B) {
	ts := NewTenantStore()
	ts.Set("key", []byte("value"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ts.Get("key")
	}
}

func BenchmarkStore_SetMultiTenant(b *testing.B) {
	store := NewStore()
	value := []byte("benchmark-value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tenantID := string(rune('A' + (i % 10)))
		store.Set(tenantID, "key", value)
	}
}

func BenchmarkStore_GetMultiTenant(b *testing.B) {
	store := NewStore()

	// Prepare data
	for i := 0; i < 10; i++ {
		tenantID := string(rune('A' + i))
		store.Set(tenantID, "key", []byte("value"))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tenantID := string(rune('A' + (i % 10)))
		store.Get(tenantID, "key")
	}
}
