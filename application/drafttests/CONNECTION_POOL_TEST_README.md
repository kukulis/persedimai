# Connection Pool Test Documentation

## Overview

This test suite (`connection_pool_test.go`) demonstrates and validates how database connection pool settings work in practice.

## Connection Pool Settings Being Tested

```go
conn.SetMaxOpenConns(25)          // Maximum 25 connections can be open at once
conn.SetMaxIdleConns(5)           // Maximum 5 idle connections kept in pool
conn.SetConnMaxLifetime(5 * time.Minute)   // Connections older than 5 minutes are closed
conn.SetConnMaxIdleTime(1 * time.Minute)   // Idle connections older than 1 minute are closed
```

## Running the Tests

Since these are draft tests, you need to use the `draft` build tag:

```bash
# Run all connection pool tests
go test -tags=draft -v ./application/drafttests -run TestConnectionPool

# Run specific test
go test -tags=draft -v ./application/drafttests -run TestConnectionPoolBehavior

# Run stress test only
go test -tags=draft -v ./application/drafttests -run TestConnectionPoolUnderStress
```

## Test Scenarios

### TestConnectionPoolBehavior

This test demonstrates pool behavior through 6 scenarios:

1. **Initial Pool Stats** - Shows baseline state
2. **Sequential Queries** - 5 queries one after another
3. **Concurrent Queries (Low)** - 10 concurrent queries (under MaxOpenConns)
4. **High Concurrent Load** - 50 concurrent queries (exceeds MaxOpenConns=25)
5. **Idle Connection Cleanup** - Monitors pool for 70 seconds to see idle timeout in action
6. **Pool Recovery** - Shows pool can recover and create new connections after cleanup

**Expected observations:**
- InUse count rises during concurrent queries
- WaitCount increases when concurrent queries exceed MaxOpenConns (25)
- MaxIdleClosed increases when idle connections exceed MaxIdleConns (5)
- MaxIdleTimeClosed increases after connections are idle for > 1 minute

### TestConnectionPoolUnderStress

This test simulates real-world sustained load:
- 100 concurrent workers
- Each worker continuously makes queries
- Runs for 30 seconds
- Monitors stats every 5 seconds

**Expected observations:**
- OpenConnections should stabilize at or near MaxOpenConns (25)
- WaitCount will increase as 100 workers compete for 25 connections
- WaitDuration shows cumulative time queries waited for connections

## Understanding the Stats

### Key Metrics

- **MaxOpenConnections**: Your configured limit (25)
- **OpenConnections**: Currently open connections to the database
- **InUse**: Connections actively executing queries
- **Idle**: Connections open but not currently used
- **WaitCount**: Number of times a query had to wait because all connections were busy
- **WaitDuration**: Total time queries spent waiting for connections
- **MaxIdleClosed**: Connections closed because pool had more than MaxIdleConns (5) idle
- **MaxLifetimeClosed**: Connections closed because they exceeded ConnMaxLifetime (5 min)
- **MaxIdleTimeClosed**: Connections closed because idle for > ConnMaxIdleTime (1 min)

### What to Look For

**Good signs:**
- Low WaitCount under normal load
- Idle connections between 0 and MaxIdleConns
- OpenConnections scales with load up to MaxOpenConns

**Potential issues:**
- High WaitCount with low OpenConnections → increase MaxOpenConns
- WaitDuration very high → queries are waiting too long
- OpenConnections always at max → may need more connections
- Many MaxIdleClosed → MaxIdleConns might be too low

## Tuning Guidelines

Based on test results, you can adjust settings:

### If queries are waiting too much (high WaitCount):
```go
conn.SetMaxOpenConns(50)  // Increase from 25
```

### If too many idle connections are kept:
```go
conn.SetMaxIdleConns(2)   // Decrease from 5
```

### If connections should be recycled faster:
```go
conn.SetConnMaxLifetime(2 * time.Minute)  // Decrease from 5
```

### If idle connections should be closed sooner:
```go
conn.SetConnMaxIdleTime(30 * time.Second)  // Decrease from 1 minute
```

## Notes

- The tests use `SELECT SLEEP()` to simulate long-running queries
- Test 5 waits for 70 seconds to demonstrate idle timeout behavior
- Stress test runs for 30 seconds and generates significant database load
- All tests properly close result sets to return connections to the pool
