package main

import (
	"flag"
	"fmt"
	"os"
	"sales-analysis/internal/analyzer"
)

func main() {

	// 1. Argument Declaration and Parsing using the flag package

	// Declare the "file" flag with a default value and usage description.
	filePath := flag.String("file", "data/sales.csv", "Path to the CSV sales data file")

	// Parse the command-line arguments, populating the filePath variable.
	flag.Parse()

	// Check if the file path is empty (though it has a default value, this ensures robustness)
	if *filePath == "" {
		fmt.Println("Error: File path is required.")
		fmt.Println("Usage: go run . --file=<path/to/file.csv>")
		os.Exit(1)
	}

	// 2. File Parsing and Critical Error Handling

	// flag.String returns a *string, so we dereference it using *
	records, err := analyzer.ParseCSV(*filePath)
	if err != nil {
		fmt.Printf("Critical parsing error: %v\n", err)
		os.Exit(1)
	}

	// Check if any valid records were found
	if len(records) == 0 {
		fmt.Printf("File '%s' read successfully, but no valid records were found for analysis.\n", *filePath)
		return
	}

	// 3. Data Analysis
	result := analyzer.AnalyzeData(records)

	// 4. Print Results
	fmt.Println("--- Sales Record Analysis Report ---")
	fmt.Printf("File Processed: %s\n", *filePath)
	fmt.Printf("Total Valid Transactions: %d\n", result.TotalTransactions)
	fmt.Printf("Total Revenue: %.2f $\n", result.TotalRevenue)
	fmt.Printf("Most Popular Product: %s (sold %d units)\n", result.MostPopularProduct, result.MaxQuantitySoldUnits)
	fmt.Println("------------------------------------")
}
