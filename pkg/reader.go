package pkg

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// DecodeRecord represents a click event from the decodes.json file
type DecodeRecord struct {
	Bitlink   string `json:"bitlink"`
	UserAgent string `json:"user_agent"`
	Timestamp string `json:"timestamp"`
	Referrer  string `json:"referrer"`
	RemoteIP  string `json:"remote_ip"`
}

// EncodeRecord represents a URL mapping from the encodes.csv file
type EncodeRecord struct {
	LongURL string
	Domain  string
	Hash    string
}

// URLMapping stores the mapping from bitlink to original URL
type URLMapping map[string]string

// ReadEncodesMappings reads the CSV file and creates a hash map for O(1) lookups
func ReadEncodesMappings(filename string) (URLMapping, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening encodes file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %w", err)
	}

	mapping := make(URLMapping)
	// Skip header row
	for i := 1; i < len(records); i++ {
		if len(records[i]) >= 3 {
			domain := records[i][1]  // e.g., "bit.ly"
			hash := records[i][2]    // e.g., "31Tt55y"
			longURL := records[i][0] // e.g., "https://google.com/"

			// Create the full bitlink URL
			bitlink := fmt.Sprintf("http://%s/%s", domain, hash)
			mapping[bitlink] = longURL
		}
	}

	return mapping, nil
}

// StreamDecodes processes the JSON file using streaming decoder
// The callback function is called for each decode record as it's read
func StreamDecodes(filename string, callback func(DecodeRecord) error) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening decodes file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	// Read the opening bracket of the JSON array
	token, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("error reading JSON array start: %w", err)
	}
	if delim, ok := token.(json.Delim); !ok || delim != '[' {
		return fmt.Errorf("expected JSON array, got %v", token)
	}

	// Process each record in the array
	recordCount := 0
	for decoder.More() {
		var record DecodeRecord
		if err := decoder.Decode(&record); err != nil {
			return fmt.Errorf("error decoding record %d: %w", recordCount, err)
		}

		// Call the callback function for each record
		if err := callback(record); err != nil {
			return fmt.Errorf("error in callback for record %d: %w", recordCount, err)
		}

		recordCount++
	}

	// Read the closing bracket of the JSON array
	token, err = decoder.Token()
	if err != nil && err != io.EOF {
		return fmt.Errorf("error reading JSON array end: %w", err)
	}

	return nil
}

// GetLongURL looks up the original URL for a given bitlink
func (mapping URLMapping) GetLongURL(bitlink string) (string, bool) {
	longURL, exists := mapping[bitlink]
	return longURL, exists
}
