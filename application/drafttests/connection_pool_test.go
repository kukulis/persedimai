//go:build draft

package drafttests

import (
	"darbelis.eu/persedimai/di"
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestConnectionPoolBehavior demonstrates and tests how connection pool settings work
func TestConnectionPoolBehavior(t *testing.T) {
	db, err := di.NewDatabase("test")
	if err != nil {
		t.Fatal(err)
	}
	defer db.CloseConnection()

	conn, err := db.GetConnection()
	if err != nil {
		t.Fatal(err)
	}

	// Configure connection pool with the settings to test
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(5 * time.Minute)
	conn.SetConnMaxIdleTime(1 * time.Minute)

	t.Log("=== Connection Pool Settings ===")
	t.Logf("MaxOpenConns: 25")
	t.Logf("MaxIdleConns: 5")
	t.Logf("ConnMaxLifetime: 5 minutes")
	t.Logf("ConnMaxIdleTime: 1 minute")
	t.Log("")

	// Test 1: Baseline stats
	t.Log("=== Test 1: Initial Pool Stats ===")
	logPoolStats(t, conn)
	t.Log("")

	// Test 2: Sequential queries
	t.Log("=== Test 2: Sequential Queries (5 queries) ===")
	for i := 0; i < 5; i++ {
		rows, err := conn.Query("SELECT 1")
		if err != nil {
			t.Errorf("Query %d failed: %v", i+1, err)
		} else {
			rows.Close()
		}
		time.Sleep(100 * time.Millisecond)
	}
	logPoolStats(t, conn)
	t.Log("")

	// Test 3: Concurrent queries (under MaxOpenConns)
	t.Log("=== Test 3: Concurrent Queries (10 concurrent) ===")
	runConcurrentQueries(t, conn, 10, 500*time.Millisecond)
	logPoolStats(t, conn)
	t.Log("")

	// Test 4: High concurrent load (exceeds MaxOpenConns)
	t.Log("=== Test 4: High Concurrent Load (50 concurrent, exceeds MaxOpenConns=25) ===")
	runConcurrentQueries(t, conn, 50, 200*time.Millisecond)
	logPoolStats(t, conn)
	t.Log("")

	// Test 5: Idle connection cleanup
	t.Log("=== Test 5: Idle Connection Cleanup (wait for idle timeout) ===")
	t.Log("Waiting 70 seconds for connections to become idle and be cleaned up...")
	t.Log("(ConnMaxIdleTime is 1 minute, so idle connections should be closed)")

	// Log stats every 10 seconds
	for i := 0; i < 7; i++ {
		time.Sleep(10 * time.Second)
		t.Logf("After %d seconds:", (i+1)*10)
		logPoolStats(t, conn)
	}
	t.Log("")

	// Test 6: Pool recovery after cleanup
	t.Log("=== Test 6: Pool Recovery (new queries after cleanup) ===")
	runConcurrentQueries(t, conn, 5, 100*time.Millisecond)
	logPoolStats(t, conn)
	t.Log("")
}

// TestConnectionPoolUnderStress tests pool behavior under sustained load
func TestConnectionPoolUnderStress(t *testing.T) {
	db, err := di.NewDatabase("test")
	if err != nil {
		t.Fatal(err)
	}
	defer db.CloseConnection()

	conn, err := db.GetConnection()
	if err != nil {
		t.Fatal(err)
	}

	// Configure pool
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(5 * time.Minute)
	conn.SetConnMaxIdleTime(1 * time.Minute)

	t.Log("=== Stress Test: Sustained Load for 30 seconds ===")
	t.Log("Running 100 concurrent workers, each making queries continuously")
	t.Log("")

	stopChan := make(chan bool)
	var wg sync.WaitGroup
	queryCount := 0
	var queryCountMutex sync.Mutex

	// Start 100 workers
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-stopChan:
					return
				default:
					rows, err := conn.Query("SELECT SLEEP(0.01)")
					if err != nil {
						t.Logf("Worker %d error: %v", workerID, err)
					} else {
						rows.Close()
						queryCountMutex.Lock()
						queryCount++
						queryCountMutex.Unlock()
					}
				}
			}
		}(i)
	}

	// Monitor stats every 5 seconds for 30 seconds
	for i := 0; i < 6; i++ {
		time.Sleep(5 * time.Second)
		queryCountMutex.Lock()
		count := queryCount
		queryCountMutex.Unlock()

		t.Logf("=== After %d seconds ===", (i+1)*5)
		t.Logf("Total queries completed: %d", count)
		logPoolStats(t, conn)
		t.Log("")
	}

	// Stop all workers
	close(stopChan)
	wg.Wait()

	queryCountMutex.Lock()
	finalCount := queryCount
	queryCountMutex.Unlock()

	t.Logf("=== Stress Test Complete ===")
	t.Logf("Total queries completed: %d", finalCount)
	t.Logf("Average queries per second: %.2f", float64(finalCount)/30.0)
}

// runConcurrentQueries executes multiple queries concurrently
func runConcurrentQueries(t *testing.T, conn *sql.DB, count int, queryDuration time.Duration) {
	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			query := fmt.Sprintf("SELECT SLEEP(%f)", queryDuration.Seconds())
			rows, err := conn.Query(query)
			if err != nil {
				t.Logf("Concurrent query %d failed: %v", id, err)
			} else {
				rows.Close()
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)
	t.Logf("Completed %d concurrent queries in %v", count, elapsed)
}

// logPoolStats logs the current database connection pool statistics
func logPoolStats(t *testing.T, conn *sql.DB) {
	stats := conn.Stats()

	t.Logf("Pool Stats:")
	t.Logf("  MaxOpenConnections: %d", stats.MaxOpenConnections)
	t.Logf("  OpenConnections: %d (currently open)", stats.OpenConnections)
	t.Logf("  InUse: %d (actively used)", stats.InUse)
	t.Logf("  Idle: %d (idle and ready)", stats.Idle)
	t.Logf("  WaitCount: %d (queries that had to wait for a connection)", stats.WaitCount)
	t.Logf("  WaitDuration: %v (total time queries waited)", stats.WaitDuration)
	t.Logf("  MaxIdleClosed: %d (connections closed due to SetMaxIdleConns)", stats.MaxIdleClosed)
	t.Logf("  MaxLifetimeClosed: %d (connections closed due to SetConnMaxLifetime)", stats.MaxLifetimeClosed)
	t.Logf("  MaxIdleTimeClosed: %d (connections closed due to SetConnMaxIdleTime)", stats.MaxIdleTimeClosed)
}
