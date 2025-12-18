# TravelDao Timeout Usage

## Overview

The `TravelDao` now has a built-in `Timeout` property that allows you to set query timeouts without changing your calling code. This provides a clean, encapsulated way to handle slow queries.

## Implementation Details

### What Was Added

1. **`Timeout` property** in `TravelDao` struct
   - Type: `time.Duration`
   - Default: `0` (no timeout)
   - When set to > 0, queries use `context.WithTimeout`

2. **Automatic timeout handling** in `FindPathSimple3`
   - Checks if `Timeout > 0`
   - If yes: creates context with timeout and uses `QueryContext`
   - If no: uses regular `Query` (backward compatible)

### File Changes

- `application/internal/dao/travel_dao.go`
  - Added `Timeout` field to struct
  - Modified `FindPathSimple3` to support timeout
  - Added `context` and `database/sql` imports

## Usage Examples

### Example 1: No Timeout (Default Behavior)

```go
// Default - no timeout, backward compatible
db, _ := di.NewDatabase("test")
travelDao := dao.NewTravelDao(db)

// Timeout is 0 by default, so query runs without timeout
sequences, err := travelDao.FindPathSimple3(filter)
```

### Example 2: With Timeout

```go
db, _ := di.NewDatabase("test")
travelDao := dao.NewTravelDao(db)

// Set 30 second timeout
travelDao.Timeout = 30 * time.Second

// This query will be cancelled if it takes more than 30 seconds
sequences, err := travelDao.FindPathSimple3(filter)
if err != nil {
    if err.Error() == "context deadline exceeded" {
        log.Println("Query timed out!")
    }
}
```

### Example 3: Different Timeouts for Different Operations

```go
db, _ := di.NewDatabase("test")
travelDao := dao.NewTravelDao(db)

// Quick search - 10 second timeout
travelDao.Timeout = 10 * time.Second
quickResults, err := travelDao.FindPathSimple3(quickFilter)

// Comprehensive search - 60 second timeout
travelDao.Timeout = 60 * time.Second
fullResults, err := travelDao.FindPathSimple3(comprehensiveFilter)

// Disable timeout for specific query
travelDao.Timeout = 0
noTimeoutResults, err := travelDao.FindPathSimple3(filter)
```

### Example 4: In Web Handler

```go
func SearchTravels(c *gin.Context) {
    db, _ := di.NewDatabase("prod")
    travelDao := dao.NewTravelDao(db)

    // Set reasonable timeout for web requests
    travelDao.Timeout = 30 * time.Second

    sequences, err := travelDao.FindPathSimple3(filter)
    if err != nil {
        if err.Error() == "context deadline exceeded" {
            c.JSON(http.StatusRequestTimeout, gin.H{
                "error": "Search took too long, please try a more specific search",
            })
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"sequences": sequences})
}
```

## Error Handling

When a query times out, you'll get:
- **Error**: `context deadline exceeded`
- **Query**: Cancelled by the MySQL driver via `KILL QUERY`
- **Connection**: Closed and not returned to pool

```go
sequences, err := travelDao.FindPathSimple3(filter)
if err != nil {
    switch err.Error() {
    case "context deadline exceeded":
        log.Println("Query timed out")
    case "context canceled":
        log.Println("Query was cancelled")
    default:
        log.Printf("Query error: %v", err)
    }
}
```

## Testing

Run the timeout test:

```bash
go test -tags=draft -v ./application/drafttests -run TestSlowQueryTimeout
```

Expected output:
```
Setting TravelDao timeout to: 30s
Executing query with timeout: 30s
[  1.000s]     .
[  2.000s]     .
...
[ 30.001s]     .
Query finished with error after 30.002s: context deadline exceeded
✅ Query was successfully cancelled by timeout!
SUCCESS: Query timed out as expected after ~30 seconds
```

## Benefits

1. **Encapsulation**: Timeout logic is inside the DAO, not scattered in calling code
2. **Backward Compatible**: Default `Timeout = 0` means existing code works unchanged
3. **Flexible**: Can change timeout per DAO instance or per operation
4. **Clean**: No context passing through multiple layers
5. **Minimal Changes**: Calling code stays simple

## Future Improvements

You can extend this pattern to other methods:

```go
// Apply same pattern to other FindPath methods
func (td *TravelDao) FindPathSimple2(filter *data.TravelFilter) ([]*tables.TransferSequence, error) {
    // ... same timeout logic
    if td.Timeout > 0 {
        ctx, cancel := context.WithTimeout(context.Background(), td.Timeout)
        defer cancel()
        rows, err = connection.QueryContext(ctx, sql)
    } else {
        rows, err = connection.Query(sql)
    }
    // ...
}
```

Or create a helper method:

```go
func (td *TravelDao) queryWithOptionalTimeout(sql string) (*sql.Rows, error) {
    connection, err := td.database.GetConnection()
    if err != nil {
        return nil, err
    }

    if td.Timeout > 0 {
        ctx, cancel := context.WithTimeout(context.Background(), td.Timeout)
        defer cancel()
        return connection.QueryContext(ctx, sql)
    }
    return connection.Query(sql)
}
```

## Notes

- The timeout applies to the entire query execution
- MySQL driver will send `KILL QUERY` when timeout occurs
- Connection is closed after timeout (not returned to pool)
- Works with both MySQL and MariaDB
