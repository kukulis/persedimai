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

## What Gets Tested

### TestPerformanceFindPath
- **Setup**: Calls `FillTestDatabase` to generate ~1024 points and ~1 million transfers
- **Test**: Runs 100 FindPath queries with random source/destination points
- **Measurement**: Measures only the `FindPath` duration (excludes database setup)
- **Metrics**: Reports min, max, avg, median, P95, P99 query times

### BenchmarkFindPath
- Standard Go benchmark using `testing.B`
- Pre-generates 100 test cases
- Runs as many iterations as possible in the benchtime
- Reports ops/sec and ns/op

## Expected Performance

Based on dataset size (~1M transfers, ~1K points):
- **Target**: < 100ms average query time
- **Acceptable**: < 200ms P95
- **Warning threshold**: Logged if exceeded

## Test Data

The tests use the same data generation as integration tests:
- Points: 1024 (32×32 grid)
- Coordinates: 0, 6000, 12000, ..., 186000 (multiples of 6000)
- Transfers: ~1 million over 2-month period
- Travel time: ~6 hours average
- Rest time: 24 hours

## Notes

- Performance tests require a running MySQL database
- Database configuration from `.env` file (test environment)
- Tests can be run in parallel with other test packages
- Cleanup happens automatically (tables are recreated)
