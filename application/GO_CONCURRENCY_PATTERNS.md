# Go Concurrency Patterns: Channels vs Shared Variables

## Go's Philosophy

> "Don't communicate by sharing memory, share memory by communicating."

But this is a **guideline**, not an absolute rule!

## Decision Matrix

### ✅ Use Channels When:

**1. Ownership Transfer**
```go
// Good: Channel transfers ownership
jobs := make(chan Job)
go worker(jobs)
jobs <- job  // Worker now owns this job
```

**2. Distributing Work**
```go
// Good: Fan-out pattern
for i := 0; i < numWorkers; i++ {
    go worker(jobsChan, resultsChan)
}
```

**3. Signaling/Coordination**
```go
// Good: Signal completion
done := make(chan struct{})
go func() {
    doWork()
    close(done)
}()
<-done
```

**4. Pipeline/Stream Processing**
```go
// Good: Data flows through stages
stage1 := producer()
stage2 := transform(stage1)
stage3 := consumer(stage2)
```

### ✅ Use Shared Variables (with Mutexes) When:

**1. Shared State/Counters**
```go
// Good: Mutex for shared counter
type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (c *SafeCounter) Inc() {
    c.mu.Lock()
    c.count++
    c.mu.Unlock()
}
```

**2. Caching/Lookups**
```go
// Good: RWMutex for read-heavy cache
type Cache struct {
    mu    sync.RWMutex
    items map[string]Item
}

func (c *Cache) Get(key string) (Item, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    item, ok := c.items[key]
    return item, ok
}
```

**3. Configuration/State That Doesn't Transfer**
```go
// Good: Atomic for simple flags
type Server struct {
    running atomic.Bool
}

func (s *Server) Start() {
    s.running.Store(true)
}

func (s *Server) IsRunning() bool {
    return s.running.Load()
}
```

**4. Performance-Critical Sections**
```go
// Good: Atomic operations for hot paths
var requestCount atomic.Int64

func handleRequest() {
    requestCount.Add(1)  // Much faster than channel
    // ...
}
```

## Real-World Examples

### Example 1: Connection Pool (Channels ✅)

```go
// Good: Channels are perfect for this
type Pool struct {
    connections chan *Connection
}

func (p *Pool) Get() *Connection {
    return <-p.connections  // Get available connection
}

func (p *Pool) Put(conn *Connection) {
    p.connections <- conn  // Return to pool
}
```

### Example 2: Metrics Collector (Mutex ✅)

```go
// Good: Mutex for aggregating metrics
type Metrics struct {
    mu            sync.Mutex
    requestCount  int
    errorCount    int
    totalDuration time.Duration
}

func (m *Metrics) RecordRequest(duration time.Duration, err error) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.requestCount++
    if err != nil {
        m.errorCount++
    }
    m.totalDuration += duration
}
```

**Why not channels here?**
- ❌ Channels would require a dedicated goroutine to aggregate
- ❌ More complex for simple increments
- ❌ Worse performance (channel overhead)

### Example 3: Rate Limiter (Both!)

```go
// Hybrid approach - uses both!
type RateLimiter struct {
    tokens chan struct{}      // Channel for token bucket
    mu     sync.Mutex          // Mutex for refill logic
    rate   int
}
```

## Common Patterns & Tools

### 1. **sync.Once** - One-Time Initialization
```go
var (
    instance *Database
    once     sync.Once
)

func GetDB() *Database {
    once.Do(func() {
        instance = &Database{...}
    })
    return instance
}
```

### 2. **atomic.Value** - Lock-Free Reads
```go
var config atomic.Value  // Store *Config

func UpdateConfig(c *Config) {
    config.Store(c)  // Atomic write
}

func GetConfig() *Config {
    return config.Load().(*Config)  // Atomic read
}
```

### 3. **sync.WaitGroup** - Wait for Multiple Goroutines
```go
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        doWork()
    }()
}
wg.Wait()
```

### 4. **context.Context** - Cancellation
```go
ctx, cancel := context.WithCancel(context.Background())
go worker(ctx)
// Later:
cancel()  // Signal worker to stop
```

## Performance Comparison

```go
// Benchmark: Atomic vs Mutex vs Channel
func BenchmarkAtomic(b *testing.B) {
    var counter atomic.Int64
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            counter.Add(1)  // ~5-10ns
        }
    })
}

func BenchmarkMutex(b *testing.B) {
    var mu sync.Mutex
    var counter int64
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            mu.Lock()
            counter++      // ~20-30ns
            mu.Unlock()
        }
    })
}

func BenchmarkChannel(b *testing.B) {
    ch := make(chan int64, 100)
    go func() {
        for range ch {}
    }()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            ch <- 1        // ~100-200ns
        }
    })
}
```

**Results**: Atomic < Mutex < Channel (for simple operations)

## Anti-Patterns to Avoid

### ❌ Don't: Channel for Simple Counter
```go
// BAD: Overkill for a counter
counter := make(chan int)
go func() {
    count := 0
    for range counter {
        count++
    }
}()

// Just use atomic!
var count atomic.Int64
count.Add(1)
```

### ❌ Don't: Mutex for Ownership Transfer
```go
// BAD: Mutex when you should transfer ownership
type Queue struct {
    mu    sync.Mutex
    items []Item
}

// Use channel instead!
items := make(chan Item)
```

### ❌ Don't: Unsynchronized Shared Variables
```go
// VERY BAD: Race condition!
var count int

go func() {
    count++  // Write
}()

println(count)  // Read - DATA RACE!
```

## Decision Tree

```
Need to share data between goroutines?
│
├─ Is it ownership transfer / pipeline?
│  └─ YES → Use Channel
│
├─ Is it a simple counter / flag?
│  └─ YES → Use atomic
│
├─ Is it read-heavy (10:1 reads:writes)?
│  └─ YES → Use sync.RWMutex
│
├─ Is it complex shared state?
│  └─ YES → Use sync.Mutex
│
└─ Just need to wait for completion?
   └─ YES → Use sync.WaitGroup or channel close
```

## Summary Table

| Use Case | Best Tool | Why |
|----------|-----------|-----|
| Worker pool | Channel | Natural ownership transfer |
| Request counter | atomic.Int64 | Fastest, simple |
| Cache | sync.RWMutex | Many reads, few writes |
| Config | atomic.Value | Lock-free reads |
| Pipeline | Channel | Data flows between stages |
| Shared complex state | sync.Mutex | Multiple related fields |
| One-time init | sync.Once | Guaranteed once |
| Wait for goroutines | sync.WaitGroup | Built for this |
| Cancellation | context.Context | Standard pattern |

## Why errChan Is Needed (Race Condition Example)

### The Problem with Shared `err`

```go
var err error  // Declared in outer scope

go func() {
    sequences, err = travelDao.FindPathSimple3(...)  // Goroutine writes to err
    stopChan <- true
}()

// Later in main goroutine:
if finished {
    // Main goroutine reads err - RACE CONDITION!
    if err != nil { ... }
}
```

### Why This Is a Race Condition

Even though `stopChan` signals that the goroutine finished, **Go's memory model doesn't guarantee** that the write to `err` is visible to the main goroutine without explicit synchronization.

**What Could Go Wrong:**

1. **CPU caching**: The `err` value might still be in the goroutine's CPU cache
2. **Compiler optimizations**: The compiler might reorder operations
3. **Memory visibility**: No guarantee the write has propagated to main goroutine's view of memory

### Testing With Race Detector:

```bash
go test -race -tags=draft ./application/drafttests -run TestSlowQueryTimeout
```

Without `errChan`, you'd likely see:
```
WARNING: DATA RACE
Write at 0x... by goroutine X:
  TestSlowQueryTimeout.func1()
Read at 0x... by main goroutine:
  TestSlowQueryTimeout()
```

### How `errChan` Solves This

Channels in Go provide **synchronization guarantees**:

```go
go func() {
    sequences, err = travelDao.FindPathSimple3(...)
    errChan <- err  // Send provides memory barrier
    stopChan <- true
}()

// Later:
if finished {
    queryErr := <-errChan  // Receive provides memory barrier
}
```

**Go's memory model guarantees:**
- All writes before `errChan <- err` are visible after `<-errChan`
- This is called a **happens-before relationship**

### Could You Use Just `err`?

Technically, **it might work** in this specific case because:

```go
go func() {
    sequences, err = travelDao.FindPathSimple3(...)  // Write 1
    stopChan <- true                                   // Write 2 (synchronized)
}()

// Main goroutine:
<-stopChan     // Read (synchronized) - happens-before relationship
use(err)       // Read - might be safe
```

The `stopChan` synchronization *might* provide transitivity, but:
- ❌ Not explicitly guaranteed by Go spec
- ❌ Race detector will complain
- ❌ Subtle and error-prone
- ❌ Not idiomatic Go

### The Correct Pattern

```go
// Good: Explicit data channel
errChan := make(chan error, 1)
go func() {
    _, err := doWork()
    errChan <- err
}()
result := <-errChan
```

Or with multiple return values:
```go
// Even better: Return all results via channel
type Result struct {
    Sequences []*tables.TransferSequence
    Err       error
}

resultChan := make(chan Result, 1)
go func() {
    sequences, err := doWork()
    resultChan <- Result{sequences, err}
}()
result := <-resultChan
```

### Comparison Table

| Approach | Race-Safe? | Race Detector? | Idiomatic? |
|----------|-----------|----------------|------------|
| Shared `err` variable | ⚠️ Unclear | ❌ Fails | ❌ No |
| `errChan` | ✅ Yes | ✅ Passes | ✅ Yes |

## Real-World Advice

1. **Start with channels** if ownership transfer is involved
2. **Use atomic** for simple counters/flags
3. **Use mutexes** for protecting complex state
4. **Profile** if performance matters
5. **Run with `-race`** to catch bugs
6. **Choose clarity** over cleverness

The Go proverb is good guidance, but **pragmatism wins**. Sometimes a mutex is just clearer and faster than spinning up a goroutine to manage a channel!

## Additional Resources

- [Go Memory Model](https://go.dev/ref/mem)
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Go Blog - Share Memory By Communicating](https://go.dev/blog/codelab-share)
- [Go Blog - Race Detector](https://go.dev/blog/race-detector)
