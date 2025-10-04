package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Lithnotep/EncodeChallange/pkg"
)

func main() {
	// Parse command line flags
	var year = flag.Int("year", 2021, "Filter clicks by year (default: 2021)")
	var help = flag.Bool("help", false, "Show usage information")
	flag.Parse()

	if *help {
		fmt.Println("Encode Challenge Data Processing Tool")
		fmt.Println("\nUsage:")
		fmt.Println("  go run main.go [flags]")
		fmt.Println("\nFlags:")
		flag.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  go run main.go              # Filter for year 2021 (default)")
		fmt.Println("  go run main.go -year=2020    # Filter for year 2020")
		fmt.Println("  go run main.go -year=0       # No year filter (all data)")
		return
	}

	fmt.Printf("Starting Encode Challenge Data Processing (Year: %d)...\n", *year)

	// Step 1: Build URL mapping index (one-time setup)
	fmt.Println("Loading URL mappings from encodes.csv...")
	mapping, err := pkg.ReadEncodesMappings("data/encodes.csv")
	if err != nil {
		log.Printf("Error reading encodes mapping: %v", err)
		return
	}
	fmt.Printf("Loaded %d URL mappings\n", len(mapping))

	// Step 2: Create aggregator with the mapping and configuration
	config := pkg.AggregationConfig{
		FilterYear: *year,
	}
	aggregator := pkg.NewAggregator(mapping, config)

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
