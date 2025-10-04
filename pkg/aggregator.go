package pkg

import (
	"fmt"
	"time"
)

// AggregationConfig holds configuration options for aggregation
type AggregationConfig struct {
	FilterYear int // Year to filter by (0 means no filter)
}

// AggregationResults holds all the computed analytics
type AggregationResults struct {
	TotalClicks      int
	ClicksByURL      map[string]int
	ClicksByReferrer map[string]int
	ClicksByDate     map[string]int // YYYY-MM-DD format
	UnknownBitlinks  []string       // Bitlinks not found in encodes mapping
	ProcessedRecords int
	FilteredOut      int // Records filtered out by year
	FilterYear       int // Year that was filtered for
}

// Aggregator handles the streaming aggregation of decode records
type Aggregator struct {
	mapping URLMapping
	config  AggregationConfig
	results AggregationResults
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

// GetResults returns the final aggregation results
func (a *Aggregator) GetResults() AggregationResults {
	return a.results
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

	fmt.Printf("\n--- Top URLs by Clicks ---\n")
	for url, clicks := range a.results.ClicksByURL {
		fmt.Printf("%s: %d clicks\n", url, clicks)
	}

	fmt.Printf("\n--- Top Referrers ---\n")
	for referrer, clicks := range a.results.ClicksByReferrer {
		fmt.Printf("%s: %d clicks\n", referrer, clicks)
	}

	fmt.Printf("\n--- Clicks by Date (first 10) ---\n")
	count := 0
	for date, clicks := range a.results.ClicksByDate {
		if count >= 10 {
			break
		}
		fmt.Printf("%s: %d clicks\n", date, clicks)
		count++
	}

	if len(a.results.UnknownBitlinks) > 0 {
		fmt.Printf("\n--- Unknown Bitlinks (first 5) ---\n")
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
	fmt.Printf("[")
	first := true
	for url, clicks := range a.results.ClicksByURL {
		// Only include actual long URLs - exclude unmapped bitlinks
		// Unmapped bitlinks will still have the bitlink format (bit.ly, es.pn, amzn.to, etc.)
		if url != "" && !a.isShortlink(url) {
			if !first {
				fmt.Printf(", ")
			}
			if first {
				first = false
			}
			fmt.Printf("{\"%s\": %d}", url, clicks)
		}
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
