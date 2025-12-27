# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Go-based travel route search system that finds flight connections with multiple transfers (persėdimai means "transfers" in Lithuanian). The system searches for optimal travel paths between locations with configurable numbers of transfers, using MySQL/MariaDB for data storage and query optimization.

## Common Commands

### Build & Run

```bash
# Navigate to application directory
cd application

# Tidy dependencies
go mod tidy

# Build the web application
go build -o bin/webapp ./cmd/webapp

# Run the web application
bin/webapp

# Build and run the seeder
go build -o bin/seeder ./cmd/seeder
bin/seeder -env test -strategy normal

# Build clusters creator
go build -o bin/createclusters ./cmd/createclusters

# Build and run schedule collector (Aviation Edge API)
go build -o bin/collectschedules ./cmd/collectschedules
bin/collectschedules -country US -start 2025-12-27
```

### Database Management

```bash
# Start Docker containers (MariaDB databases)
docker-compose up -d

# Production database: localhost:23313 (user: persedimai, db: persedimai)
# Test database: localhost:23314 (user: test, db: test)

# Dump clustered data
mysqldump -P 23314 -u root -h 127.0.0.1 -p test > clusters_32.sql
```

### Testing

```bash
# Run all tests (excludes drafttests)
cd application
go test ./...

# Run a single test
cd integration_tests  # or drafttests, performance_tests
go test -run TestLoadDbConfig

# Run with draft tests (special tag)
go test -tags=draft ./...

# Run tests with no timeout (for long-running tests)
go test -v -timeout 0 -run TestClustersCreator

# Run benchmarks
cd performance_tests
go test -bench=. -benchmem
```

## Architecture

### Core Components

**Strategy Pattern for Travel Search**
- `TravelSearchStrategy` interface defines search contracts
- `SimpleTravelSearchStrategy` - Direct SQL joins for 1-3 transfers
- `ClusteredTravelSearchStrategy` - Uses pre-computed clustered data for faster 2-4 transfer searches
- Time-based clustering groups travels into 1-hour buckets for optimization

**Database Layer**
- `Database` - Connection wrapper with version detection (MySQL vs MariaDB)
- `TravelDao` - Travel/transfer data access with timeout support
- `PointDao` - Location/point data access
- DAO methods support configurable query timeouts (see Timeout Handling below)

**Data Models**
- `Transfer` - Single flight/travel segment (from, to, departure, arrival)
- `TransferSequence` - Ordered sequence of transfers forming a complete path
- `TravelPath` - UI-friendly representation of a journey
- `TravelFilter` - Search parameters (source, destination, date range, transfer count)

**Dependency Injection**
- `di.InitializeSingletons(env)` - Initializes global instances based on environment
- Environment determined by .env file selection (.env, .env.test, etc.)
- DatabasesContainer pattern allows multiple database connections

**Web Layer**
- Gin framework for HTTP routing
- Controllers: `TravelSearchController`, `FlightsSearchController`, `HomeController`
- API endpoints under `/api` prefix (e.g., `/api/points`)
- HTML templates rendered with template functions (e.g., `add` function)

### Timeout Handling

**Current Implementation (Basic)**:
- `TravelDao` has a `Timeout` field (time.Duration) that can be set per-instance
- `TravelDao.executeQueryWithTimeout()` (internal/dao/travel_dao.go:166) uses Go context to cancel queries client-side
- `Database.AddTimeoutToQuery()` (internal/database/database.go:121) adds server-side timeout hints:
  - MariaDB: `SET STATEMENT max_statement_time=X FOR ...`
  - MySQL: `SELECT /*+ MAX_EXECUTION_TIME(X) */ ...`
- Server-side hints applied to FindPathSimple* and FindPathClustered* methods

**Proposed Enhancement (Not Yet Implemented)**:
- A two-tier timeout strategy is documented in `TIMEOUT_SUMMARY.md`, `TIMEOUT_SOLUTION_PROPOSAL.md`
- Example implementation exists in `internal/database/timeout_example.go.example`
- Proposed design: 15s user timeout (return error to user) + 5min background timeout (continue query)
- Currently NOT in production code - timeout system is simpler than documented proposal

### Directory Structure

```
application/
├── cmd/                          # Executable entry points
│   ├── webapp/                   # Main web application
│   ├── seeder/                   # Database seeding tool
│   ├── createdb/                 # Database creation utility
│   └── createclusters/           # Clustering data generator
├── internal/
│   ├── dao/                      # Data access objects
│   ├── data/                     # Data models and filters
│   ├── database/                 # Database connection layer
│   ├── generator/                # Test data generators
│   ├── migrations/               # Database schema migrations
│   ├── tables/                   # Table/entity definitions
│   ├── travel_finder/            # Search strategy implementations
│   ├── util/                     # Utilities (array ops, env helpers, time)
│   └── web/                      # Web controllers and routing
│       └── api/                  # REST API controllers
├── di/                           # Dependency injection
├── integration_tests/            # Tests requiring database
├── performance_tests/            # Performance benchmarks
├── drafttests/                   # Draft tests (tag: draft)
└── templates/                    # HTML templates
```

### Key Patterns

**Transfer vs Trip vs Travel**
- `Transfer`: Original model name for a single flight segment
- `Trip`: Legacy name, avoid in new code
- Use `Transfer` consistently for single segments and `TransferSequence` for paths

**Clustered Search Optimization**
- Pre-computed table `clustered_arrival_travels32` groups travels by hour
- Cluster ID = `floor(unix_timestamp(date)/3600)`
- Dramatically faster for multi-transfer searches on large datasets
- Falls back to regular `travels` table for simple searches

**Environment Configuration**
- Production: `.env` file
- Test: `.env.test` file
- Environment files contain DBUSER, DBPASS, DBNAME, DBPORT, DBHOST
- DI factory handles env file resolution (checks current dir, then parent dir)

**Singleton Pattern**
- `DatabaseInstance` - Default database connection
- `DatabasesContainerInstance` - Multi-database container
- `ApiPointsControllerInstance` - Shared API controller
- Initialized via `di.InitializeSingletons(env)` in web router setup

## Code Style Notes

- String escaping: Use `database.MysqlRealEscapeString()` for dynamic SQL (though parameterized queries are preferred)
- Date parsing: `util.ParseDate()` and `util.ParseDateTime()` for consistent formats
- Array operations: Use `util.ArrayMap`, `util.ArrayFilter`, `util.ArrayReduce` for functional transforms
- Error handling: Return errors up the stack; log.Fatal only in main/cmd packages

## Special Test Tags

Draft tests use build tag `//go:build draft` at the top of test files. These are experimental or long-running tests excluded from normal test runs. Run them explicitly with `-tags=draft`.

## Important Constraints

- Maximum 3 transfers for SimpleTravelSearchStrategy (5 for Clustered)
- Time filtering uses arrival time windows (ArrivalTimeFrom/To)
- Transit wait time controlled by MaxWaitHoursBetweenTransits
- Max connection time hours: 2, 4, 8, 16, or 32 (controlled by MaxConnectionTimeHours)
- Results limited by filter.Limit parameter

## Lessons Learned

### API Integration: Always Verify Endpoint URLs from Official Documentation

When integrating with external APIs (e.g., Aviation Edge API), **never assume or guess endpoint URLs** based on feature names or general descriptions. Always:

1. **Consult the official API documentation** for each specific endpoint
2. **Verify the exact URL path** and required parameters
3. **Test with the documented examples** before writing production code

**Example from Aviation Edge API integration:**
- ❌ **Wrong assumption:** Feature named "Historical Schedules API" → guessed endpoint `/scheduleDatabase` (returned 404)
- ❌ **Wrong assumption:** Feature named "Future Schedules API" → guessed endpoint `/schedulesFuture` (returned 404)
- ✅ **Correct approach:** Read documentation → actual endpoints are `/flightsHistory` and `/flightsFuture`

**Cost of assumption:** Hours of debugging, incorrect error diagnosis (thinking it was an API tier limitation when it was just wrong URLs), and wasted API calls.

**Best practice:** When building API clients, always cross-reference with official documentation pages for each endpoint, not just the general overview page.
