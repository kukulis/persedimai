# Database Timeout Solution Proposal

## Problem Statement

Current implementation:
- No timeout control on SQL queries
- Long-running queries block user requests indefinitely
- Database can become overloaded with slow queries
- No user feedback when queries take too long

## Requirements

1. **User-facing timeout**: 15 seconds
2. **User feedback**: Clear message when query exceeds 15s
3. **Background continuation**: Query continues for up to 5 minutes total
4. **Final cancellation**: Query canceled after 5 minutes
5. **Database protection**: Prevent overload from long-running queries

## Solution Architecture

### Option 1: Two-Tier Context Timeout (Recommended)

**Architecture:**
```
User Request (15s timeout)
    ├─> Return "query taking too long" message after 15s
    └─> Background goroutine (5min timeout)
            ├─> Query continues execution
            ├─> Log results if completed
            └─> Cancel after 5 minutes
```

**Key Components:**

1. **Request Context** - 15 second timeout for user response
2. **Background Context** - 5 minute timeout for actual query
3. **Query Cancellation** - MySQL connection kill on timeout
4. **Result Channel** - For async result handling

**Pros:**
- Immediate user feedback
- Efficient resource utilization
- Proper query cancellation
- Best practice for Go

**Cons:**
- More complex implementation
- Requires goroutine management
- Background results may not be useful if user has left

### Option 2: Simple Timeout Only

**Architecture:**
```
User Request (15s timeout)
    └─> Cancel query and return error after 15s
```

**Pros:**
- Simple to implement
- Immediate resource cleanup
- No background processing complexity

**Cons:**
- No opportunity for slow queries to complete
- May waste work done by database
- Less graceful degradation

---

## Recommended Implementation (Option 1)

### 1. Database Layer Enhancement

Add context-aware query execution to `Database` struct:

```go
// database/database.go

import (
    "context"
    "database/sql"
    "time"
)

type QueryResult struct {
    Rows  *sql.Rows
    Error error
}

// ExecuteQueryWithTimeout executes a query with two-tier timeout
func (db *Database) ExecuteQueryWithTimeout(
    query string,
    args []interface{},
    userTimeout time.Duration,      // 15 seconds
    backgroundTimeout time.Duration, // 5 minutes
) (*sql.Rows, error, bool) {

    conn, err := db.GetConnection()
    if err != nil {
        return nil, err, false
    }

    // Create background context with 5-minute timeout
    bgCtx, bgCancel := context.WithTimeout(context.Background(), backgroundTimeout)
    defer bgCancel()

    // Create user-facing context with 15-second timeout
    userCtx, userCancel := context.WithTimeout(bgCtx, userTimeout)
    defer userCancel()

    // Channel for query results
    resultChan := make(chan QueryResult, 1)

    // Execute query in goroutine
    go func() {
        rows, err := conn.QueryContext(bgCtx, query, args...)
        resultChan <- QueryResult{Rows: rows, Error: err}
    }()

    // Wait for either result or user timeout
    select {
    case result := <-resultChan:
        // Query completed within 15s
        return result.Rows, result.Error, false

    case <-userCtx.Done():
        // User timeout exceeded
        log.Printf("Query exceeded user timeout, continuing in background: %s", query)

        // Continue in background
        go func() {
            select {
            case result := <-resultChan:
                if result.Error != nil {
                    log.Printf("Background query failed: %v", result.Error)
                } else {
                    log.Printf("Background query completed successfully")
                    // Optionally: cache results, update metrics, etc.
                    result.Rows.Close()
                }
            case <-bgCtx.Done():
                log.Printf("Background query exceeded 5min timeout, canceling")
            }
        }()

        return nil, ErrQueryTimeout, true
    }
}
```

### 2. DAO Layer Integration

Modify DAOs to use context-aware methods:

```go
// dao/travel_dao.go

func (td *TravelDao) FindPathWithTimeout(
    filter *data.TravelFilter,
) ([]*tables.TransferSequence, error) {

    query := td.buildQuerySQL(filter)

    rows, err, timedOut := td.database.ExecuteQueryWithTimeout(
        query,
        []interface{}{}, // args
        15*time.Second,  // user timeout
        5*time.Minute,   // background timeout
    )

    if timedOut {
        return nil, &TimeoutError{
            Message: "Query is taking longer than expected. Please try with narrower search criteria.",
            UserTimeout: true,
        }
    }

    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Process results...
    return results, nil
}
```

### 3. Error Types

Define custom error types for timeout handling:

```go
// database/errors.go

type TimeoutError struct {
    Message      string
    UserTimeout  bool // true if user-facing timeout, false if background
}

func (e *TimeoutError) Error() string {
    return e.Message
}

var ErrQueryTimeout = &TimeoutError{
    Message:     "Query execution timeout",
    UserTimeout: true,
}
```

### 4. Web Controller Integration

Update controllers to handle timeout errors:

```go
// web/travel_search_controller.go

func (controller *TravelSearchController) SearchResult(c *gin.Context) {
    // ... setup ...

    paths, err := strategy.FindPath(filter)

    if err != nil {
        var timeoutErr *database.TimeoutError
        if errors.As(err, &timeoutErr) && timeoutErr.UserTimeout {
            c.HTML(http.StatusOK, "travel-search-result.html", gin.H{
                "data": SearchResultData{
                    Error: "Your search is taking longer than expected. " +
                           "The query is still running in the background. " +
                           "Please try narrowing your search criteria or try again later.",
                },
            })
            return
        }

        // Handle other errors...
    }

    // Display results...
}
```

### 5. Configuration

Make timeouts configurable:

```go
// database/db_config.go

type DBConfig struct {
    // ... existing fields ...

    // Timeout configurations
    UserQueryTimeout       time.Duration // Default: 15s
    BackgroundQueryTimeout time.Duration // Default: 5min
}

func DefaultTimeouts() *DBConfig {
    return &DBConfig{
        UserQueryTimeout:       15 * time.Second,
        BackgroundQueryTimeout: 5 * time.Minute,
    }
}
```

---

## Advanced Features

### 1. Query Cancellation via Connection Kill

For MySQL, implement actual query cancellation:

```go
func (db *Database) killQuery(connectionID uint32) error {
    killConn, err := db.GetConnection()
    if err != nil {
        return err
    }

    _, err = killConn.Exec(fmt.Sprintf("KILL QUERY %d", connectionID))
    return err
}
```

### 2. Query Metrics and Monitoring

Track slow queries:

```go
type QueryMetrics struct {
    Query        string
    Duration     time.Duration
    TimedOut     bool
    Timestamp    time.Time
}

func (db *Database) logSlowQuery(metrics QueryMetrics) {
    if metrics.Duration > 1*time.Second || metrics.TimedOut {
        log.Printf("SLOW QUERY [%v]: %s", metrics.Duration, metrics.Query)
        // Send to monitoring system (Prometheus, DataDog, etc.)
    }
}
```

### 3. Connection Pool Tuning

Configure connection pool to handle timeouts better:

```go
func (db *Database) configureConnectionPool() {
    db.connection.SetMaxOpenConns(25)
    db.connection.SetMaxIdleConns(5)
    db.connection.SetConnMaxLifetime(5 * time.Minute)
    db.connection.SetConnMaxIdleTime(1 * time.Minute)
}
```

---

## Migration Strategy

### Phase 1: Infrastructure (Week 1)
1. Add context support to Database layer
2. Implement timeout error types
3. Add configuration options

### Phase 2: DAO Updates (Week 2)
1. Update TravelDao methods
2. Update PointDao methods
3. Add tests for timeout behavior

### Phase 3: Controller Integration (Week 3)
1. Update web controllers
2. Improve error messages
3. Add user feedback

### Phase 4: Monitoring (Week 4)
1. Add query metrics
2. Implement slow query logging
3. Dashboard for timeout monitoring

---

## Testing Strategy

### Unit Tests
```go
func TestExecuteQueryWithTimeout_UserTimeout(t *testing.T) {
    // Test that query returns timeout error after 15s
}

func TestExecuteQueryWithTimeout_BackgroundSuccess(t *testing.T) {
    // Test that background query completes after user timeout
}

func TestExecuteQueryWithTimeout_BackgroundTimeout(t *testing.T) {
    // Test that background query is canceled after 5min
}
```

### Integration Tests
```go
func TestSlowQuery_RealDatabase(t *testing.T) {
    // Test with actual slow query (e.g., SLEEP(20))
    query := "SELECT SLEEP(20)"
    // Verify timeout behavior
}
```

---

## Security Considerations

1. **Resource Limits**: Prevent DoS from many slow queries
2. **Query Validation**: Validate queries before execution
3. **Connection Limits**: Enforce max connections per user
4. **Audit Logging**: Log all timeout events

---

## Performance Impact

**Pros:**
- Faster user response (always under 15s)
- Better resource utilization
- Improved user experience

**Cons:**
- Small overhead from goroutines (~2KB per query)
- Additional complexity in error handling
- Need to manage background goroutines

**Estimated Impact:**
- Memory: +2KB per concurrent query
- CPU: Negligible (<1%)
- Latency: No change for fast queries, much better for slow queries

---

## Alternative: Simple Circuit Breaker

If two-tier timeout is too complex, implement a simple circuit breaker:

```go
type CircuitBreaker struct {
    maxConcurrent int
    current       atomic.Int32
}

func (cb *CircuitBreaker) Allow() bool {
    current := cb.current.Load()
    if current >= int32(cb.maxConcurrent) {
        return false
    }
    cb.current.Add(1)
    return true
}

func (cb *CircuitBreaker) Release() {
    cb.current.Add(-1)
}
```

---

## Recommendation

**Implement Option 1 (Two-Tier Context Timeout)** with the following priorities:

1. **Immediate** (Week 1):
   - Add basic context timeout to critical queries
   - Simple error messages for users

2. **Short-term** (Weeks 2-3):
   - Full two-tier timeout implementation
   - Background continuation logic
   - Improved error handling

3. **Long-term** (Month 2+):
   - Query metrics and monitoring
   - Automatic query optimization suggestions
   - Adaptive timeout based on query patterns

This approach balances user experience, database protection, and implementation complexity.
