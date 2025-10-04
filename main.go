package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
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

func main() {
	fmt.Println("Starting Encode Challenge Data Reading...")

	// Read decode data (JSON)
	decodes, err := readDecodesFile("data/decodes.json")
	if err != nil {
		log.Printf("Error reading decodes file: %v", err)
		return
	}
	fmt.Printf("Read %d decode records\n", len(decodes))

	// Read encode data (CSV)
	encodes, err := readEncodesFile("data/encodes.csv")
	if err != nil {
		log.Printf("Error reading encodes file: %v", err)
		return
	}
	fmt.Printf("Read %d encode records\n", len(encodes))

	// Display sample data
	if len(decodes) > 0 {
		fmt.Printf("Sample decode: %+v\n", decodes[0])
	}
	if len(encodes) > 0 {
		fmt.Printf("Sample encode: %+v\n", encodes[0])
	}
}

// readDecodesFile reads and parses the decodes JSON file
func readDecodesFile(filename string) ([]DecodeRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening decodes file: %w", err)
	}
	defer file.Close()

	var records []DecodeRecord
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&records); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w", err)
	}

	return records, nil
}

// readEncodesFile reads and parses the encodes CSV file
func readEncodesFile(filename string) ([]EncodeRecord, error) {
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

	var encodeRecords []EncodeRecord
	// Skip header row
	for i := 1; i < len(records); i++ {
		if len(records[i]) >= 3 {
			encodeRecords = append(encodeRecords, EncodeRecord{
				LongURL: records[i][0],
				Domain:  records[i][1],
				Hash:    records[i][2],
			})
		}
	}

	return encodeRecords, nil
}
