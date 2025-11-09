package analyzer

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

// SaleRecord represents a single sales record from the CSV file.
type SaleRecord struct {
	Date     string  // Transaction date (e.g., "2023-10-25")
	Product  string  // Name of the product sold
	Quantity int     // Quantity of units sold
	Price    float64 // Price per unit
}

// AnalysisResult contains the aggregated results of the data analysis.
type AnalysisResult struct {
	TotalTransactions    int
	TotalRevenue         float64
	MostPopularProduct   string // Name of the product with the highest total quantity sold
	MaxQuantitySoldUnits int    // Max quantity of units sold for that product
}

// --- CSV PROCESSING FUNCTIONS ---

// ParseCSV reads data from the file specified by filePath,
// parses it into SaleRecord structs, and returns a slice of valid records.
// It skips records with invalid data (e.g., non-numeric quantity or price)
// and handles file I/O errors.
func ParseCSV(filePath string) ([]SaleRecord, error) {
	// 1. Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %w", filePath, err)
	}
	// The defer statement ensures that the file is closed when the function exits.
	defer file.Close()

	reader := csv.NewReader(file)
	// Setting FieldsPerRecord to -1 allows reading records with a variable number of fields.
	// This lets our code handle missing/extra fields gracefully instead of relying on a critical CSV error.
	reader.FieldsPerRecord = -1
	records := []SaleRecord{}

	// Skip the header (first line)
	if _, err := reader.Read(); err != nil && err != io.EOF {
		return nil, fmt.Errorf("error reading CSV header: %w", err)
	}

	// 2. Iterate and parse rows
	for lineNumber := 2; ; lineNumber++ { // Start from line 2 (after header)
		row, err := reader.Read()
		if err == io.EOF {
			break // End of file reached
		}
		if err != nil {
			// Handle critical I/O errors that are not EOF
			return nil, fmt.Errorf("critical error reading line %d: %w", lineNumber, err)
		}

		// Ensure the record has the expected 4 columns. If not, skip it.
		if len(row) != 4 {
			fmt.Printf("Warning: line %d skipped - incorrect number of fields (%d instead of 4)\n", lineNumber, len(row))
			continue
		}

		// Parse Quantity (string to int)
		quantity, err := strconv.Atoi(row[2])
		if err != nil {
			fmt.Printf("Warning: line %d skipped - error parsing Quantity (%s): %v\n", lineNumber, row[2], err)
			continue
		}

		// Parse Price (string to float64)
		price, err := strconv.ParseFloat(row[3], 64)
		if err != nil {
			// Note: strconv.ParseFloat handles "NaN" and "Inf" without error.
			// We only skip on true parsing errors (e.g., "invalid_price").
			fmt.Printf("Warning: line %d skipped - error parsing Price (%s): %v\n", lineNumber, row[3], err)
			continue
		}

		// Successful record creation only occurs after all validation checks pass.
		record := SaleRecord{
			Date:     row[0],
			Product:  row[1],
			Quantity: quantity,
			Price:    price,
		}
		records = append(records, record)
	}

	return records, nil
}

// --- DATA ANALYSIS FUNCTION ---

// AnalyzeData performs analysis on the provided slice of SaleRecord structs.
// It calculates the total revenue and identifies the most popular product
// based on the total quantity of units sold.
func AnalyzeData(records []SaleRecord) AnalysisResult {
	result := AnalysisResult{
		TotalTransactions:    len(records),
		TotalRevenue:         0.0,
		MaxQuantitySoldUnits: 0,
	}

	// productQuantities aggregates the total units sold for each product.
	productQuantities := make(map[string]int)

	// 1. Calculate Total Revenue and aggregate quantity
	for _, record := range records {
		// Calculate Total Revenue (Quantity * Price)
		result.TotalRevenue += float64(record.Quantity) * record.Price

		// Aggregate quantity sold for each product
		productQuantities[record.Product] += record.Quantity
	}

	// 2. Determine the Most Popular Product
	for product, totalQuantity := range productQuantities {
		if totalQuantity > result.MaxQuantitySoldUnits {
			result.MaxQuantitySoldUnits = totalQuantity
			result.MostPopularProduct = product
		}
	}

	return result
}
