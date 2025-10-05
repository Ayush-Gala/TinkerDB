package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ayushgala/tinkerdb/pkg/client"
)

func main() {
	// Connect to TinkerDB
	cfg := &client.Config{
		Address:  "localhost:50051",
		TenantID: "interactive",
	}

	c, err := client.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to TinkerDB: %v", err)
	}
	defer c.Close()

	scanner := bufio.NewScanner(os.Stdin)
	ctx := context.Background()

	fmt.Println("╔════════════════════════════════════════════════════════╗")
	fmt.Println("║        TinkerDB Interactive Client v1.0                ║")
	fmt.Println("╚════════════════════════════════════════════════════════╝")
	fmt.Printf("\nConnected to: %s\n", cfg.Address)
	fmt.Printf("Current tenant: %s\n", c.GetTenant())
	fmt.Println("\nCommands:")
	fmt.Println("  set <key> <value>  - Set a key-value pair")
	fmt.Println("  get <key>          - Get value for a key")
	fmt.Println("  delete <key>       - Delete a key")
	fmt.Println("  exists <key>       - Check if key exists")
	fmt.Println("  keys               - List all keys")
	fmt.Println("  tenant <id>        - Switch tenant (or show current)")
	fmt.Println("  help               - Show this help")
	fmt.Println("  quit               - Exit")
	fmt.Println()

	for {
		fmt.Printf("[%s]> ", c.GetTenant())
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		cmd := parts[0]

		switch cmd {
		case "set":
			if len(parts) < 3 {
				fmt.Println("❌ Usage: set <key> <value>")
				continue
			}
			key := parts[1]
			value := strings.Join(parts[2:], " ")
			err := c.SetString(ctx, key, value)
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			} else {
				fmt.Printf("✓ Set '%s' = '%s'\n", key, value)
			}

		case "get":
			if len(parts) < 2 {
				fmt.Println("❌ Usage: get <key>")
				continue
			}
			key := parts[1]
			value, err := c.GetString(ctx, key)
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			} else {
				fmt.Printf("✓ %s = '%s'\n", key, value)
			}

		case "delete":
			if len(parts) < 2 {
				fmt.Println("❌ Usage: delete <key>")
				continue
			}
			key := parts[1]
			err := c.Delete(ctx, key)
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			} else {
				fmt.Printf("✓ Deleted '%s'\n", key)
			}

		case "exists":
			if len(parts) < 2 {
				fmt.Println("❌ Usage: exists <key>")
				continue
			}
			key := parts[1]
			exists, err := c.Exists(ctx, key)
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			} else {
				if exists {
					fmt.Printf("✓ Key '%s' exists\n", key)
				} else {
					fmt.Printf("✗ Key '%s' does not exist\n", key)
				}
			}

		case "keys":
			keys, err := c.Keys(ctx)
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			} else {
				if len(keys) == 0 {
					fmt.Println("No keys found")
				} else {
					fmt.Printf("✓ Found %d key(s):\n", len(keys))
					for i, k := range keys {
						fmt.Printf("  %d. %s\n", i+1, k)
					}
				}
			}

		case "tenant":
			if len(parts) < 2 {
				fmt.Printf("Current tenant: %s\n", c.GetTenant())
			} else {
				c.SetTenant(parts[1])
				fmt.Printf("✓ Switched to tenant '%s'\n", parts[1])
			}

		case "help":
			fmt.Println("\nCommands:")
			fmt.Println("  set <key> <value>  - Set a key-value pair")
			fmt.Println("  get <key>          - Get value for a key")
			fmt.Println("  delete <key>       - Delete a key")
			fmt.Println("  exists <key>       - Check if key exists")
			fmt.Println("  keys               - List all keys")
			fmt.Println("  tenant <id>        - Switch tenant (or show current)")
			fmt.Println("  help               - Show this help")
			fmt.Println("  quit               - Exit")
			fmt.Println()

		case "quit", "exit":
			fmt.Println("\n👋 Goodbye!")
			return

		default:
			fmt.Printf("❌ Unknown command: %s (type 'help' for commands)\n", cmd)
		}
	}
}
