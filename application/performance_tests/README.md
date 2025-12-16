# Performance Tests

This package contains performance tests for the travel pathfinding system.

## Overview

These tests are separated from regular integration tests because they:
- Take a long time to run (45+ seconds for setup)
- Use large datasets (~1000 points, ~1 million transfers)
- Are intended for performance analysis, not regular CI/CD

## Running Performance Tests

### Run the performance test:
```bash
go test ./performance_tests -v
```

### Skip performance tests (short mode):
```bash
go test ./performance_tests -short
```

### Run the benchmark:
```bash
go test ./performance_tests -bench=. -benchtime=30s
```

### run the particular benchmark
```bash
go test ./performance_tests -bench=BenchmarkFindPaths -benchtime=30s
```