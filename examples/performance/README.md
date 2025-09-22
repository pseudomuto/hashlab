# Performance Analysis Demo

This demo provides performance benchmarking and analysis for consistent hashing operations at scale.

## What it does

- Tests hash ring performance with configurable parameters
- Analyzes key distribution across multiple servers
- Measures performance metrics with varying load sizes
- Demonstrates scalability of the consistent hashing implementation

## Run the demo

```bash
task demo:performance
```

## Configuration options

- `--keys`: Number of keys to test (default: 10,000)
- `--servers`: Number of servers to distribute across (default: 3)
- `--vnodes`: Number of virtual nodes per server (default: 150)

## Key concepts

- **Performance metrics**: Analyze throughput and distribution quality
- **Scalability testing**: Test with varying numbers of keys and servers
- **Distribution analysis**: Examine how evenly keys are distributed across servers
