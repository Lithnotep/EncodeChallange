package main

import (
	"os"
	"testing"

	"github.com/Lithnotep/EncodeChallange/pkg"
)

// Integration test for the main workflow
func TestMainWorkflow(t *testing.T) {
	// Create test CSV file
	testCSV := `long_url,domain,hash
https://google.com/,bit.ly,31Tt55y
https://github.com/,bit.ly,2kJO0qS`

	csvFile, err := os.CreateTemp("", "test_encodes*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp CSV file: %v", err)
	}
	defer os.Remove(csvFile.Name())
	defer csvFile.Close()

	if _, err := csvFile.WriteString(testCSV); err != nil {
		t.Fatalf("Failed to write to CSV file: %v", err)
	}
	csvFile.Close()

	// Create test JSON file
	testJSON := `[
		{"bitlink": "http://bit.ly/31Tt55y", "user_agent": "Mozilla/5.0", "timestamp": "2020-02-15T00:00:00Z", "referrer": "t.co", "remote_ip": "4.14.247.63"},
		{"bitlink": "http://bit.ly/2kJO0qS", "user_agent": "Chrome", "timestamp": "2020-02-16T00:00:00Z", "referrer": "direct", "remote_ip": "192.168.1.1"},
		{"bitlink": "http://bit.ly/31Tt55y", "user_agent": "Safari", "timestamp": "2020-02-17T00:00:00Z", "referrer": "facebook.com", "remote_ip": "3.3.3.3"}
	]`

	jsonFile, err := os.CreateTemp("", "test_decodes*.json")
	if err != nil {
		t.Fatalf("Failed to create temp JSON file: %v", err)
	}
	defer os.Remove(jsonFile.Name())
	defer jsonFile.Close()

	if _, err := jsonFile.WriteString(testJSON); err != nil {
		t.Fatalf("Failed to write to JSON file: %v", err)
	}
	jsonFile.Close()

	// Test the main workflow
	// Step 1: Read encodes mapping
	mapping, err := pkg.ReadEncodesMappings(csvFile.Name())
	if err != nil {
		t.Fatalf("ReadEncodesMappings failed: %v", err)
	}

	if len(mapping) != 2 {
		t.Errorf("Expected 2 mappings, got %d", len(mapping))
	}

	// Step 2: Create aggregator
	config := pkg.AggregationConfig{FilterYear: 0} // No year filter for testing
	aggregator := pkg.NewAggregator(mapping, config)
	if aggregator == nil {
		t.Fatal("NewAggregator returned nil")
	}

	// Step 3: Stream and process decodes
	err = pkg.StreamDecodes(jsonFile.Name(), aggregator.ProcessRecord)
	if err != nil {
		t.Fatalf("StreamDecodes failed: %v", err)
	}

	// Step 4: Verify results
	results := aggregator.GetResults()

	if results.TotalClicks != 3 {
		t.Errorf("Expected 3 total clicks, got %d", results.TotalClicks)
	}

	if results.ProcessedRecords != 3 {
		t.Errorf("Expected 3 processed records, got %d", results.ProcessedRecords)
	}

	// Check URL aggregation (31Tt55y should have 2 clicks, 2kJO0qS should have 1)
	googleClicks := results.ClicksByURL["https://google.com/"]
	if googleClicks != 2 {
		t.Errorf("Expected 2 clicks for google.com, got %d", googleClicks)
	}

	githubClicks := results.ClicksByURL["https://github.com/"]
	if githubClicks != 1 {
		t.Errorf("Expected 1 click for github.com, got %d", githubClicks)
	}

	// Check referrer aggregation
	if results.ClicksByReferrer["t.co"] != 1 {
		t.Errorf("Expected 1 click from t.co, got %d", results.ClicksByReferrer["t.co"])
	}

	if results.ClicksByReferrer["direct"] != 1 {
		t.Errorf("Expected 1 click from direct, got %d", results.ClicksByReferrer["direct"])
	}

	if results.ClicksByReferrer["facebook.com"] != 1 {
		t.Errorf("Expected 1 click from facebook.com, got %d", results.ClicksByReferrer["facebook.com"])
	}

	// Check date aggregation
	if results.ClicksByDate["2020-02-15"] != 1 {
		t.Errorf("Expected 1 click on 2020-02-15, got %d", results.ClicksByDate["2020-02-15"])
	}

	if results.ClicksByDate["2020-02-16"] != 1 {
		t.Errorf("Expected 1 click on 2020-02-16, got %d", results.ClicksByDate["2020-02-16"])
	}

	if results.ClicksByDate["2020-02-17"] != 1 {
		t.Errorf("Expected 1 click on 2020-02-17, got %d", results.ClicksByDate["2020-02-17"])
	}

	// No unknown bitlinks in this test
	if len(results.UnknownBitlinks) != 0 {
		t.Errorf("Expected 0 unknown bitlinks, got %d", len(results.UnknownBitlinks))
	}
}

// Test the main workflow with unknown bitlinks
func TestMainWorkflow_WithUnknownBitlinks(t *testing.T) {
	// Create test CSV file with limited mappings
	testCSV := `long_url,domain,hash
https://google.com/,bit.ly,31Tt55y`

	csvFile, err := os.CreateTemp("", "test_limited_encodes*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp CSV file: %v", err)
	}
	defer os.Remove(csvFile.Name())
	defer csvFile.Close()

	if _, err := csvFile.WriteString(testCSV); err != nil {
		t.Fatalf("Failed to write to CSV file: %v", err)
	}
	csvFile.Close()

	// Create test JSON file with some unknown bitlinks
	testJSON := `[
		{"bitlink": "http://bit.ly/31Tt55y", "user_agent": "Mozilla/5.0", "timestamp": "2020-02-15T00:00:00Z", "referrer": "t.co", "remote_ip": "4.14.247.63"},
		{"bitlink": "http://bit.ly/unknown1", "user_agent": "Chrome", "timestamp": "2020-02-16T00:00:00Z", "referrer": "direct", "remote_ip": "192.168.1.1"},
		{"bitlink": "http://bit.ly/unknown2", "user_agent": "Safari", "timestamp": "2020-02-17T00:00:00Z", "referrer": "facebook.com", "remote_ip": "3.3.3.3"}
	]`

	jsonFile, err := os.CreateTemp("", "test_unknown_decodes*.json")
	if err != nil {
		t.Fatalf("Failed to create temp JSON file: %v", err)
	}
	defer os.Remove(jsonFile.Name())
	defer jsonFile.Close()

	if _, err := jsonFile.WriteString(testJSON); err != nil {
		t.Fatalf("Failed to write to JSON file: %v", err)
	}
	jsonFile.Close()

	// Process the data
	mapping, err := pkg.ReadEncodesMappings(csvFile.Name())
	if err != nil {
		t.Fatalf("ReadEncodesMappings failed: %v", err)
	}

	config := pkg.AggregationConfig{FilterYear: 0} // No year filter for testing
	aggregator := pkg.NewAggregator(mapping, config)
	err = pkg.StreamDecodes(jsonFile.Name(), aggregator.ProcessRecord)
	if err != nil {
		t.Fatalf("StreamDecodes failed: %v", err)
	}

	results := aggregator.GetResults()

	// Should track 2 unknown bitlinks
	if len(results.UnknownBitlinks) != 2 {
		t.Errorf("Expected 2 unknown bitlinks, got %d", len(results.UnknownBitlinks))
	}

	// Total clicks should still be 3
	if results.TotalClicks != 3 {
		t.Errorf("Expected 3 total clicks, got %d", results.TotalClicks)
	}

	// Known URL should have 1 click
	if results.ClicksByURL["https://google.com/"] != 1 {
		t.Errorf("Expected 1 click for google.com, got %d", results.ClicksByURL["https://google.com/"])
	}

	// Unknown bitlinks should have 1 click each as fallback
	if results.ClicksByURL["http://bit.ly/unknown1"] != 1 {
		t.Errorf("Expected 1 click for unknown1, got %d", results.ClicksByURL["http://bit.ly/unknown1"])
	}

	if results.ClicksByURL["http://bit.ly/unknown2"] != 1 {
		t.Errorf("Expected 1 click for unknown2, got %d", results.ClicksByURL["http://bit.ly/unknown2"])
	}
}
