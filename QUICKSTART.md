# TinkerDB Quick Start Guide

Great! Before moving on to the next steps, make sure you have set up the reporsitory and installed all the required packages and softwares. You are all set to start using TinkerDB. 🚀

## Step 1: Start the Server (Terminal 1)

From the root of the repository run:
```bash
go run cmd/server/main.go
```

✅ You should see:
```
TinkerDB server starting on port 50051...
Server is ready to accept connections
```

**Keep this terminal running!**

## Step 2: Run your queries (Terminal 2)

### Option A: Run the Example (Easiest)

From the root directory, run:
```bash
go run examples/basic_usage/main.go
```

### Option B: Interactive Client (Most Fun!)

Launch the TinkerDB CLI:
```bash
go run interactive_client.go
```

Try these commands:
```
[interactive]> set name <your name>
[interactive]> get name
[interactive]> keys
[interactive]> tenant app2
[app2]> set name "Different Ayush"
[app2]> get name
[app2]> quit
```

## Next Steps

**Full testing guide**: See `MANUAL_TESTING.md`   
**Complete docs**: See `docs/Milestone1.md`  

## Stop the Server

Press `Ctrl+C` in the server terminal to gracefully shut down.
