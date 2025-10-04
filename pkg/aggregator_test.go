package pkg

import (
	"testing"
)

func TestNewAggregator(t *testing.T) {
	mapping := URLMapping{
		"http://bit.ly/test": "https://example.com/",
	}

	config := AggregationConfig{FilterYear: 0}
	aggregator := NewAggregator(mapping, config)

	if aggregator == nil {
		t.Fatal("NewAggregator returned nil")
	}

	if aggregator.results.TotalClicks != 0 {
		t.Errorf("Expected initial TotalClicks to be 0, got %d", aggregator.results.TotalClicks)
	}

	if len(aggregator.results.ClicksByURL) != 0 {
		t.Errorf("Expected initial ClicksByURL to be empty, got %v", aggregator.results.ClicksByURL)
	}
}

func TestAggregator_ProcessRecord(t *testing.T) {
	mapping := URLMapping{
		"http://bit.ly/31Tt55y": "https://google.com/",
		"http://bit.ly/2kJO0qS": "https://github.com/",
	}

	config := AggregationConfig{FilterYear: 0}
	aggregator := NewAggregator(mapping, config)

	// Test record with known bitlink
	record1 := DecodeRecord{
		Bitlink:   "http://bit.ly/31Tt55y",
		UserAgent: "Mozilla/5.0",
		Timestamp: "2020-02-15T00:00:00Z",
		Referrer:  "t.co",
		RemoteIP:  "4.14.247.63",
	}

	err := aggregator.ProcessRecord(record1)
	if err != nil {
		t.Fatalf("ProcessRecord failed: %v", err)
	}

	// Check results
	if aggregator.results.TotalClicks != 1 {
		t.Errorf("Expected TotalClicks to be 1, got %d", aggregator.results.TotalClicks)
	}

	if aggregator.results.ProcessedRecords != 1 {
		t.Errorf("Expected ProcessedRecords to be 1, got %d", aggregator.results.ProcessedRecords)
	}

	if aggregator.results.ClicksByURL["https://google.com/"] != 1 {
		t.Errorf("Expected clicks for google.com to be 1, got %d", aggregator.results.ClicksByURL["https://google.com/"])
	}

	if aggregator.results.ClicksByReferrer["t.co"] != 1 {
		t.Errorf("Expected clicks from t.co to be 1, got %d", aggregator.results.ClicksByReferrer["t.co"])
	}

	if aggregator.results.ClicksByDate["2020-02-15"] != 1 {
		t.Errorf("Expected clicks on 2020-02-15 to be 1, got %d", aggregator.results.ClicksByDate["2020-02-15"])
	}
}

func TestAggregator_ProcessRecord_UnknownBitlink(t *testing.T) {
	mapping := URLMapping{
		"http://bit.ly/known": "https://known.com/",
	}

	config := AggregationConfig{FilterYear: 0}
	aggregator := NewAggregator(mapping, config)

	// Test record with unknown bitlink
	record := DecodeRecord{
		Bitlink:   "http://bit.ly/unknown",
		UserAgent: "Mozilla/5.0",
		Timestamp: "2020-02-15T00:00:00Z",
		Referrer:  "direct",
		RemoteIP:  "1.1.1.1",
	}

	err := aggregator.ProcessRecord(record)
	if err != nil {
		t.Fatalf("ProcessRecord failed: %v", err)
	}

	// Check that unknown bitlink is tracked
	if len(aggregator.results.UnknownBitlinks) != 1 {
		t.Errorf("Expected 1 unknown bitlink, got %d", len(aggregator.results.UnknownBitlinks))
	}

	if aggregator.results.UnknownBitlinks[0] != "http://bit.ly/unknown" {
		t.Errorf("Expected unknown bitlink to be 'http://bit.ly/unknown', got %s", aggregator.results.UnknownBitlinks[0])
	}

	// Check that it still counts as a click using the bitlink as fallback
	if aggregator.results.ClicksByURL["http://bit.ly/unknown"] != 1 {
		t.Errorf("Expected clicks for unknown bitlink to be 1, got %d", aggregator.results.ClicksByURL["http://bit.ly/unknown"])
	}
}

func TestAggregator_ProcessRecord_InvalidTimestamp(t *testing.T) {
	mapping := URLMapping{
		"http://bit.ly/test": "https://test.com/",
	}

	config := AggregationConfig{FilterYear: 0}
	aggregator := NewAggregator(mapping, config)

	// Test record with invalid timestamp
	record := DecodeRecord{
		Bitlink:   "http://bit.ly/test",
		UserAgent: "Mozilla/5.0",
		Timestamp: "invalid-timestamp",
		Referrer:  "direct",
		RemoteIP:  "1.1.1.1",
	}

	err := aggregator.ProcessRecord(record)
	if err == nil {
		t.Error("Expected error for invalid timestamp, got nil")
	}
}

func TestAggregator_MultipleRecords(t *testing.T) {
	mapping := URLMapping{
		"http://bit.ly/google": "https://google.com/",
		"http://bit.ly/github": "https://github.com/",
	}

	config := AggregationConfig{FilterYear: 0}
	aggregator := NewAggregator(mapping, config)

	records := []DecodeRecord{
		{Bitlink: "http://bit.ly/google", UserAgent: "Chrome", Timestamp: "2020-01-01T00:00:00Z", Referrer: "direct", RemoteIP: "1.1.1.1"},
		{Bitlink: "http://bit.ly/google", UserAgent: "Firefox", Timestamp: "2020-01-01T00:00:00Z", Referrer: "t.co", RemoteIP: "2.2.2.2"},
		{Bitlink: "http://bit.ly/github", UserAgent: "Safari", Timestamp: "2020-01-02T00:00:00Z", Referrer: "direct", RemoteIP: "3.3.3.3"},
	}

	for _, record := range records {
		err := aggregator.ProcessRecord(record)
		if err != nil {
			t.Fatalf("ProcessRecord failed: %v", err)
		}
	}

	// Check aggregated results
	if aggregator.results.TotalClicks != 3 {
		t.Errorf("Expected TotalClicks to be 3, got %d", aggregator.results.TotalClicks)
	}

	if aggregator.results.ClicksByURL["https://google.com/"] != 2 {
		t.Errorf("Expected clicks for google.com to be 2, got %d", aggregator.results.ClicksByURL["https://google.com/"])
	}

	if aggregator.results.ClicksByURL["https://github.com/"] != 1 {
		t.Errorf("Expected clicks for github.com to be 1, got %d", aggregator.results.ClicksByURL["https://github.com/"])
	}

	if aggregator.results.ClicksByReferrer["direct"] != 2 {
		t.Errorf("Expected clicks from direct to be 2, got %d", aggregator.results.ClicksByReferrer["direct"])
	}

	if aggregator.results.ClicksByDate["2020-01-01"] != 2 {
		t.Errorf("Expected clicks on 2020-01-01 to be 2, got %d", aggregator.results.ClicksByDate["2020-01-01"])
	}

	if aggregator.results.ClicksByDate["2020-01-02"] != 1 {
		t.Errorf("Expected clicks on 2020-01-02 to be 1, got %d", aggregator.results.ClicksByDate["2020-01-02"])
	}
}

func TestAggregator_GetResults(t *testing.T) {
	mapping := URLMapping{
		"http://bit.ly/test": "https://test.com/",
	}

	config := AggregationConfig{FilterYear: 0}
	aggregator := NewAggregator(mapping, config)

	record := DecodeRecord{
		Bitlink:   "http://bit.ly/test",
		UserAgent: "Mozilla/5.0",
		Timestamp: "2020-01-01T00:00:00Z",
		Referrer:  "direct",
		RemoteIP:  "1.1.1.1",
	}

	aggregator.ProcessRecord(record)
	results := aggregator.GetResults()

	if results.TotalClicks != 1 {
		t.Errorf("Expected TotalClicks to be 1, got %d", results.TotalClicks)
	}

	if results.ProcessedRecords != 1 {
		t.Errorf("Expected ProcessedRecords to be 1, got %d", results.ProcessedRecords)
	}
}

// Test year filtering functionality
func TestAggregator_YearFiltering(t *testing.T) {
	mapping := URLMapping{
		"http://bit.ly/test2020": "https://example2020.com/",
		"http://bit.ly/test2021": "https://example2021.com/",
	}

	// Test filtering for 2020 only
	config := AggregationConfig{FilterYear: 2020}
	aggregator := NewAggregator(mapping, config)

	records := []DecodeRecord{
		{Bitlink: "http://bit.ly/test2020", UserAgent: "Mozilla", Timestamp: "2020-05-01T00:00:00Z", Referrer: "direct", RemoteIP: "1.1.1.1"},
		{Bitlink: "http://bit.ly/test2021", UserAgent: "Chrome", Timestamp: "2021-05-01T00:00:00Z", Referrer: "direct", RemoteIP: "2.2.2.2"},
		{Bitlink: "http://bit.ly/test2020", UserAgent: "Safari", Timestamp: "2020-06-01T00:00:00Z", Referrer: "t.co", RemoteIP: "3.3.3.3"},
	}

	for _, record := range records {
		err := aggregator.ProcessRecord(record)
		if err != nil {
			t.Fatalf("ProcessRecord failed: %v", err)
		}
	}

	results := aggregator.GetResults()

	// Should only count 2020 records (2 out of 3)
	if results.TotalClicks != 2 {
		t.Errorf("Expected 2 total clicks (2020 only), got %d", results.TotalClicks)
	}

	if results.FilteredOut != 1 {
		t.Errorf("Expected 1 filtered out record, got %d", results.FilteredOut)
	}

	if results.ProcessedRecords != 3 {
		t.Errorf("Expected 3 processed records, got %d", results.ProcessedRecords)
	}

	// Should only have clicks for 2020 URL
	if results.ClicksByURL["https://example2020.com/"] != 2 {
		t.Errorf("Expected 2 clicks for 2020 URL, got %d", results.ClicksByURL["https://example2020.com/"])
	}

	if results.ClicksByURL["https://example2021.com/"] != 0 {
		t.Errorf("Expected 0 clicks for 2021 URL, got %d", results.ClicksByURL["https://example2021.com/"])
	}
}

// Test the isShortlink method - now based on actual unmapped bitlinks
func TestAggregator_IsShortlink(t *testing.T) {
	mapping := URLMapping{
		"http://bit.ly/mapped": "https://mapped.com/",
	}
	config := AggregationConfig{FilterYear: 0}
	aggregator := NewAggregator(mapping, config)

	// Process some records to populate the UnknownBitlinks list
	records := []DecodeRecord{
		// This will be mapped
		{Bitlink: "http://bit.ly/mapped", UserAgent: "Chrome", Timestamp: "2020-01-01T00:00:00Z", Referrer: "direct", RemoteIP: "1.1.1.1"},
		// These will be unmapped and added to UnknownBitlinks
		{Bitlink: "http://bit.ly/unknown1", UserAgent: "Firefox", Timestamp: "2020-01-02T00:00:00Z", Referrer: "t.co", RemoteIP: "2.2.2.2"},
		{Bitlink: "http://es.pn/unknown2", UserAgent: "Safari", Timestamp: "2020-01-03T00:00:00Z", Referrer: "direct", RemoteIP: "3.3.3.3"},
		{Bitlink: "http://amzn.to/unknown3", UserAgent: "Edge", Timestamp: "2020-01-04T00:00:00Z", Referrer: "facebook.com", RemoteIP: "4.4.4.4"},
	}

	for _, record := range records {
		err := aggregator.ProcessRecord(record)
		if err != nil {
			t.Fatalf("ProcessRecord failed: %v", err)
		}
	}

	// Test that unmapped bitlinks are identified as shortlinks
	if !aggregator.isShortlink("http://bit.ly/unknown1") {
		t.Error("Expected unknown1 to be identified as shortlink")
	}

	if !aggregator.isShortlink("http://es.pn/unknown2") {
		t.Error("Expected unknown2 to be identified as shortlink")
	}

	if !aggregator.isShortlink("http://amzn.to/unknown3") {
		t.Error("Expected unknown3 to be identified as shortlink")
	}

	// Test that mapped long URLs are NOT identified as shortlinks
	if aggregator.isShortlink("https://mapped.com/") {
		t.Error("Expected mapped.com to NOT be identified as shortlink")
	}

	// Test that URLs not in our data are NOT identified as shortlinks
	if aggregator.isShortlink("http://bit.ly/notprocessed") {
		t.Error("Expected unprocessed bitlink to NOT be identified as shortlink")
	}

	if aggregator.isShortlink("https://google.com/") {
		t.Error("Expected google.com to NOT be identified as shortlink")
	}

	// Test edge cases
	if aggregator.isShortlink("") {
		t.Error("Expected empty string to NOT be identified as shortlink")
	}
}

// Test that the final summary only includes mapped long URLs
func TestAggregator_FinalSummaryFiltering(t *testing.T) {
	mapping := URLMapping{
		"http://bit.ly/google": "https://google.com/",
		"http://bit.ly/github": "https://github.com/",
	}

	config := AggregationConfig{FilterYear: 0}
	aggregator := NewAggregator(mapping, config)

	records := []DecodeRecord{
		// Mapped bitlinks (should appear in final summary)
		{Bitlink: "http://bit.ly/google", UserAgent: "Chrome", Timestamp: "2020-01-01T00:00:00Z", Referrer: "direct", RemoteIP: "1.1.1.1"},
		{Bitlink: "http://bit.ly/github", UserAgent: "Firefox", Timestamp: "2020-01-02T00:00:00Z", Referrer: "t.co", RemoteIP: "2.2.2.2"},
		{Bitlink: "http://bit.ly/google", UserAgent: "Safari", Timestamp: "2020-01-03T00:00:00Z", Referrer: "direct", RemoteIP: "3.3.3.3"},

		// Unmapped bitlinks (should NOT appear in final summary)
		{Bitlink: "http://bit.ly/unknown", UserAgent: "Edge", Timestamp: "2020-01-04T00:00:00Z", Referrer: "direct", RemoteIP: "4.4.4.4"},
		{Bitlink: "http://es.pn/unknown", UserAgent: "Opera", Timestamp: "2020-01-05T00:00:00Z", Referrer: "facebook.com", RemoteIP: "5.5.5.5"},
	}

	for _, record := range records {
		err := aggregator.ProcessRecord(record)
		if err != nil {
			t.Fatalf("ProcessRecord failed: %v", err)
		}
	}

	results := aggregator.GetResults()

	// Verify total counts include all records
	if results.TotalClicks != 5 {
		t.Errorf("Expected 5 total clicks, got %d", results.TotalClicks)
	}

	// Verify that mapped URLs have correct counts
	if results.ClicksByURL["https://google.com/"] != 2 {
		t.Errorf("Expected 2 clicks for google.com, got %d", results.ClicksByURL["https://google.com/"])
	}

	if results.ClicksByURL["https://github.com/"] != 1 {
		t.Errorf("Expected 1 click for github.com, got %d", results.ClicksByURL["https://github.com/"])
	}

	// Verify that unmapped URLs are also tracked (as fallback)
	if results.ClicksByURL["http://bit.ly/unknown"] != 1 {
		t.Errorf("Expected 1 click for unknown bitlink, got %d", results.ClicksByURL["http://bit.ly/unknown"])
	}

	if results.ClicksByURL["http://es.pn/unknown"] != 1 {
		t.Errorf("Expected 1 click for es.pn unknown, got %d", results.ClicksByURL["http://es.pn/unknown"])
	}

	// Test that isShortlink correctly identifies the URLs
	if !aggregator.isShortlink("http://bit.ly/unknown") {
		t.Error("Expected bit.ly URL to be identified as shortlink")
	}

	if !aggregator.isShortlink("http://es.pn/unknown") {
		t.Error("Expected es.pn URL to be identified as shortlink")
	}

	if aggregator.isShortlink("https://google.com/") {
		t.Error("Expected google.com to NOT be identified as shortlink")
	}

	if aggregator.isShortlink("https://github.com/") {
		t.Error("Expected github.com to NOT be identified as shortlink")
	}
}

// Test that PrintSummary works without panicking (integration-style test)
func TestAggregator_PrintSummary(t *testing.T) {
	mapping := URLMapping{
		"http://bit.ly/test": "https://example.com/",
	}

	config := AggregationConfig{FilterYear: 2020}
	aggregator := NewAggregator(mapping, config)

	record := DecodeRecord{
		Bitlink:   "http://bit.ly/test",
		UserAgent: "Mozilla/5.0",
		Timestamp: "2020-01-01T00:00:00Z",
		Referrer:  "direct",
		RemoteIP:  "1.1.1.1",
	}

	err := aggregator.ProcessRecord(record)
	if err != nil {
		t.Fatalf("ProcessRecord failed: %v", err)
	}

	// This test mainly ensures PrintSummary doesn't panic
	// In a more sophisticated test setup, we could capture stdout to verify output
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintSummary panicked: %v", r)
		}
	}()

	aggregator.PrintSummary()
}
