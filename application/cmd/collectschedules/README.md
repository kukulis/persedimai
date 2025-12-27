# collectschedules

A command-line tool to collect flight schedules from Aviation Edge API for all airports in a country.

## Prerequisites

1. **Aviation Edge API Key**: Set the `AVIATION_EDGE_API_KEY` environment variable
2. **Go 1.16+**: Required to build the application

## Installation

### Build

```bash
cd application
go build -o bin/collectschedules ./cmd/collectschedules
```

### Configuration

Create or update your `.env` file in the project root:

```bash
AVIATION_EDGE_API_KEY=your-api-key-here
```

Or export the environment variable:

```bash
export AVIATION_EDGE_API_KEY=your-api-key-here
```

## Usage

### Basic Usage

```bash
# Collect schedules for all US airports for today
./bin/collectschedules -country US

# Collect schedules for a date range
./bin/collectschedules -country GB -start 2025-12-20 -end 2025-12-22

# Specify output file
./bin/collectschedules -country FR -output french_schedules.json
```

### Command-Line Options

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-country` | string | **REQUIRED** | Country code (e.g., US, GB, FR, DE) |
| `-start` | string | today | Start date in format YYYY-MM-DD |
| `-end` | string | today | End date in format YYYY-MM-DD |
| `-output` | string | `schedules.json` | Output JSON file path |
| `-departures` | bool | `true` | Include departure schedules |
| `-arrivals` | bool | `true` | Include arrival schedules |
| `-rate-limit` | int | `1` | Delay between API calls in seconds |
| `-print` | bool | `false` | Print to stdout instead of saving to file |

### Examples

#### Example 1: Collect US schedules for a specific date

```bash
./bin/collectschedules \
  -country US \
  -start 2025-12-27 \
  -output us_schedules_dec27.json
```

#### Example 2: Collect only departures for UK

```bash
./bin/collectschedules \
  -country GB \
  -start 2025-12-20 \
  -end 2025-12-22 \
  -arrivals=false \
  -output uk_departures.json
```

#### Example 3: Print schedules to stdout

```bash
./bin/collectschedules \
  -country FR \
  -start 2025-12-27 \
  -print
```

#### Example 4: Collect with slower rate limiting

```bash
./bin/collectschedules \
  -country DE \
  -start 2025-12-25 \
  -end 2025-12-26 \
  -rate-limit 3
```

## How It Works

1. **Loads API Key**: Reads `AVIATION_EDGE_API_KEY` from environment
2. **Fetches Airports**: Retrieves all airports for the specified country
3. **Iterates Dates**: Loops through each date in the range
4. **Collects Schedules**: For each airport and date:
   - Fetches departure schedules (if enabled)
   - Fetches arrival schedules (if enabled)
   - Passes data to the consumer (file or print)
5. **Rate Limiting**: Adds delay between API calls to avoid hitting rate limits
6. **Saves Results**: Writes all collected schedules to JSON file (or prints to stdout)

## Output Format

The output JSON file contains an array of schedule objects:

```json
[
  {
    "type": "departure",
    "status": "scheduled",
    "departure": {
      "iataCode": "JFK",
      "terminal": "4",
      "gate": "B23",
      "scheduledTime": "2025-12-27T14:30:00.000"
    },
    "arrival": {
      "iataCode": "LAX",
      "terminal": "5",
      "scheduledTime": "2025-12-27T17:45:00.000"
    },
    "airline": {
      "name": "American Airlines",
      "iataCode": "AA",
      "icaoCode": "AAL"
    },
    "flight": {
      "number": "100",
      "iataNumber": "AA100",
      "icaoNumber": "AAL100"
    }
  }
]
```

## Performance Considerations

- **API Rate Limits**: Developer accounts have rate limits. Use `-rate-limit` to add delays
- **Large Countries**: Countries with many airports (e.g., US) will take significant time
- **Date Range**: Each additional date multiplies the number of API calls
- **Auto-Flush**: File consumer auto-flushes every 1000 schedules to prevent memory issues

## Troubleshooting

### Error: AVIATION_EDGE_API_KEY environment variable not set

**Solution**: Set the API key in your `.env` file or export it:

```bash
export AVIATION_EDGE_API_KEY=your-api-key
```

### Error: country code is required

**Solution**: Provide the `-country` flag:

```bash
./bin/collectschedules -country US
```

### Warning: .env file not found

**Solution**: Create a `.env` file in the `application` directory or export the API key manually.

### API Rate Limiting Errors

**Solution**: Increase the rate limit delay:

```bash
./bin/collectschedules -country US -rate-limit 3
```

## Common Country Codes

| Code | Country |
|------|---------|
| US | United States |
| GB | United Kingdom |
| FR | France |
| DE | Germany |
| IT | Italy |
| ES | Spain |
| CA | Canada |
| AU | Australia |
| JP | Japan |
| CN | China |

## Notes

- Historical schedules use `GetHistoricalSchedules()` API endpoint
- The tool collects data for **all airports** in the country
- Progress is logged to stdout
- Final statistics are displayed at completion

## See Also

- Aviation Edge API Documentation: https://aviation-edge.com/developers/
- Main package documentation: `application/internal/aviation_edge/README.md`
