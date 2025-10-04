# EncodeChallenge

A high-performance Go program that efficiently processes and aggregates URL shortener click data from CSV and JSON files using streaming techniques and Test-Driven Development (TDD).

## Features

- **Memory-efficient streaming**: Processes 10,000+ records without loading all data into memory
- **Single-pass processing**: Reads decode data only once for optimal performance
- **Year-based filtering**: Filter click data by specific years with command-line arguments
- **Comprehensive analytics**: Click aggregation by URL, referrer, and date
- **Unknown link tracking**: Identifies and reports bitlinks not found in the mapping
- **Test-Driven Development**: 100% test coverage with unit and integration tests

## Prerequisites

- Go 1.21 or higher
- Git (for cloning the repository)

## Installation & Setup

### 1. Install Go

**On macOS (using Homebrew):**
```bash
brew install go
```

**On Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install golang-go
```

**On Windows:**
- Download from https://golang.org/dl/
- Run the installer and follow the setup wizard

**On other systems:**
- Download from https://golang.org/dl/ and follow the installation guide

### 2. Verify Go Installation
```bash
go version
```
Should show Go 1.21 or higher.

### 3. Clone and Run the Project

```bash
# Clone the repository
git clone https://github.com/Lithnotep/EncodeChallange.git
cd EncodeChallange

# Download Go dependencies (if any)
go mod tidy

# Run all tests to verify everything works
go test ./... -v

# Run the program with default settings
go run main.go
```

## Usage

The program analyzes URL shortener click data with flexible command-line options:

### Basic Usage

```bash
# Default: Filter for year 2021, sort descending
go run main.go

# Filter for specific year
go run main.go -year=2020
go run main.go -year=2021

# Process all data (no year filter)
go run main.go -year=0

# Sort results in ascending order (lowest to highest clicks)
go run main.go -sort-desc=false

# Combine options: 2020 data with ascending sort
go run main.go -year=2020 -sort-desc=false

# Show help and available options
go run main.go -help
```

### Command-Line Options

| Flag | Default | Description |
|------|---------|-------------|
| `-year` | 2021 | Filter clicks by year (0 = no filter) |
| `-sort-desc` | true | Sort results in descending order (false = ascending) |
| `-help` | false | Show usage information |

### Data Format

**Input Files:**
- `data/encodes.csv`: URL mappings (long_url, domain, hash)
- `data/decodes.json`: Click events (bitlink, user_agent, timestamp, referrer, remote_ip)

**Sample encodes.csv:**
```csv
long_url,domain,hash
https://google.com/,bit.ly,31Tt55y
https://github.com/,bit.ly,2kJO0qS
```

**Sample decodes.json:**
```json
[
  {
    "bitlink": "http://bit.ly/31Tt55y",
    "user_agent": "Mozilla/5.0...",
    "timestamp": "2021-02-15T00:00:00Z",
    "referrer": "t.co",
    "remote_ip": "4.14.247.63"
  }
]
```

## Project Structure

```
EncodeChallange/
├── go.mod              # Go module definition
├── main.go             # Main program and CLI interface
├── main_test.go        # Integration tests
├── Makefile           # Build automation (optional)
├── README.md          # This file
├── pkg/               # Core packages
│   ├── reader.go      # CSV/JSON streaming readers
│   ├── reader_test.go # Reader unit tests
│   ├── aggregator.go  # Data aggregation logic
│   └── aggregator_test.go # Aggregator unit tests  
└── data/              # Data files
    ├── encodes.csv    # URL mappings
    └── decodes.json   # Click event data
```

## Output Example

```
=== Aggregation Results ===
Filter Year: 2021
Records Filtered Out: 4918
Total Records Processed: 10000
Total Clicks: 5082
Unknown Bitlinks: 2018
Processing Time: 37ms

--- Top URLs by Clicks ---
https://youtube.com/: 557 clicks
https://twitter.com/: 512 clicks
...

--- Top Referrers ---
direct: 2039 clicks
facebook.com: 541 clicks
...

--- Clicks by Date (first 10) ---
2021-12-15: 25 clicks
2021-01-02: 23 clicks
...

Final Summary:
[{"https://youtube.com/": 557}, {"https://twitter.com/": 512}, {"https://reddit.com/": 510}]
```

### Sorting Behavior

The `-sort-desc` flag controls the sort order for **all summary sections**:

- **`-sort-desc=true` (default)**: All results sorted by click count, highest to lowest
- **`-sort-desc=false`**: All results sorted by click count, lowest to highest

**Sections affected by sorting:**
- Top URLs by Clicks
- Top Referrers  
- Clicks by Date
- Final Summary (JSON output)

**Example with ascending sort:**
```bash
go run main.go -sort-desc=false
# Output shows lowest click counts first in all sections
```

## Architecture

### Streaming Processing
- **Memory Efficient**: Uses JSON streaming decoder to process large files
- **Single Pass**: Reads decode data only once
- **Real-time Filtering**: Filters records during streaming, not after

### Modular Design
- **`pkg/reader.go`**: Handles CSV and JSON file reading with streaming
- **`pkg/aggregator.go`**: Processes and aggregates click data with filtering
- **`main.go`**: Command-line interface and orchestration

## Alternative: Using Make (if available)

If you have `make` installed:

```bash
# Run tests
make test

# Build the program
make build

# Run the program
make run

# Clean build artifacts
make clean
```

## Testing

The project follows Test-Driven Development (TDD) with comprehensive test coverage:

```bash
# Run all tests (recommended)
go test ./... -v

# Run tests with coverage
go test ./... -cover

# Run only unit tests (pkg package)
go test ./pkg/... -v

# Run only integration tests (main package)
go test -v
```

### Test Coverage

- **16 total tests** across all packages
- **Unit Tests**: Reader functions, aggregation logic, year filtering
- **Integration Tests**: End-to-end workflow testing
- **Edge Cases**: Invalid files, malformed data, callback errors

All tests should pass before using the program.

## Performance

- **Memory Usage**: Constant memory usage regardless of file size
- **Processing Speed**: ~10,000 records processed efficiently
- **Scalability**: Designed to handle larger datasets without memory issues

### Benchmarks
- **10,000 records**: Processed in ~0.4 seconds
- **Memory footprint**: Minimal (streaming approach)
- **No data duplication**: Single-pass processing

## Development

### Adding New Features

1. **Write tests first** (TDD approach)
2. **Implement functionality** in appropriate package
3. **Update command-line interface** if needed
4. **Update documentation**

### Code Organization

- **Separation of concerns**: I/O operations separate from business logic
- **Testable units**: Each function has focused responsibility
- **Configuration-driven**: Easy to extend with new options

## Troubleshooting

**Go version too old:**
```bash
go version
```
If less than 1.21, update Go using the installation methods above.

**Module issues:**
```bash
go mod tidy
go clean -modcache
```

**Permission issues on Unix systems:**
```bash
chmod +x main.go
```

**Test failures:**
```bash
# Clean test cache and rerun
go clean -testcache
go test ./... -v
```

**Large file processing:**
- The streaming approach handles large files efficiently
- If you encounter memory issues, check for infinite loops in callback functions

## Contributing

1. Fork the repository
2. Create a feature branch
3. Write tests for new functionality (TDD)
4. Implement the feature
5. Ensure all tests pass: `go test ./... -v`
6. Submit a pull request

## License

This project is available under the MIT License.
