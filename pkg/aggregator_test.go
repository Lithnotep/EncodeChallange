package pkg

import (
	"testing"
)

func TestNewAggregator(t *testing.T) {
	mapping := URLMapping{
		"http://bit.ly/test": "https://example.com/",
	}

	aggregator := NewAggregator(mapping)

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

	aggregator := NewAggregator(mapping)

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

	aggregator := NewAggregator(mapping)

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

	aggregator := NewAggregator(mapping)

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

	aggregator := NewAggregator(mapping)

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

	aggregator := NewAggregator(mapping)

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

func TestExtractDate(t *testing.T) {
	testCases := []struct {
		timestamp string
		expected  string
		shouldErr bool
	}{
		{"2020-02-15T00:00:00Z", "2020-02-15", false},
		{"2021-12-31T23:59:59Z", "2021-12-31", false},
		{"invalid-timestamp", "", true},
		{"2020-02-15", "", true}, // Missing time part
	}

	for _, tc := range testCases {
		result, err := extractDate(tc.timestamp)

		if tc.shouldErr {
			if err == nil {
				t.Errorf("Expected error for timestamp %s, got nil", tc.timestamp)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for timestamp %s: %v", tc.timestamp, err)
			}
			if result != tc.expected {
				t.Errorf("Expected date %s for timestamp %s, got %s", tc.expected, tc.timestamp, result)
			}
		}
	}
}
