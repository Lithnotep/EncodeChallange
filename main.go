package main

import (
	"fmt"
	"log"

	"github.com/Lithnotep/EncodeChallange/pkg"
)

func main() {
	fmt.Println("Starting Encode Challenge Data Processing...")

	// Step 1: Build URL mapping index (one-time setup)
	fmt.Println("Loading URL mappings from encodes.csv...")
	mapping, err := pkg.ReadEncodesMappings("data/encodes.csv")
	if err != nil {
		log.Printf("Error reading encodes mapping: %v", err)
		return
	}
	fmt.Printf("Loaded %d URL mappings\n", len(mapping))

	// Step 2: Create aggregator with the mapping
	aggregator := pkg.NewAggregator(mapping)

	// Step 3: Stream process decode records (single pass)
	fmt.Println("Streaming decode records...")
	err = pkg.StreamDecodes("data/decodes.json", aggregator.ProcessRecord)
	if err != nil {
		log.Printf("Error streaming decodes: %v", err)
		return
	}

	// Step 4: Display results
	fmt.Println("Processing complete!")
	aggregator.PrintSummary()
}
