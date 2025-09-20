# Consistent Hash Ring Workshop

A (very roughly) 3-hour hands-on workshop for building a consistent hash ring in Go.

## Workshop Overview

Learn how to design and implement a consistent hash ring from scratch, solving the problem of efficient data
distribution in distributed systems.

## Repository Structure

```
hashlab/
├── cmd/
│   └── demo/
│       └── main.go              # Main demo application
├── hashring/
│   ├── hashing.go               # Core hash ring implementation
│   ├── hashing_test.go          # Unit tests
│   ├── hashing_bench_test.go    # Performance benchmarks
│   └── metrics.go               # Performance metrics and analysis
└── examples/
    ├── cache/                   # Cache distribution demo
    ├── compare/                 # Comparison of hashing strategies
    ├── loadbalancer/            # Load balancing demo
    ├── performance/             # Performance analysis demo
    └── sharding/                # Database sharding demo
```

## Git Branches

This repository contains multiple branches representing different stages of implementation. Each one builds upon the
preceding branch. For example, _02-consistent-hashing_ represents the completion of _01-naive-hashing_ and the starting
point for the second hands-on session.

- **`01-naive-hashing`** - Empty structs/tests and outlined API (starting point)
- **`02-consistent-hashing`** - Basic modulo hashing (naive approach)
- **`03-virtual-nodes`** - Consistent hashing without virtual nodes
- **`04-perf-analysis`** - Full implementation with virtual nodes
- **`main`** - Complete implementation with all features

## Getting Started

1. Clone this repository
2. Check out the `01-naive-hashing` branch to start from scratch:

   ```bash
   git checkout 01-naive-hashing
   ```

3. Follow along with the workshop slides

## Running the Code

You'll need Go (1.24 or newer) installed. See [installation instructions](https://go.dev/doc/install) if necessary.

You'll also need a recent version of `task`. You can install it like this (where `<dir>` should be replaced by a
directory on your PATH - e.g. _/usr/local/bin_):

```bash
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b <dir> v3.44.1
```

### Working with the code base

```bash
# Run the demo
task run

# Run tests
task test

# Run benchmarks
task test:bench

# Run example app
task demo:<dir> [-- <args>]
```

## Learning Objectives

By the end of this workshop, you will:

1. Understand the limitations of simple hash-based sharding
2. Implement a consistent hash ring from scratch
3. Add virtual nodes for better key distribution
4. Write comprehensive tests for distributed systems
5. Optimize performance for production use

## Prerequisites

- Basic Go knowledge
- Understanding of hash functions
- Familiarity with distributed systems concepts (helpful but not required)

## Key Concepts Covered

- **Problem:** Single database overload
- **Solution:** Sharding with consistent hashing
- **Challenge:** Minimizing data movement during scaling
- **Technique:** Virtual nodes for even distribution

## Real-World Applications

- Distributed caching (Memcached, Redis)
- Load balancing (HAProxy, Nginx)
- Distributed databases (Cassandra, DynamoDB)
- CDN edge selection
- Microservice discovery

## Estimated Workshop Timeline

1. **Introduction & Problem Statement** (15 min)
2. **Sharding Basics** (15 min)
3. **Consistent Hashing Theory** (20 min)
4. **Hands-on: Building the Hash Ring** (45 min)
5. **Virtual Nodes Implementation** (30 min)
6. **Testing & Optimization** (30 min)
7. **Q&A & Advanced Topics** (25 min)

## Resources

- [Consistent Hashing and Random Trees (Karger et al.)](https://www.akamai.com/us/en/multimedia/documents/technical-publication/consistent-hashing-and-random-trees-distributed-caching-protocols-for-relieving-hot-spots-on-the-world-wide-web-technical-publication.pdf)
- [Cassandra Architecture](https://cassandra.apache.org/doc/latest/architecture/dynamo.html)
- [Go hash/crc32 Documentation](https://pkg.go.dev/hash/crc32)

## License

MIT License - Feel free to use this for learning and teaching purposes.
