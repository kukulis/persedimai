# Large Bulk Import Optimization Guide

## For 4GB+ SQL File Imports

This guide helps you optimize MariaDB for importing very large SQL files (4GB+).

## Configuration Strategy

We provide two configurations:

1. **performance.cnf** - Balanced for queries and moderate imports (default)
   - Good for: Daily operations, 167MB imports (20s)
   - Buffer pool: 4GB

2. **bulk-import-4gb.cnf** - Optimized for large bulk imports
   - Good for: Initial database setup, 4GB+ imports
   - Buffer pool: 8GB (adjustable)
   - More aggressive write settings

## Quick Start: Import 4GB SQL File

### Step 1: Switch to Bulk Import Configuration

```bash
# Backup current config
cp .docker/tst_db/conf.d/performance.cnf .docker/tst_db/conf.d/performance.cnf.backup

# Replace with bulk import config
cp .docker/tst_db/conf.d/bulk-import-4gb.cnf .docker/tst_db/conf.d/performance.cnf

# Fix permissions
chmod 644 .docker/tst_db/conf.d/performance.cnf

# Restart container
docker compose restart tst_db

# Verify settings (wait 3 seconds for startup)
sleep 3
docker exec tst_db mariadb -uroot -ptest -e "SHOW VARIABLES LIKE 'innodb_buffer_pool_size';"
# Should show: 8589934592 (8GB)

docker exec tst_db mariadb -uroot -ptest -e "SHOW VARIABLES LIKE 'innodb_log_file_size';"
# Should show: 2147483648 (2GB)
```

### Step 2: Prepare Database for Import

```bash
# Connect to database and disable safety checks
docker exec -i tst_db mariadb -uroot -ptest test << 'EOF'
-- Disable foreign key checks (huge speedup!)
SET FOREIGN_KEY_CHECKS=0;

-- Disable unique key checks temporarily
SET UNIQUE_CHECKS=0;

-- Set autocommit off for better batching
SET autocommit=0;

-- Use bulk insert optimization
SET SESSION sql_log_bin=0;

-- Show current settings
SHOW VARIABLES LIKE 'foreign_key_checks';
SHOW VARIABLES LIKE 'unique_checks';
EOF
```

### Step 3: Import the SQL File

```bash
# Method 1: Direct pipe (fastest)
docker exec -i tst_db mariadb -uroot -ptest test < /path/to/your/large-file.sql

# Method 2: With timing and progress (if pv is installed)
pv /path/to/your/large-file.sql | docker exec -i tst_db mariadb -uroot -ptest test

# Method 3: Using source command (if file is in container)
docker cp /path/to/your/large-file.sql tst_db:/tmp/import.sql
docker exec tst_db mariadb -uroot -ptest test -e "SOURCE /tmp/import.sql"
```

### Step 4: Re-enable Safety Checks

```bash
# After import completes, re-enable checks
docker exec -i tst_db mariadb -uroot -ptest test << 'EOF'
SET FOREIGN_KEY_CHECKS=1;
SET UNIQUE_CHECKS=1;
COMMIT;

-- Optimize tables after bulk insert
OPTIMIZE TABLE travels;
OPTIMIZE TABLE clustered_arrival_travels32;
OPTIMIZE TABLE points;

-- Update statistics for query optimizer
ANALYZE TABLE travels;
ANALYZE TABLE clustered_arrival_travels32;
ANALYZE TABLE points;
EOF
```

### Step 5: Switch Back to Query-Optimized Configuration

```bash
# Restore original performance config
cp .docker/tst_db/conf.d/performance.cnf.backup .docker/tst_db/conf.d/performance.cnf

# Restart container
docker compose restart tst_db

# Verify
sleep 3
docker exec tst_db mariadb -uroot -ptest -e "SHOW VARIABLES LIKE 'innodb_doublewrite';"
# Should show: ON (re-enabled for safety)
```

## Key Configuration Differences

| Setting | performance.cnf | bulk-import-4gb.cnf | Impact |
|---------|----------------|---------------------|---------|
| innodb_buffer_pool_size | 4GB | 8GB | 2x more memory cache |
| innodb_log_file_size | 512MB | 2GB | 4x fewer checkpoints |
| innodb_flush_log_at_trx_commit | 2 | 0 | ~2x faster writes |
| innodb_doublewrite | ON | OFF | ~2x faster (unsafe!) |
| tmp_table_size | 512MB | 2GB | 4x larger temp ops |
| sort_buffer_size | 16MB | 64MB | 4x faster sorts |

**Expected Import Speed Improvements**: 5-10x faster than defaults

## Advanced: SQL-Level Optimizations

### For mysqldump Files

If you're importing a mysqldump file, add these flags when creating the dump:

```bash
# Create optimized dump
mysqldump -uroot -p \
  --quick \
  --single-transaction \
  --skip-lock-tables \
  --disable-keys \
  --extended-insert \
  --max_allowed_packet=1G \
  persedimai > dump.sql
```

### For Custom SQL Scripts

Add these at the beginning of your SQL file:

```sql
SET FOREIGN_KEY_CHECKS=0;
SET UNIQUE_CHECKS=0;
SET autocommit=0;
SET SESSION sql_log_bin=0;

-- Your INSERT statements here
-- Use extended inserts: INSERT INTO table VALUES (...), (...), (...);

COMMIT;
SET FOREIGN_KEY_CHECKS=1;
SET UNIQUE_CHECKS=1;
```

## Troubleshooting Large Imports

### Import Fails with "Packet Too Large" Error

```bash
# Increase max_allowed_packet (already 1G in bulk-import-4gb.cnf)
docker exec tst_db mariadb -uroot -ptest -e "SET GLOBAL max_allowed_packet=1073741824;"
```

### Out of Memory During Import

```bash
# Check current memory usage
docker stats tst_db --no-stream

# Reduce buffer pool in bulk-import-4gb.cnf
# For 8GB RAM system: innodb_buffer_pool_size = 5G
# For 16GB RAM system: innodb_buffer_pool_size = 10G

# Ensure Docker has enough memory allocated (Docker Desktop settings)
```

### Import is Still Slow

```bash
# Verify settings are applied
docker exec tst_db mariadb -uroot -ptest -e "
SHOW VARIABLES WHERE Variable_name IN (
  'innodb_buffer_pool_size',
  'innodb_log_file_size',
  'innodb_flush_log_at_trx_commit',
  'innodb_doublewrite',
  'foreign_key_checks'
);"

# Check if import is CPU or I/O bound
docker stats tst_db --no-stream

# Monitor import progress
docker exec tst_db mariadb -uroot -ptest -e "SHOW PROCESSLIST;"
```

### Container Crashes During Import

```bash
# Check logs
docker logs tst_db --tail 100

# Common causes:
# 1. innodb_buffer_pool_size too large (reduce it)
# 2. innodb_log_file_size too large (reduce to 1G)
# 3. Not enough disk space for data files
df -h  # Check disk space
```

## Performance Monitoring During Import

```bash
# Real-time monitoring script
watch -n 5 'docker exec tst_db mariadb -uroot -ptest -e "
SELECT
  table_name,
  table_rows,
  ROUND((data_length + index_length)/1024/1024, 2) AS size_mb
FROM information_schema.tables
WHERE table_schema = \"test\"
ORDER BY (data_length + index_length) DESC;"'

# Monitor InnoDB status
docker exec tst_db mariadb -uroot -ptest -e "SHOW ENGINE INNODB STATUS\G" | less

# Check buffer pool usage during import
docker exec tst_db mariadb -uroot -ptest -e "
SHOW STATUS WHERE Variable_name LIKE 'Innodb_buffer_pool%';"
```

## System Requirements for 4GB Import

**Minimum**:
- 8GB RAM
- 20GB free disk space (for data + temp files)
- SSD storage (highly recommended)

**Recommended**:
- 16GB+ RAM
- 50GB+ free disk space
- NVMe SSD storage
- Docker allocated: 10GB+ memory

## Estimated Import Times

Based on testing with 167MB → 20s with performance.cnf:

| File Size | Default Config | performance.cnf | bulk-import-4gb.cnf |
|-----------|---------------|-----------------|---------------------|
| 167MB     | 112s          | 20s             | ~15s                |
| 1GB       | ~670s (11m)   | ~120s (2m)      | ~60s (1m)           |
| 4GB       | ~2680s (45m)  | ~480s (8m)      | ~240s (4m)          |

*Note: Times vary based on hardware, data complexity, and indexes*

## Safety Warnings

⚠️ **IMPORTANT**: The bulk-import-4gb.cnf disables `innodb_doublewrite`!

- This makes imports faster but **increases data corruption risk**
- Only use for:
  - Initial database setup
  - Test environments
  - Situations where you can re-import if needed
- **Always switch back to performance.cnf after import**
- **Always backup before large imports**

## Complete Import Script Example

```bash
#!/bin/bash
# complete-import.sh - Full import workflow

set -e  # Exit on error

IMPORT_FILE="$1"
DB_NAME="test"

echo "Step 1: Switching to bulk import configuration..."
cp .docker/tst_db/conf.d/performance.cnf .docker/tst_db/conf.d/performance.cnf.backup
cp .docker/tst_db/conf.d/bulk-import-4gb.cnf .docker/tst_db/conf.d/performance.cnf
chmod 644 .docker/tst_db/conf.d/performance.cnf
docker compose restart tst_db
sleep 5

echo "Step 2: Preparing database..."
docker exec tst_db mariadb -uroot -ptest $DB_NAME -e "
SET FOREIGN_KEY_CHECKS=0;
SET UNIQUE_CHECKS=0;
"

echo "Step 3: Importing $IMPORT_FILE..."
time docker exec -i tst_db mariadb -uroot -ptest $DB_NAME < "$IMPORT_FILE"

echo "Step 4: Re-enabling safety checks and optimizing..."
docker exec tst_db mariadb -uroot -ptest $DB_NAME -e "
SET FOREIGN_KEY_CHECKS=1;
SET UNIQUE_CHECKS=1;
COMMIT;
ANALYZE TABLE travels;
ANALYZE TABLE clustered_arrival_travels32;
ANALYZE TABLE points;
"

echo "Step 5: Switching back to query-optimized configuration..."
cp .docker/tst_db/conf.d/performance.cnf.backup .docker/tst_db/conf.d/performance.cnf
docker compose restart tst_db
sleep 5

echo "Import complete! Verifying configuration..."
docker exec tst_db mariadb -uroot -ptest -e "SHOW VARIABLES LIKE 'innodb_doublewrite';"

echo "Done!"
```

Usage:
```bash
chmod +x complete-import.sh
./complete-import.sh /path/to/large-file.sql
```
