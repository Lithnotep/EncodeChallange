# EncodeChallenge

A simple Go program that aggregates data from CSV and JSON files.

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

# Run tests to verify everything works
go test -v

# Run the program
go run main.go
```

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

## Project Structure

```
EncodeChallange/
├── go.mod              # Go module definition
├── main.go             # Main program logic
├── main_test.go        # Test suite
├── Makefile           # Build automation (optional)
├── README.md          # This file
└── data/              # Sample data files
    ├── data.csv       # Sample CSV file
    └── data.json      # Sample JSON file
```

## Usage

The program reads CSV and JSON files from the `data/` directory and aggregates them by ID. 

- JSON files should contain records with `id` and `name` fields
- CSV files should contain records with `id` and `value` fields
- The program outputs aggregated data combining both sources

## Testing

Run the test suite to ensure everything works correctly:

```bash
go test -v
```

All tests should pass before using the program.

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
