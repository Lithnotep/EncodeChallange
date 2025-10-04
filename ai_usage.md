
# AI Usage Documentation

This document outlines how AI tools (GitHub Copilot and Claude) were used throughout the development of the EncodeChallenge project.

## Development Timeline

### 1. Initial Project Setup
**Tool Used:** GitHub Copilot  
**Objective:** Create a foundational Go project structure  

- Generated simple Go skeleton for JSON/CSV reading
- Set up testing framework for immediate TDD workflow
- Ensured project can be easily run when pulled down
- Added sample data files (`decodes.json` and `encodes.csv`)

### 2. Architecture Planning & Streaming Implementation
**Tool Used:** GitHub Copilot  
**Objective:** Design memory-efficient data processing  

**Key Focus:** *"Use JSON streaming decoder instead of loading all 10k records into memory"*

- Implemented streaming JSON processing for large datasets
- Avoided memory-intensive approaches that load entire files
- Designed for scalability with larger data volumes

### 3. Modular Code Organization  
**Tool Used:** Claude 4  
**Objective:** Separate concerns into testable modules

- Created separate `reader` and `aggregator` packages
- Developed comprehensive test suites for each module
- Implemented Test-Driven Development (TDD) approach
- **Result:** Claude effectively created packages with complete aggregation understanding

**Note:** Significant manual code review and validation was performed to ensure correct implementation direction.

### 4. Dynamic Command-Line Interface
**Tool Used:** Claude 4  
**Objective:** Add flexible year-based filtering  

**Prompt:** *"I want to be able to run a command in the terminal that will run the aggregator with dynamic settings. First option: I want to be able to send a specific year and have it get the clicks for that year. It should default to 2021 if no year is sent."*

**Implementation:**
- Added command-line flag support with Go's `flag` package
- Default year: 2021 (as per project requirements)
- Dynamic year filtering: any year can be queried
- Special case: `year=0` processes ALL years
- **Result:** Complete solution meeting project requirements

### 5. Final Output Formatting
**Tool Used:** Claude 4  
**Objective:** Create clean JSON-formatted output  

**Prompt:** *"I want the final log to be a super simple breakdown of the long URLs and the amount of clicks for the specified year. Right before this log, note that shortlinks without mapping are left out."*

**Example Format Requested:**
```json
[{"https://google.com": 3}, {"https://www.twitter.com": 2}]
```

**Implementation:**
- Clean JSON-like output format
- Automatic filtering of unmapped shortlinks
- Clear notification about excluded data
- Dynamic shortlink detection based on actual unmapped data

### 6. Data-Driven Shortlink Detection
**Tool Used:** Claude 4  
**Objective:** Replace hardcoded shortlink domains with dynamic detection

**Enhancement:** Instead of hardcoding shortlink domains (bit.ly, es.pn, etc.), the system now dynamically identifies shortlinks based on actual unmapped bitlinks encountered in the data.

**Benefits:**
- Adapts to new shortlink services automatically
- Works with any data source
- No maintenance of hardcoded domain lists

### 7. Performance Monitoring & Timing Implementation
**Tool Used:** Claude 4  
**Objective:** Add processing time measurement for scalability validation

**Rationale:** *"I added [a process timer] so that if a massive file were used I could see if the time was scaling linearly and not going exponential."*

**Implementation:**
- Added `ProcessingTime` field to `AggregationResults` struct  
- Created `StartTiming()` and `StopTiming()` methods on aggregator
- Integrated timing around the core streaming process in main.go
- Display processing time in final results summary

**Performance Results Observed:**
- **10,000 records**: ~37ms average processing time
- **Processing rate**: ~270,000 records/second
- **Scaling**: Linear time complexity confirmed (O(n))

**Benefits:**
- **Scalability validation**: Confirms streaming approach scales linearly
- **Performance monitoring**: Easy to detect performance regressions
- **Optimization guidance**: Provides metrics for future improvements
- **Production readiness**: Essential monitoring for large dataset processing

### 8. Comprehensive Sorting Architecture & Code Consolidation
**Tool Used:** Claude 4  
**Objective:** Implement configurable sorting across all summary sections with architectural optimization

**Initial Request:** *"The final output should be in descending order of clicks"*  
**Follow-up Analysis:** *"Question, is this the best and most efficient place to sort? If there was a flag for ascending or descending sort should it be done in the aggregator?"*

**Evolution of Implementation:**

**Phase 1: Basic Final Output Sorting**
- Added descending sort to final summary JSON output only
- Sorting logic embedded directly in `PrintSummary()` method
- **Issue:** Mixed presentation logic with data transformation

**Phase 2: Architecture Analysis & Comprehensive Sorting**
- **Problem Identified:** Multiple unsorted sections (Top URLs, Top Referrers, Clicks by Date)
- **Request:** *"I also noticed Top Referrers, Clicks by Date, Top URLs by Clicks are not sorted either. Have them use the config new config as well. As efficiently as possible. If ascending then all lists are ascending is the idea."*

**Implementation:**
- Added `SortDesc` boolean to `AggregationConfig` (default: true)
- Added `-sort-desc` CLI flag with clear documentation
- Created separate sorting methods:
  - `GetSortedURLs()` - with shortlink filtering capability
  - `getSortedKeyValues()` - for generic key-value sorting
- Applied consistent sorting to ALL summary sections

**Phase 3: Code Consolidation Optimization**
- **Key Insight:** *"Question. Couldn't getSortedKeyValues work for urls too? Since it's a key value string to int?"*
- **Analysis:** Both methods performed identical operations (map[string]int → sorted slice)
- **Refactoring:** Unified into single `getSortedKeyValues()` method with optional filter function

**Final Architecture:**
```go
// One unified method handles all sorting with optional filtering
getSortedKeyValues(data map[string]int, filter func(string) bool) []KeyValue

// URL sorting uses same method with shortlink filtering
GetSortedURLs(excludeShortlinks bool) []KeyValue
```

**CLI Usage:**
```bash
go run main.go                           # Default: 2021, descending
go run main.go -sort-desc=false          # Ascending sort
go run main.go -year=2020 -sort-desc=false # Year 2020, ascending
```

**Results Achieved:**
- **Code Reduction:** Eliminated ~20 lines of duplicate sorting logic
- **Consistency:** All sections use identical sorting behavior
- **Maintainability:** Single point of change for sorting modifications
- **Flexibility:** Filter function enables custom exclusion logic
- **Testing:** Added comprehensive sorting tests (25 total tests passing)
- **Performance:** No performance impact, sorting happens post-aggregation

**Benefits:**
- **DRY Principle:** No duplicate sorting code
- **Separation of Concerns:** Data processing separated from presentation
- **Configurability:** User controls sort direction across entire application
- **Architectural Cleanliness:** Unified approach reduces cognitive complexity
- **Future-Proof:** Easy to extend with additional sorting criteria

