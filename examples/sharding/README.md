# Database Sharding Demo

This demo demonstrates how consistent hashing can be used for database sharding, distributing data across multiple database instances.

## What it does

- Sets up a hash ring with 4 database shards
- Routes user data to specific shards based on user ID
- Shows the impact of adding a new shard (minimal data migration)
- Demonstrates that only ~20% of data needs to move when scaling from 4 to 5 shards

## Run the demo

```bash
task demo:sharding
```

## Key concepts

- **Data distribution**: Users are evenly distributed across database shards
- **Horizontal scaling**: Easy to add new shards with minimal data migration
- **Predictable routing**: User data location is deterministic based on user ID
- **Efficient resharding**: Only affected data moves when adding/removing shards