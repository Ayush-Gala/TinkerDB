package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ayushgala/tinkerdb/pkg/client"
)

func main() {
	// Create a client configuration
	cfg := &client.Config{
		Address:  "localhost:50051",
		TenantID: "example-tenant",
		Timeout:  5 * time.Second,
	}

	// Create a new client
	c, err := client.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	ctx := context.Background()

	// Example 1: Set and Get a string value
	fmt.Println("=== Example 1: Basic Set/Get ===")
	err = c.SetString(ctx, "greeting", "Hello, TinkerDB!")
	if err != nil {
		log.Printf("Set failed: %v", err)
	} else {
		fmt.Println("✓ Set key 'greeting'")
	}

	value, err := c.GetString(ctx, "greeting")
	if err != nil {
		log.Printf("Get failed: %v", err)
	} else {
		fmt.Printf("✓ Got value: %s\n", value)
	}

	// Example 2: Check if key exists
	fmt.Println("\n=== Example 2: Check Existence ===")
	exists, err := c.Exists(ctx, "greeting")
	if err != nil {
		log.Printf("Exists check failed: %v", err)
	} else {
		fmt.Printf("✓ Key 'greeting' exists: %v\n", exists)
	}

	exists, err = c.Exists(ctx, "nonexistent")
	if err != nil {
		log.Printf("Exists check failed: %v", err)
	} else {
		fmt.Printf("✓ Key 'nonexistent' exists: %v\n", exists)
	}

	// Example 3: Store multiple keys
	fmt.Println("\n=== Example 3: Multiple Keys ===")
	keys := map[string]string{
		"name":    "TinkerDB",
		"version": "1.0.0",
		"author":  "Ayush Gala",
	}

	for k, v := range keys {
		err = c.SetString(ctx, k, v)
		if err != nil {
			log.Printf("Set failed for key '%s': %v", k, err)
		} else {
			fmt.Printf("✓ Set key '%s' = '%s'\n", k, v)
		}
	}

	// Example 4: List all keys
	fmt.Println("\n=== Example 4: List All Keys ===")
	allKeys, err := c.Keys(ctx)
	if err != nil {
		log.Printf("Keys retrieval failed: %v", err)
	} else {
		fmt.Printf("✓ Total keys: %d\n", len(allKeys))
		for _, k := range allKeys {
			val, _ := c.GetString(ctx, k)
			fmt.Printf("  - %s: %s\n", k, val)
		}
	}

	// Example 5: Delete a key
	fmt.Println("\n=== Example 5: Delete Key ===")
	err = c.Delete(ctx, "greeting")
	if err != nil {
		log.Printf("Delete failed: %v", err)
	} else {
		fmt.Println("✓ Deleted key 'greeting'")
	}

	exists, _ = c.Exists(ctx, "greeting")
	fmt.Printf("✓ Key 'greeting' exists after deletion: %v\n", exists)

	// Example 6: Multi-tenant support
	fmt.Println("\n=== Example 6: Multi-Tenant Support ===")

	// Set a key in current tenant
	c.SetString(ctx, "tenant-specific", "tenant-1-value")
	fmt.Printf("✓ Set key in tenant '%s'\n", c.GetTenant())

	// Switch to different tenant
	c.SetTenant("another-tenant")
	c.SetString(ctx, "tenant-specific", "tenant-2-value")
	fmt.Printf("✓ Set same key in tenant '%s'\n", c.GetTenant())

	// Retrieve from both tenants
	c.SetTenant("example-tenant")
	val1, _ := c.GetString(ctx, "tenant-specific")
	fmt.Printf("✓ Value in 'example-tenant': %s\n", val1)

	c.SetTenant("another-tenant")
	val2, _ := c.GetString(ctx, "tenant-specific")
	fmt.Printf("✓ Value in 'another-tenant': %s\n", val2)

	// Example 7: Binary data
	fmt.Println("\n=== Example 7: Binary Data ===")
	binaryData := []byte{0x00, 0x01, 0x02, 0x03, 0xFF}
	err = c.Set(ctx, "binary-key", binaryData)
	if err != nil {
		log.Printf("Set binary data failed: %v", err)
	} else {
		fmt.Println("✓ Stored binary data")
	}

	retrieved, err := c.Get(ctx, "binary-key")
	if err != nil {
		log.Printf("Get binary data failed: %v", err)
	} else {
		fmt.Printf("✓ Retrieved binary data: %v\n", retrieved)
	}

	fmt.Println("\n=== All examples completed successfully! ===")
}
