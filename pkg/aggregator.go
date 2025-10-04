package pkg

import (
	"fmt"
	"sort"
	"time"
)

// AggregationConfig holds configuration options for aggregation
type AggregationConfig struct {
	FilterYear int  // Year to filter by (0 means no filter)
	SortDesc   bool // true for descending sort, false for ascending
}

// AggregationResults holds all the computed analytics
type AggregationResults struct {
	TotalClicks      int
	ClicksByURL      map[string]int
	ClicksByReferrer map[string]int
	ClicksByDate     map[string]int // YYYY-MM-DD format
	UnknownBitlinks  []string       // Bitlinks not found in encodes mapping
	ProcessedRecords int
	FilteredOut      int           // Records filtered out by year
	FilterYear       int           // Year that was filtered for
	ProcessingTime   time.Duration // Total time taken for streaming and processing
}

// Aggregator handles the streaming aggregation of decode records
type Aggregator struct {
	mapping   URLMapping
	config    AggregationConfig
	results   AggregationResults
	startTime time.Time // Track when processing started
}

// NewAggregator creates a new aggregator with the URL mapping and configuration
func NewAggregator(mapping URLMapping, config AggregationConfig) *Aggregator {
	return &Aggregator{
		mapping: mapping,
		config:  config,
		results: AggregationResults{
			ClicksByURL:      make(map[string]int),
			ClicksByReferrer: make(map[string]int),
			ClicksByDate:     make(map[string]int),
			UnknownBitlinks:  make([]string, 0),
			FilterYear:       config.FilterYear,
		},
	}
}

// ProcessRecord processes a single decode record and updates aggregations
func (a *Aggregator) ProcessRecord(record DecodeRecord) error {
	a.results.ProcessedRecords++

	// Parse the timestamp to check the year
	recordTime, err := time.Parse("2006-01-02T15:04:05Z", record.Timestamp)
	if err != nil {
		return fmt.Errorf("error parsing timestamp %s: %w", record.Timestamp, err)
	}

	// Filter by year if specified
	if a.config.FilterYear > 0 && recordTime.Year() != a.config.FilterYear {
		a.results.FilteredOut++
		return nil // Skip this record
	}

	a.results.TotalClicks++

	// Look up the original URL
	longURL, found := a.mapping.GetLongURL(record.Bitlink)
	if !found {
		// Track unknown bitlinks for debugging
		a.results.UnknownBitlinks = append(a.results.UnknownBitlinks, record.Bitlink)
		longURL = record.Bitlink // Use bitlink as fallback
	}

	// Aggregate clicks by original URL
	a.results.ClicksByURL[longURL]++

	// Aggregate clicks by referrer
	a.results.ClicksByReferrer[record.Referrer]++

	// Aggregate clicks by date
	date := recordTime.Format("2006-01-02")
	a.results.ClicksByDate[date]++

	return nil
}

// StartTiming begins tracking processing time
func (a *Aggregator) StartTiming() {
	a.startTime = time.Now()
}

// StopTiming ends tracking and stores the total processing time
func (a *Aggregator) StopTiming() {
	if !a.startTime.IsZero() {
		a.results.ProcessingTime = time.Since(a.startTime)
	}
}

// GetResults returns the final aggregation results
func (a *Aggregator) GetResults() AggregationResults {
	return a.results
}

// GetSortedURLs returns URLs sorted by click count according to config
func (a *Aggregator) GetSortedURLs(excludeShortlinks bool) []KeyValue {
	var filter func(string) bool
	if excludeShortlinks {
		filter = func(url string) bool {
			return a.isShortlink(url) // Return true to EXCLUDE shortlinks
		}
	}

	return a.getSortedKeyValues(a.results.ClicksByURL, filter)
}

// KeyValue represents a generic key-value pair for sorting
type KeyValue struct {
	Key   string
	Value int
}

// getSortedKeyValues converts a map to sorted slice based on config
// Optional filter function can be provided to exclude certain keys
func (a *Aggregator) getSortedKeyValues(data map[string]int, filter func(string) bool) []KeyValue {
	var items []KeyValue
	for key, value := range data {
		// Apply filter if provided (return true to EXCLUDE the item)
		if filter != nil && filter(key) {
			continue
		}
		if key != "" {
			items = append(items, KeyValue{Key: key, Value: value})
		}
	}

	// Sort based on configuration
	sort.Slice(items, func(i, j int) bool {
		if a.config.SortDesc {
			return items[i].Value > items[j].Value // Descending
		}
		return items[i].Value < items[j].Value // Ascending
	})

	return items
}

// PrintSummary prints a human-readable summary of the results
func (a *Aggregator) PrintSummary() {
	fmt.Printf("\n=== Aggregation Results ===\n")
	if a.results.FilterYear > 0 {
		fmt.Printf("Filter Year: %d\n", a.results.FilterYear)
		fmt.Printf("Records Filtered Out: %d\n", a.results.FilteredOut)
	}
	fmt.Printf("Total Records Processed: %d\n", a.results.ProcessedRecords)
	fmt.Printf("Total Clicks: %d\n", a.results.TotalClicks)
	fmt.Printf("Unknown Bitlinks: %d\n", len(a.results.UnknownBitlinks))
	if a.results.ProcessingTime > 0 {
		fmt.Printf("Processing Time: %v\n", a.results.ProcessingTime)
	}

	fmt.Printf("\n--- Top URLs by Clicks ---\n")
	sortedURLs := a.GetSortedURLs(false) // Include all URLs
	for _, urlClick := range sortedURLs {
		fmt.Printf("%s: %d clicks\n", urlClick.Key, urlClick.Value)
	}

	fmt.Printf("\n--- Top Referrers ---\n")
	sortedReferrers := a.getSortedKeyValues(a.results.ClicksByReferrer, nil)
	for _, referrer := range sortedReferrers {
		fmt.Printf("%s: %d clicks\n", referrer.Key, referrer.Value)
	}

	fmt.Printf("\n--- Clicks by Date (first 10) ---\n")
	sortedDates := a.getSortedKeyValues(a.results.ClicksByDate, nil)
	for i, date := range sortedDates {
		if i >= 10 {
			break
		}
		fmt.Printf("%s: %d clicks\n", date.Key, date.Value)
	}

	if len(a.results.UnknownBitlinks) > 0 {
		fmt.Printf("\n--- Unknown Bitlink Clicks (first 5) ---\n")
		for i, bitlink := range a.results.UnknownBitlinks {
			if i >= 5 {
				break
			}
			fmt.Printf("%s\n", bitlink)
		}
	}

	// Print final summary - only mapped long URLs (shortlinks without mapping are excluded)
	fmt.Printf("\nNote: Shortlinks without mapping are excluded from the final summary.\n")
	fmt.Printf("\nFinal Summary:\n")

	// Get sorted URLs excluding shortlinks
	sortedFinalURLs := a.GetSortedURLs(true) // Exclude shortlinks

	// Print sorted results
	fmt.Printf("[")
	for i, urlClick := range sortedFinalURLs {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("{\"%s\": %d}", urlClick.Key, urlClick.Value)
	}
	fmt.Printf("]\n")
}

// isShortlink checks if a URL is a shortlink (not a mapped long URL)
// It determines this by checking if the URL appears in our list of unknown bitlinks
func (a *Aggregator) isShortlink(url string) bool {
	// Check if this URL is in our list of unmapped bitlinks
	for _, unmappedLink := range a.results.UnknownBitlinks {
		if url == unmappedLink {
			return true
		}
	}
	return false
}
