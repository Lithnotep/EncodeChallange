package pkg

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestReadEncodesMappings(t *testing.T) {
	// Test CSV content
	testCSV := `long_url,domain,hash
https://google.com/,bit.ly,31Tt55y
https://github.com/,bit.ly,2kJO0qS
https://twitter.com/,bit.ly,2kkAHNs`

	// Create temporary CSV file
	tmpFile, err := os.CreateTemp("", "test_encodes*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(testCSV); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Test the function
	mapping, err := ReadEncodesMappings(tmpFile.Name())
	if err != nil {
		t.Fatalf("ReadEncodesMappings failed: %v", err)
	}

	// Expected mapping
	expected := URLMapping{
		"http://bit.ly/31Tt55y": "https://google.com/",
		"http://bit.ly/2kJO0qS": "https://github.com/",
		"http://bit.ly/2kkAHNs": "https://twitter.com/",
	}

	if !reflect.DeepEqual(mapping, expected) {
		t.Errorf("Expected mapping %+v, got %+v", expected, mapping)
	}
}

func TestReadEncodesMappings_EmptyFile(t *testing.T) {
	// Test with empty CSV (only header)
	testCSV := `long_url,domain,hash`

	tmpFile, err := os.CreateTemp("", "test_empty*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(testCSV); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	mapping, err := ReadEncodesMappings(tmpFile.Name())
	if err != nil {
		t.Fatalf("ReadEncodesMappings failed: %v", err)
	}

	if len(mapping) != 0 {
		t.Errorf("Expected empty mapping, got %+v", mapping)
	}
}

func TestReadEncodesMappings_InvalidFile(t *testing.T) {
	_, err := ReadEncodesMappings("nonexistent_file.csv")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestStreamDecodes(t *testing.T) {
	// Test JSON content
	testJSON := `[
		{"bitlink": "http://bit.ly/2kkAHNs", "user_agent": "Mozilla/5.0", "timestamp": "2020-02-15T00:00:00Z", "referrer": "t.co", "remote_ip": "4.14.247.63"},
		{"bitlink": "http://bit.ly/31Tt55y", "user_agent": "Chrome", "timestamp": "2020-02-16T00:00:00Z", "referrer": "direct", "remote_ip": "192.168.1.1"}
	]`

	// Create temporary JSON file
	tmpFile, err := os.CreateTemp("", "test_decodes*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(testJSON); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Collect records via callback
	var records []DecodeRecord
	err = StreamDecodes(tmpFile.Name(), func(record DecodeRecord) error {
		records = append(records, record)
		return nil
	})

	if err != nil {
		t.Fatalf("StreamDecodes failed: %v", err)
	}

	// Expected records
	expected := []DecodeRecord{
		{Bitlink: "http://bit.ly/2kkAHNs", UserAgent: "Mozilla/5.0", Timestamp: "2020-02-15T00:00:00Z", Referrer: "t.co", RemoteIP: "4.14.247.63"},
		{Bitlink: "http://bit.ly/31Tt55y", UserAgent: "Chrome", Timestamp: "2020-02-16T00:00:00Z", Referrer: "direct", RemoteIP: "192.168.1.1"},
	}

	if !reflect.DeepEqual(records, expected) {
		t.Errorf("Expected records %+v, got %+v", expected, records)
	}
}

func TestStreamDecodes_CallbackError(t *testing.T) {
	testJSON := `[{"bitlink": "http://bit.ly/test", "user_agent": "test", "timestamp": "2020-01-01T00:00:00Z", "referrer": "test", "remote_ip": "1.1.1.1"}]`

	tmpFile, err := os.CreateTemp("", "test_callback_error*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(testJSON); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Test callback that returns an error
	callbackError := fmt.Errorf("callback test error")
	err = StreamDecodes(tmpFile.Name(), func(record DecodeRecord) error {
		return callbackError
	})

	if err == nil {
		t.Error("Expected error from callback, got nil")
	}
}

func TestStreamDecodes_InvalidFile(t *testing.T) {
	err := StreamDecodes("nonexistent_file.json", func(record DecodeRecord) error {
		return nil
	})
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestStreamDecodes_InvalidJSON(t *testing.T) {
	// Test with invalid JSON
	testJSON := `{"invalid": "json structure"}`

	tmpFile, err := os.CreateTemp("", "test_invalid*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(testJSON); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	err = StreamDecodes(tmpFile.Name(), func(record DecodeRecord) error {
		return nil
	})

	if err == nil {
		t.Error("Expected error for invalid JSON structure, got nil")
	}
}

func TestURLMapping_GetLongURL(t *testing.T) {
	mapping := URLMapping{
		"http://bit.ly/31Tt55y": "https://google.com/",
		"http://bit.ly/2kJO0qS": "https://github.com/",
	}

	// Test existing URL
	longURL, exists := mapping.GetLongURL("http://bit.ly/31Tt55y")
	if !exists {
		t.Error("Expected URL to exist in mapping")
	}
	if longURL != "https://google.com/" {
		t.Errorf("Expected 'https://google.com/', got '%s'", longURL)
	}

	// Test non-existing URL
	_, exists = mapping.GetLongURL("http://bit.ly/nonexistent")
	if exists {
		t.Error("Expected URL to not exist in mapping")
	}
}
