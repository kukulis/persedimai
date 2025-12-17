# Database Timeout Quick Reference

## Visual Flow

```
┌─────────────────────────────────────────────────────────────┐
│                     User Request                             │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│           ExecuteQueryWithTimeout()                          │
│  ┌────────────────────────────────────────────────────┐    │
│  │ Context: User (15s) + Background (5min)            │    │
│  └────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
                              │
                    ┌─────────┴─────────┐
                    │                   │
                    ▼                   ▼
           ┌────────────────┐  ┌────────────────┐
           │ Query Results  │  │ User Timeout   │
           │  (< 15s)       │  │  (≥ 15s)       │
           └────────────────┘  └────────────────┘
                    │                   │
                    │                   ├──► Return timeout error to user
                    │                   │
                    │                   ├──► Continue in background goroutine
                    │                   │           │
                    │                   │           ▼
                    │                   │    ┌──────────────┐
                    │                   │    │ Query Result │
                    │                   │    │  (< 5min)    │
                    │                   │    └──────────────┘
                    │                   │           │
                    │                   │           ├──► Log success
                    │                   │           └──► Close resources
                    │                   │
                    │                   └──► Or: Background timeout (≥ 5min)
                    │                              │
                    │                              ├──► Cancel query
                    │                              └──► Log timeout
                    │
                    └──► Return results to user
```

## Key Components

### 1. TimeoutConfig
```go
type TimeoutConfig struct {
    UserTimeout       time.Duration // 15 seconds
    BackgroundTimeout time.Duration // 5 minutes
}
```

### 2. Error Types
```go
type TimeoutError struct {
    Message     string
    UserTimeout bool
}
```

### 3. Main Method
```go
func (db *Database) ExecuteQueryWithTimeout(
    query string,
    args []interface{},
    config *TimeoutConfig,
) (*sql.Rows, error, bool)
```

**Returns:**
- `*sql.Rows` - Query results (nil if timeout)
- `error` - Error if query failed
- `bool` - true if user timeout (continuing in background)

## Implementation Checklist

### Phase 1: Core Infrastructure
- [ ] Add `timeout.go` to database package
- [ ] Define `TimeoutConfig` struct
- [ ] Define `TimeoutError` type
- [ ] Implement `ExecuteQueryWithTimeout()`
- [ ] Implement `ExecuteWithTimeout()` for non-queries
- [ ] Add tests for timeout behavior

### Phase 2: DAO Integration
- [ ] Update `TravelDao.FindPath*()` methods
- [ ] Update `PointDao` methods if needed
- [ ] Add timeout configuration to DAO constructors
- [ ] Update integration tests

### Phase 3: Controller Updates
- [ ] Update web controllers to handle `TimeoutError`
- [ ] Add user-friendly error messages
- [ ] Add HTTP status code handling (408 Request Timeout)
- [ ] Update templates with helpful timeout messages

### Phase 4: Monitoring
- [ ] Add query duration logging
- [ ] Track timeout metrics
- [ ] Set up alerts for frequent timeouts
- [ ] Create dashboard for slow queries

## Usage Examples

### In DAO
```go
func (td *TravelDao) FindPathSimple1(filter *data.TravelFilter) ([]*tables.TransferSequence, error) {
    query := "SELECT ... WHERE ..."

    rows, err, timedOut := td.database.ExecuteQueryWithTimeout(
        query,
        []interface{}{filter.Source, filter.Destination},
        nil, // Use defaults: 15s user, 5min background
    )

    if timedOut {
        return nil, &database.TimeoutError{
            Message: "Search taking too long. Try narrower criteria.",
            UserTimeout: true,
        }
    }

    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Process rows...
}
```

### In Controller
```go
func (controller *TravelSearchController) SearchResult(c *gin.Context) {
    paths, err := strategy.FindPath(filter)

    if err != nil {
        var timeoutErr *database.TimeoutError
        if errors.As(err, &timeoutErr) && timeoutErr.UserTimeout {
            c.HTML(http.StatusRequestTimeout, "travel-search-result.html", gin.H{
                "data": SearchResultData{
                    Error: "Search is taking longer than 15 seconds. " +
                           "Try reducing the date range or number of transfers.",
                },
            })
            return
        }
        // Other errors...
    }

    // Display results...
}
```

### Custom Timeouts
```go
// For critical queries that need longer timeout
config := &database.TimeoutConfig{
    UserTimeout:       30 * time.Second, // 30s for user
    BackgroundTimeout: 10 * time.Minute, // 10min background
}

rows, err, timedOut := db.ExecuteQueryWithTimeout(query, args, config)
```

## Testing Strategy

### Unit Test Example
```go
func TestExecuteQueryWithTimeout_UserTimeout(t *testing.T) {
    db := setupTestDB()

    // Query that takes 20 seconds
    query := "SELECT SLEEP(20)"

    config := &TimeoutConfig{
        UserTimeout:       1 * time.Second,
        BackgroundTimeout: 30 * time.Second,
    }

    rows, err, timedOut := db.ExecuteQueryWithTimeout(query, nil, config)

    assert.Nil(t, rows)
    assert.Error(t, err)
    assert.True(t, timedOut)
    assert.IsType(t, &TimeoutError{}, err)
}
```

### Integration Test
```go
func TestRealSlowQuery(t *testing.T) {
    db, _ := di.NewDatabase("test")

    // Simulate slow query
    query := `
        SELECT t1.*, t2.*, t3.*
        FROM travels t1
        CROSS JOIN travels t2
        CROSS JOIN travels t3
        WHERE SLEEP(0.1) OR 1=1
    `

    rows, err, timedOut := db.ExecuteQueryWithTimeout(query, nil, nil)

    if timedOut {
        t.Log("Query timed out as expected")
    } else if err != nil {
        t.Errorf("Unexpected error: %v", err)
    } else {
        t.Log("Query completed successfully")
        rows.Close()
    }
}
```

## Monitoring Queries

### Log Slow Queries
```go
func (db *Database) logSlowQuery(query string, duration time.Duration, timedOut bool) {
    if duration > 1*time.Second || timedOut {
        log.Printf("[SLOW QUERY] Duration: %v, TimedOut: %v, Query: %s",
            duration, timedOut, truncateQuery(query))
    }
}
```

### Metrics Collection
```go
var (
    queryDurationHistogram = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "db_query_duration_seconds",
            Help:    "Database query duration in seconds",
            Buckets: []float64{0.1, 0.5, 1, 5, 15, 60, 300},
        },
        []string{"query_type"},
    )

    queryTimeoutCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "db_query_timeouts_total",
            Help: "Total number of database query timeouts",
        },
        []string{"timeout_type"}, // "user" or "background"
    )
)
```

## Best Practices

### DO ✅
- Use default timeouts for most queries
- Log all timeout events
- Provide helpful error messages to users
- Monitor timeout frequency
- Close resources in background goroutines
- Use appropriate timeouts based on query type

### DON'T ❌
- Set timeouts too short (< 5s for complex queries)
- Ignore background query results completely
- Forget to close result sets
- Use blocking operations in background goroutines
- Set background timeout shorter than user timeout
- Allow unlimited concurrent background queries

## Configuration Recommendations

### By Query Type

| Query Type | User Timeout | Background Timeout | Rationale |
|------------|--------------|-------------------|-----------|
| Simple SELECT | 5s | 30s | Should be fast |
| Complex JOIN | 15s | 5min | May need time |
| Aggregations | 10s | 2min | Usually optimized |
| INSERT/UPDATE | 10s | 1min | Should be quick |
| BULK Operations | 30s | 10min | Expected to be slow |

### By Environment

| Environment | User Timeout | Background Timeout | Notes |
|-------------|--------------|-------------------|-------|
| Development | 30s | 10min | Allow for debugging |
| Testing | 15s | 5min | Standard timeouts |
| Production | 10s | 2min | Keep responsive |

## Troubleshooting

### High Timeout Rate
1. Check query optimization (EXPLAIN)
2. Review database indexes
3. Analyze query patterns
4. Consider caching frequently-accessed data
5. Add query result pagination

### Background Queries Not Completing
1. Check database connection pool settings
2. Review max_execution_time on MySQL
3. Monitor database CPU/memory
4. Check for lock contention

### Memory Leaks from Background Goroutines
1. Ensure all result sets are closed
2. Use context cancellation properly
3. Monitor goroutine count
4. Implement goroutine pool if needed

## Migration Path

1. **Week 1**: Implement core timeout logic
2. **Week 2**: Update critical DAOs (TravelDao)
3. **Week 3**: Update web controllers and error handling
4. **Week 4**: Add monitoring and fine-tune timeouts

## Additional Resources

- Go Context Documentation: https://pkg.go.dev/context
- MySQL Query Timeout: https://dev.mysql.com/doc/refman/8.0/en/server-system-variables.html#sysvar_max_execution_time
- Go Database/SQL Best Practices: https://go.dev/doc/database/
