# Cache Distribution Demo

This demo illustrates how consistent hashing can be used for distributing cache keys across multiple cache servers.

## What it does

- Sets up a hash ring with 3 cache servers
- Distributes various cache keys (user profiles, sessions, products, etc.) across the servers
- Demonstrates consistency by showing the same key always maps to the same server

## Run the demo

```bash
task demo:cache
```

## Key concepts

- **Consistent key mapping**: Cache keys are consistently routed to the same server
- **Balanced distribution**: Keys are distributed evenly across available cache servers
- **Predictable routing**: Essential for cache hit rates in distributed caching systems