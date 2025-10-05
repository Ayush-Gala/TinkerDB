package storage

import (
	"fmt"
	"sync"
)

// TenantStore represents a key-value store for a single tenant
type TenantStore struct {
	data map[string][]byte
	mu   sync.RWMutex
}

// NewTenantStore creates a new tenant store
func NewTenantStore() *TenantStore {
	return &TenantStore{
		data: make(map[string][]byte),
	}
}

// Set stores a key-value pair in the tenant store
func (ts *TenantStore) Set(key string, value []byte) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	ts.mu.Lock()
	defer ts.mu.Unlock()

	// Create a copy of the value to avoid external modifications
	valueCopy := make([]byte, len(value))
	copy(valueCopy, value)
	ts.data[key] = valueCopy

	return nil
}

// Get retrieves a value for a key from the tenant store
func (ts *TenantStore) Get(key string) ([]byte, bool) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	value, exists := ts.data[key]
	if !exists {
		return nil, false
	}

	// Return a copy to prevent external modifications
	valueCopy := make([]byte, len(value))
	copy(valueCopy, value)
	return valueCopy, true
}

// Delete removes a key from the tenant store
func (ts *TenantStore) Delete(key string) bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	_, exists := ts.data[key]
	if exists {
		delete(ts.data, key)
	}
	return exists
}

// Exists checks if a key exists in the tenant store
func (ts *TenantStore) Exists(key string) bool {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	_, exists := ts.data[key]
	return exists
}

// Keys returns all keys in the tenant store
func (ts *TenantStore) Keys() []string {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	keys := make([]string, 0, len(ts.data))
	for key := range ts.data {
		keys = append(keys, key)
	}
	return keys
}

// Size returns the number of keys in the tenant store
func (ts *TenantStore) Size() int {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	return len(ts.data)
}

// Store represents the multi-tenant key-value store
type Store struct {
	tenants map[string]*TenantStore
	mu      sync.RWMutex
}

// NewStore creates a new multi-tenant store
func NewStore() *Store {
	return &Store{
		tenants: make(map[string]*TenantStore),
	}
}

// getTenantStore retrieves or creates a tenant store
func (s *Store) getTenantStore(tenantID string) *TenantStore {
	// First try with read lock for performance
	s.mu.RLock()
	tenantStore, exists := s.tenants[tenantID]
	s.mu.RUnlock()

	if exists {
		return tenantStore
	}

	// Need to create, acquire write lock
	s.mu.Lock()
	defer s.mu.Unlock()

	// Double-check after acquiring write lock (another goroutine might have created it)
	tenantStore, exists = s.tenants[tenantID]
	if exists {
		return tenantStore
	}

	// Create new tenant store
	tenantStore = NewTenantStore()
	s.tenants[tenantID] = tenantStore
	return tenantStore
}

// Set stores a key-value pair for a specific tenant
func (s *Store) Set(tenantID, key string, value []byte) error {
	if tenantID == "" {
		return fmt.Errorf("tenant ID cannot be empty")
	}

	tenantStore := s.getTenantStore(tenantID)
	return tenantStore.Set(key, value)
}

// Get retrieves a value for a key from a specific tenant
func (s *Store) Get(tenantID, key string) ([]byte, bool) {
	if tenantID == "" {
		return nil, false
	}

	s.mu.RLock()
	tenantStore, exists := s.tenants[tenantID]
	s.mu.RUnlock()

	if !exists {
		return nil, false
	}

	return tenantStore.Get(key)
}

// Delete removes a key from a specific tenant
func (s *Store) Delete(tenantID, key string) bool {
	if tenantID == "" {
		return false
	}

	s.mu.RLock()
	tenantStore, exists := s.tenants[tenantID]
	s.mu.RUnlock()

	if !exists {
		return false
	}

	return tenantStore.Delete(key)
}

// Exists checks if a key exists for a specific tenant
func (s *Store) Exists(tenantID, key string) bool {
	if tenantID == "" {
		return false
	}

	s.mu.RLock()
	tenantStore, exists := s.tenants[tenantID]
	s.mu.RUnlock()

	if !exists {
		return false
	}

	return tenantStore.Exists(key)
}

// Keys returns all keys for a specific tenant
func (s *Store) Keys(tenantID string) []string {
	if tenantID == "" {
		return []string{}
	}

	s.mu.RLock()
	tenantStore, exists := s.tenants[tenantID]
	s.mu.RUnlock()

	if !exists {
		return []string{}
	}

	return tenantStore.Keys()
}

// TenantCount returns the number of tenants in the store
func (s *Store) TenantCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.tenants)
}

// DeleteTenant removes an entire tenant and all its data
func (s *Store) DeleteTenant(tenantID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.tenants[tenantID]
	if exists {
		delete(s.tenants, tenantID)
	}
	return exists
}
