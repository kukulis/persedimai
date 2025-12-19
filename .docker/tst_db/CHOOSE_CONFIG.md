# Choose the Right Configuration for Your Server

## Quick Start: Check Your Server RAM

```bash
# Check total RAM
free -h

# Check available RAM for Docker
docker info | grep -i memory
```

## Configuration Selection Guide

Choose the configuration based on your server's **total RAM**:

| Server RAM | Configuration File | Buffer Pool | Best For |
|------------|-------------------|-------------|----------|
| 2-4 GB | performance.cnf | 4GB | ⚠️ May fail, use adaptive |
| 4-6 GB | bulk-import-adaptive.cnf | 2GB | Small-medium imports |
| 8+ GB | bulk-import-4gb.cnf.back | 8GB | Large 4GB+ imports |
| 16+ GB | bulk-import-4gb.cnf.back (edit to 10-12GB) | 10-12GB | Very large imports |

## Current Error: Out of Memory (8GB buffer pool on small server)

Your server doesn't have enough RAM for the 8GB buffer pool in `bulk-import-4gb.cnf.back`.

**Solution: Use the adaptive configuration**

```bash
# On the server, check RAM first
free -h
# Look at the "total" column under "Mem:"

# If you have 4-6GB RAM, use adaptive config:
cd ~/persedimai
rm .docker/tst_db/conf.d/bulk-import-4gb.cnf  # Remove if it exists
cp .docker/tst_db/conf.d/bulk-import-adaptive.cnf .docker/tst_db/conf.d/performance.cnf
chmod 644 .docker/tst_db/conf.d/performance.cnf

# Restart container
docker compose down
docker compose up -d tst_db

# Verify it started successfully
docker logs tst_db
docker exec tst_db mariadb -uroot -ptest -e "SHOW VARIABLES LIKE 'innodb_buffer_pool_size';"
# Should show: 2147483648 (2GB)
```

## If You Have 8GB+ RAM but Still Get Error

This could be due to Docker memory limits. Increase Docker's memory allocation:

### For Docker Desktop (Windows/Mac)
1. Open Docker Desktop Settings
2. Go to Resources → Advanced
3. Increase Memory to at least 10GB
4. Apply & Restart

### For Linux Docker
```bash
# Check current cgroup limits
docker info | grep -i memory

# If no limit shown, check system memory
free -h

# Make sure you have at least 10GB free
# Then use bulk-import-4gb.cnf.back with 8GB buffer pool
```

## Performance Comparison

Based on 167MB SQL import test (112s default → 20s optimized):

| Configuration | Buffer Pool | 167MB Import | 4GB Import (est.) | RAM Required |
|--------------|-------------|--------------|-------------------|--------------|
| Default (no config) | 128MB | 112s | ~45min | 1GB+ |
| performance.cnf | 4GB | 20s | ~8min | 5GB+ |
| bulk-import-adaptive.cnf | 2GB | ~30s | ~12min | 3GB+ |
| bulk-import-4gb.cnf.back | 8GB | ~15s | ~4min | 10GB+ |

## Quick Commands Reference

### Check what config is active
```bash
ls -lh .docker/tst_db/conf.d/
docker exec tst_db mariadb -uroot -ptest -e "SHOW VARIABLES LIKE '%buffer_pool%';"
```

### Switch to performance.cnf (balanced, 4GB)
```bash
cp .docker/tst_db/conf.d/performance.cnf.backup .docker/tst_db/conf.d/performance.cnf 2>/dev/null || echo "Already using performance.cnf"
chmod 644 .docker/tst_db/conf.d/performance.cnf
docker compose restart tst_db
```

### Switch to adaptive (low RAM, 2GB)
```bash
cp .docker/tst_db/conf.d/bulk-import-adaptive.cnf .docker/tst_db/conf.d/performance.cnf
chmod 644 .docker/tst_db/conf.d/performance.cnf
docker compose restart tst_db
```

### Switch to 4GB bulk import (high RAM, 8GB)
```bash
cp .docker/tst_db/conf.d/bulk-import-4gb.cnf.back .docker/tst_db/conf.d/performance.cnf
chmod 644 .docker/tst_db/conf.d/performance.cnf
docker compose restart tst_db
```

## Troubleshooting

### Container keeps crashing with "Out of memory"
→ Your buffer pool size is too large. Use a smaller configuration:
1. Try adaptive config (2GB buffer pool)
2. Or edit performance.cnf and reduce innodb_buffer_pool_size to 1G or 512M

### Container starts but imports are still slow
→ You're using too small a buffer pool. Check your RAM and increase:
```bash
free -h  # Check available RAM
# Edit .docker/tst_db/conf.d/performance.cnf
# Increase innodb_buffer_pool_size to 50-70% of available RAM
docker compose restart tst_db
```

### How to verify current settings
```bash
docker exec tst_db mariadb -uroot -ptest -e "
SELECT
  @@innodb_buffer_pool_size / 1024 / 1024 / 1024 AS buffer_pool_gb,
  @@innodb_log_file_size / 1024 / 1024 AS log_file_mb,
  @@tmp_table_size / 1024 / 1024 AS tmp_table_mb,
  @@innodb_flush_log_at_trx_commit AS flush_mode,
  @@innodb_doublewrite AS doublewrite;
"
```
