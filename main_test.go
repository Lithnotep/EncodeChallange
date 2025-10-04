package main

import (
	"os"
	"reflect"
	"testing"
)

func TestReadDecodesFile(t *testing.T) {
	// Create a temporary JSON file for testing
	testJSON := `[
		{"bitlink": "http://bit.ly/2kkAHNs", "user_agent": "Mozilla/5.0", "timestamp": "2020-02-15T00:00:00Z", "referrer": "t.co", "remote_ip": "4.14.247.63"},
		{"bitlink": "http://bit.ly/31Tt55y", "user_agent": "Chrome", "timestamp": "2020-02-16T00:00:00Z", "referrer": "direct", "remote_ip": "192.168.1.1"}
	]`

	tmpFile, err := os.CreateTemp("", "test*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(testJSON); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Test reading the JSON file
	records, err := readDecodesFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("readDecodesFile failed: %v", err)
	}

	expected := []DecodeRecord{
		{Bitlink: "http://bit.ly/2kkAHNs", UserAgent: "Mozilla/5.0", Timestamp: "2020-02-15T00:00:00Z", Referrer: "t.co", RemoteIP: "4.14.247.63"},
		{Bitlink: "http://bit.ly/31Tt55y", UserAgent: "Chrome", Timestamp: "2020-02-16T00:00:00Z", Referrer: "direct", RemoteIP: "192.168.1.1"},
	}

	if !reflect.DeepEqual(records, expected) {
		t.Errorf("Expected %+v, got %+v", expected, records)
	}
}

func TestReadEncodesFile(t *testing.T) {
	// Create a temporary CSV file for testing
	testCSV := `long_url,domain,hash
https://google.com/,bit.ly,31Tt55y
https://github.com/,bit.ly,2kJO0qS`

	tmpFile, err := os.CreateTemp("", "test*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(testCSV); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Test reading the CSV file
	records, err := readEncodesFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("readEncodesFile failed: %v", err)
	}

	expected := []EncodeRecord{
		{LongURL: "https://google.com/", Domain: "bit.ly", Hash: "31Tt55y"},
		{LongURL: "https://github.com/", Domain: "bit.ly", Hash: "2kJO0qS"},
	}

	if !reflect.DeepEqual(records, expected) {
		t.Errorf("Expected %+v, got %+v", expected, records)
	}
}
