# Consistent Hash vs Modulo Comparison Demo

This demo compares consistent hashing with naive modulo hashing to demonstrate the advantages when scaling distributed systems.

## What it does

- Creates 100 keys distributed across 3 servers
- Compares key redistribution when adding a 4th server:
  - **Modulo hashing**: Keys are redistributed using `key % num_servers`
  - **Consistent hashing**: Keys are redistributed using a hash ring
- Shows how consistent hashing minimizes key movement

## Run the demo

```bash
task demo:compare
```

## Key concepts

- **Minimal redistribution**: Consistent hashing moves significantly fewer keys when servers are added/removed
- **Scalability**: Critical for distributed systems that need to scale without massive cache invalidation
- **Efficiency**: Typically moves only ~1/N keys when adding a server to N servers (vs ~N-1/N for modulo)