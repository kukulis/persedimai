# Test Database Performance Optimization

## Current Configuration

The `performance.cnf` file is automatically loaded by MariaDB and contains optimizations for the travel search workload.

**Multiple configurations available**:
- `performance.cnf` - Balanced for queries and moderate imports (default, recommended)
- `bulk-import-4gb.cnf` - Optimized for large 4GB+ SQL imports (see BULK_IMPORT_GUIDE.md)

For large SQL imports (4GB+), see **[BULK_IMPORT_GUIDE.md](./BULK_IMPORT_GUIDE.md)** for detailed instructions.

## Key Performance Settings

### 1. InnoDB Buffer Pool (Most Important)
```
innodb_buffer_pool_size = 4G
```
- **Default value**: 4GB (adjust based on your RAM)
- **Recommended**: 50-70% of available system RAM
- **For 8GB RAM**: Set to 4G
- **For 16GB RAM**: Set to 8G
- **For 32GB RAM**: Set to 16G

### 2. Large Temporary Tables
```
tmp_table_size = 512M
max_heap_table_size = 512M
```
- Critical for complex JOIN queries with clustered travel data
- Prevents disk-based temporary tables for better performance

### 3. Join Buffers
```
join_buffer_size = 16M
```
- Optimized for multi-transfer searches (1-5 transfers)
- Larger buffers = faster joins

## Performance Impact

With the optimized configuration, you'll see significant improvements:

| Setting | Default (Before) | Optimized (After) | Improvement |
|---------|------------------|-------------------|-------------|
| InnoDB Buffer Pool | 128 MB | 4 GB | **32x larger** |
| Temp Table Size | 16 MB | 512 MB | **32x larger** |
| Max Heap Table | 16 MB | 512 MB | **32x larger** |
| Join Buffer | 256 KB | 16 MB | **64x larger** |
| Query Cache | 0 MB | 256 MB | **New feature** |
| Log File Size | 48 MB | 512 MB | **10x larger** |

**Expected Performance Gains**:
- Complex multi-transfer searches: **2-10x faster**
- JOIN operations on large datasets: **5-20x faster**
- Repeated queries: **Near-instant** (from cache)
- Large result sets: **3-5x faster**

## Applying Changes

After modifying `performance.cnf`:

```bash
# Ensure correct file permissions (must be readable by container)
chmod 644 .docker/tst_db/conf.d/performance.cnf

# Restart the container
docker compose restart tst_db

# Or recreate to ensure clean state
docker compose down
docker compose up -d tst_db

# Verify settings are applied (use mariadb command, not mysql)
docker exec tst_db mariadb -uroot -ptest -e "SHOW VARIABLES LIKE 'innodb_buffer_pool_size';"
docker exec tst_db mariadb -uroot -ptest -e "SHOW VARIABLES LIKE 'tmp_table_size';"
docker exec tst_db mariadb -uroot -ptest -e "SHOW VARIABLES LIKE 'join_buffer_size';"

# Or connect from host (if mariadb-client installed)
mariadb -uroot -ptest -P23314 -h 127.0.0.1 -e "SHOW VARIABLES LIKE 'innodb_buffer_pool_size';"
```

### Expected Values After Configuration

| Setting | Value | Bytes |
|---------|-------|-------|
| innodb_buffer_pool_size | 4GB | 4,294,967,296 |
| tmp_table_size | 512MB | 536,870,912 |
| max_heap_table_size | 512MB | 536,870,912 |
| join_buffer_size | 16MB | 16,777,216 |
| max_statement_time | 300s | 300.000000 |
| query_cache_size | 256MB | 268,435456 |
| innodb_log_file_size | 512MB | 536,870,912 |

## Additional Performance Recommendations

### Database-Level Optimizations

#### 1. Add Indexes for Common Queries
```sql
-- For travel searches (if not already indexed)
CREATE INDEX idx_travels_from_departure ON travels(from_id, departure_time);
CREATE INDEX idx_travels_to_arrival ON travels(to_id, arrival_time);
CREATE INDEX idx_travels_departure ON travels(departure_time);

-- For clustered searches
CREATE INDEX idx_clustered_cluster_from ON clustered_arrival_travels32(cluster_id, from_id);
CREATE INDEX idx_clustered_to_arrival ON clustered_arrival_travels32(to_id, arrival_time);
```

#### 2. Table Statistics
```sql
-- Keep statistics updated for better query plans
ANALYZE TABLE travels;
ANALYZE TABLE clustered_arrival_travels32;
ANALYZE TABLE points;
```

#### 3. Query Optimization
```sql
-- Enable extended statistics
SET GLOBAL optimizer_use_condition_selectivity = 4;

-- Check slow queries
SELECT * FROM mysql.slow_log ORDER BY query_time DESC LIMIT 10;
```

### Application-Level Optimizations

#### 1. Connection Pooling
Ensure your Go application uses connection pooling efficiently:
```go
db.SetMaxOpenConns(100)
db.SetMaxIdleConns(25)
db.SetConnMaxLifetime(5 * time.Minute)
```

#### 2. Prepared Statements
Use prepared statements for repeated queries to reduce parsing overhead.

#### 3. Batch Operations
When seeding data, use batch inserts instead of individual inserts.

### System-Level Optimizations

#### 1. Use SSD Storage
- Mount the Docker volume on SSD storage
- Significantly faster I/O for InnoDB

#### 2. Increase Docker Resources
In Docker Desktop settings:
- Memory: Allocate enough for innodb_buffer_pool_size + OS overhead
- CPU: 4+ cores for parallel query execution

#### 3. Disable Unnecessary Services
```bash
# In performance.cnf (already configured)
performance_schema = OFF
skip-log-bin
```

## Monitoring Performance

### Check Buffer Pool Efficiency
```bash
# Check buffer pool usage and hit ratio
docker exec tst_db mariadb -uroot -ptest -e "SHOW STATUS LIKE 'Innodb_buffer_pool%';"

# Calculate hit ratio (should be > 99%)
docker exec tst_db mariadb -uroot -ptest -e "
SELECT
  ROUND(100 - (Innodb_buffer_pool_reads / Innodb_buffer_pool_read_requests * 100), 2) AS hit_ratio_percent
FROM
  (SELECT VARIABLE_VALUE AS Innodb_buffer_pool_reads FROM information_schema.GLOBAL_STATUS WHERE VARIABLE_NAME = 'Innodb_buffer_pool_reads') AS reads,
  (SELECT VARIABLE_VALUE AS Innodb_buffer_pool_read_requests FROM information_schema.GLOBAL_STATUS WHERE VARIABLE_NAME = 'Innodb_buffer_pool_read_requests') AS requests;
"

# Look for:
# - Innodb_buffer_pool_read_requests (high is good)
# - Innodb_buffer_pool_reads (low is good)
# - Hit ratio > 99% means buffer pool is well-sized
```

### Monitor Query Performance
```bash
# Show currently running queries
docker exec tst_db mariadb -uroot -ptest -e "SHOW FULL PROCESSLIST;"

# Check slow query log (if enabled in performance.cnf)
docker exec tst_db cat /var/log/mysql/slow.log

# Analyze a specific query's execution plan
docker exec tst_db mariadb -uroot -ptest test -e "EXPLAIN SELECT ...your query...;"

# Show query cache statistics
docker exec tst_db mariadb -uroot -ptest -e "SHOW STATUS LIKE 'Qcache%';"
```

### Check Table Sizes
```bash
docker exec tst_db mariadb -uroot -ptest test -e "
SELECT
    table_name,
    ROUND((data_length + index_length) / 1024 / 1024, 2) AS size_mb,
    table_rows,
    ROUND(index_length / 1024 / 1024, 2) AS index_size_mb
FROM information_schema.tables
WHERE table_schema = 'test'
ORDER BY (data_length + index_length) DESC;
"
```

### Monitor Overall Performance
```bash
# Container resource usage
docker stats tst_db --no-stream

# Database uptime and key metrics
docker exec tst_db mariadb -uroot -ptest -e "
SHOW GLOBAL STATUS WHERE
  Variable_name IN ('Uptime', 'Threads_connected', 'Questions', 'Slow_queries', 'Queries');
"

# Check for table locks
docker exec tst_db mariadb -uroot -ptest -e "SHOW OPEN TABLES WHERE In_use > 0;"
```

## Troubleshooting

### Configuration Not Loading (Settings Show Default Values)

**Symptom**: After restarting container, `SHOW VARIABLES` still shows default values (e.g., innodb_buffer_pool_size = 134217728 instead of 4294967296)

**Common Causes**:

1. **File Permissions Issue** (Most Common)
   ```bash
   # Check current permissions
   ls -la .docker/tst_db/conf.d/performance.cnf

   # If permissions are -rw------- (600), MariaDB can't read it
   # Fix: Make file readable by all
   chmod 644 .docker/tst_db/conf.d/performance.cnf

   # Verify inside container
   docker exec tst_db ls -la /etc/mysql/conf.d/

   # Should show: -rw-r--r-- (644)
   ```

2. **File Not Mounted Correctly**
   ```bash
   # Check if file exists in container
   docker exec tst_db cat /etc/mysql/conf.d/performance.cnf

   # Verify docker-compose.yml has volume mount
   # Should have: - ./.docker/tst_db/conf.d:/etc/mysql/conf.d
   ```

3. **Syntax Error in Config File**
   ```bash
   # Check MariaDB logs for errors
   docker logs tst_db 2>&1 | grep -i error

   # Look for lines like: "unknown variable" or "incorrect"
   ```

**Solution Steps**:
```bash
# 1. Fix permissions
chmod 644 .docker/tst_db/conf.d/performance.cnf

# 2. Restart container
docker compose restart tst_db

# 3. Wait a few seconds for startup
sleep 3

# 4. Verify settings loaded
docker exec tst_db mariadb -uroot -ptest -e "SHOW VARIABLES LIKE 'innodb_buffer_pool_size';"
```

### Container Won't Start
```bash
# Check logs for errors
docker logs tst_db

# Common issue: innodb_buffer_pool_size too large for available memory
# Solution: Reduce it in performance.cnf
# Example: Change from 4G to 2G
```

### Out of Memory Errors
```bash
# Check container memory usage
docker stats tst_db

# Check Docker memory limit
docker inspect tst_db | grep -i memory

# Solutions:
# 1. Reduce buffer pool size in performance.cnf
# 2. Increase Docker Desktop memory allocation
# 3. Check total system RAM availability
```

### Settings Applied But Queries Still Slow
```bash
# 1. Check if buffer pool is being utilized
docker exec tst_db mariadb -uroot -ptest -e "SHOW STATUS LIKE 'Innodb_buffer_pool%';"

# 2. Check for slow queries
docker exec tst_db mariadb -uroot -ptest -e "SHOW FULL PROCESSLIST;"

# 3. Verify slow query log (if enabled)
docker exec tst_db cat /var/log/mysql/slow.log

# 4. Check if queries are using indexes
docker exec tst_db mariadb -uroot -ptest -e "EXPLAIN SELECT ...your query...;"
```

### No 'mysql' Executable in Container
Newer MariaDB images use `mariadb` command instead of `mysql`:

```bash
# Use mariadb command instead
docker exec tst_db mariadb -uroot -ptest -e "SHOW VARIABLES;"

# From host (if mariadb-client installed)
mariadb -uroot -ptest -P23314 -h 127.0.0.1
```

## Further Reading

- [MariaDB Performance Tuning](https://mariadb.com/kb/en/optimization-and-tuning/)
- [InnoDB Buffer Pool Configuration](https://mariadb.com/kb/en/innodb-buffer-pool/)
- [Query Optimization](https://mariadb.com/kb/en/query-optimization/)
