# Load Balancer with Sticky Sessions Demo

This demo shows how consistent hashing enables sticky sessions in a load balancer, ensuring the same client always connects to the same backend server.

## What it does

- Sets up a hash ring with 3 backend servers
- Routes requests based on session IDs
- Demonstrates sticky session behavior where the same session ID always routes to the same backend

## Run the demo

```bash
task demo:loadbalancer
```

## Key concepts

- **Sticky sessions**: Same session ID always routes to the same backend server
- **Session affinity**: Important for stateful applications or when backend servers maintain session state
- **Consistent routing**: Ensures user experience continuity and reduces backend state synchronization needs