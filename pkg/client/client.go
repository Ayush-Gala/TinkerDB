package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/ayushgala/tinkerdb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client represents a TinkerDB client
type Client struct {
	conn     *grpc.ClientConn
	client   pb.KVStoreClient
	tenantID string
}

// Config holds client configuration
type Config struct {
	Address  string
	TenantID string
	Timeout  time.Duration
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Address:  "localhost:50051",
		TenantID: "default",
		Timeout:  5 * time.Second,
	}
}

// NewClient creates a new TinkerDB client
func NewClient(cfg *Config) (*Client, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	// Create gRPC connection
	conn, err := grpc.NewClient(cfg.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	client := pb.NewKVStoreClient(conn)

	return &Client{
		conn:     conn,
		client:   client,
		tenantID: cfg.TenantID,
	}, nil
}

// Close closes the client connection
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Set stores a key-value pair
func (c *Client) Set(ctx context.Context, key string, value []byte) error {
	resp, err := c.client.Set(ctx, &pb.SetRequest{
		TenantId: c.tenantID,
		Key:      key,
		Value:    value,
	})
	if err != nil {
		return fmt.Errorf("set failed: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("set failed: %s", resp.Message)
	}

	return nil
}

// SetString stores a key with a string value
func (c *Client) SetString(ctx context.Context, key, value string) error {
	return c.Set(ctx, key, []byte(value))
}

// Get retrieves a value for a key
func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	resp, err := c.client.Get(ctx, &pb.GetRequest{
		TenantId: c.tenantID,
		Key:      key,
	})
	if err != nil {
		return nil, fmt.Errorf("get failed: %w", err)
	}

	if !resp.Found {
		return nil, fmt.Errorf("key not found: %s", key)
	}

	return resp.Value, nil
}

// GetString retrieves a string value for a key
func (c *Client) GetString(ctx context.Context, key string) (string, error) {
	value, err := c.Get(ctx, key)
	if err != nil {
		return "", err
	}
	return string(value), nil
}

// Delete removes a key
func (c *Client) Delete(ctx context.Context, key string) error {
	resp, err := c.client.Delete(ctx, &pb.DeleteRequest{
		TenantId: c.tenantID,
		Key:      key,
	})
	if err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("delete failed: %s", resp.Message)
	}

	return nil
}

// Exists checks if a key exists
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	resp, err := c.client.Exists(ctx, &pb.ExistsRequest{
		TenantId: c.tenantID,
		Key:      key,
	})
	if err != nil {
		return false, fmt.Errorf("exists check failed: %w", err)
	}

	return resp.Exists, nil
}

// Keys retrieves all keys in the tenant namespace
func (c *Client) Keys(ctx context.Context) ([]string, error) {
	resp, err := c.client.Keys(ctx, &pb.KeysRequest{
		TenantId: c.tenantID,
	})
	if err != nil {
		return nil, fmt.Errorf("keys retrieval failed: %w", err)
	}

	return resp.Keys, nil
}

// SetTenant changes the tenant ID for subsequent operations
func (c *Client) SetTenant(tenantID string) {
	c.tenantID = tenantID
}

// GetTenant returns the current tenant ID
func (c *Client) GetTenant() string {
	return c.tenantID
}
