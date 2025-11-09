package analyzer_test

import (
	"os"
	"path/filepath"
	"sales-analysis/internal/analyzer"
	"testing"
)

// Helper function to create a temporary CSV file for testing ParseCSV
func createTestFile(t *testing.T, filename string, content string) string {
	t.Helper()
	// Create a temporary directory for the test file
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, filename)

	// Write content to the file
	err := os.WriteFile(filePath, []byte(content), 0666)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	return filePath
}

// TestParseCSV verifies the file reading, error handling, and data parsing logic.
func TestParseCSV(t *testing.T) {
	// Note: The header is required for a valid CSV file.
	header := "Date,Product,Quantity,Price\n"

	tests := []struct {
		name        string
		csvContent  string
		fileName    string
		expectErr   bool
		expectedLen int
		// Check the last record's data integrity
		checkData func(t *testing.T, records []analyzer.SaleRecord)
	}{
		{
			name: "Success_StandardData",
			csvContent: header +
				"2023-10-01,Laptop,2,1200.50\n" +
				"2023-10-02,Mouse,10,25.99\n",
			fileName:    "valid.csv",
			expectErr:   false,
			expectedLen: 2,
			checkData: func(t *testing.T, records []analyzer.SaleRecord) {
				if records[0].Product != "Laptop" || records[1].Quantity != 10 {
					t.Errorf("Parsed data mismatch. Got %v", records)
				}
			},
		},
		{
			name:        "Error_FileNotFound",
			fileName:    "nonexistent.csv",
			expectErr:   true, // Should return a file opening error
			expectedLen: 0,
			checkData:   func(t *testing.T, records []analyzer.SaleRecord) {},
		},
		{
			name: "Warning_BadQuantitySkipped",
			// "Two" is not a number, this record should be skipped with a warning
			csvContent: header +
				"2023-10-01,Laptop,Two,1200.50\n" +
				"2023-10-02,Mouse,10,25.99",
			fileName:    "bad_qty.csv",
			expectErr:   false,
			expectedLen: 1, // Only one record remains (Mouse)
			checkData: func(t *testing.T, records []analyzer.SaleRecord) {
				if records[0].Product != "Mouse" {
					t.Errorf("Bad record was not successfully skipped.")
				}
			},
		},
		{
			name: "Warning_BadPriceSkipped",
			// "INVALID_PRICE" is not a float, this record should be skipped with a warning
			csvContent: header +
				"2023-10-01,Keyboard,5,INVALID_PRICE\n" +
				"2023-10-02,Mouse,10,25.99\n",
			fileName:    "bad_price.csv",
			expectErr:   false,
			expectedLen: 1, // Only one record remains (Mouse)
			checkData:   func(t *testing.T, records []analyzer.SaleRecord) {},
		},
		{
			name: "Warning_WrongColumnCountSkipped",
			csvContent: header +
				"2023-10-01,ProductA,5,10.0,EXTRA_FIELD\n" + // 5 columns
				"2023-10-02,Mouse,10,25.99",
			fileName:    "bad_cols.csv",
			expectErr:   false,
			expectedLen: 1, // Only one record remains (Mouse)
			checkData:   func(t *testing.T, records []analyzer.SaleRecord) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test file if content is provided
			filePath := tt.fileName
			if tt.csvContent != "" {
				filePath = createTestFile(t, tt.fileName, tt.csvContent)
			}

			// Execute the function under test
			records, err := analyzer.ParseCSV(filePath)

			// 1. Check for expected error state
			if (err != nil) != tt.expectErr {
				t.Fatalf("ParseCSV() error = %v, expectErr %v", err, tt.expectErr)
			}

			// If error was expected, stop further checks
			if tt.expectErr {
				return
			}

			// 2. Check the number of parsed records
			if len(records) != tt.expectedLen {
				t.Fatalf("ParseCSV() got %d records, expected %d", len(records), tt.expectedLen)
			}

			// 3. Run specific data integrity checks
			tt.checkData(t, records)
		})
	}
}

// TestAnalyzeData verifies the core business logic: revenue, total count, and most popular product.
func TestAnalyzeData(t *testing.T) {
	// Helper for safe float comparison in tests
	const floatTolerance = 0.0001

	tests := []struct {
		name         string
		inputRecords []analyzer.SaleRecord
		expected     analyzer.AnalysisResult
	}{
		{
			name: "Success_NormalData",
			inputRecords: []analyzer.SaleRecord{
				{"2023", "ProductA", 10, 5.0}, // Rev 50.0
				{"2023", "ProductB", 5, 20.0}, // Rev 100.0
				{"2023", "ProductA", 2, 5.0},  // Rev 10.0
			},
			expected: analyzer.AnalysisResult{
				TotalTransactions:    3,
				TotalRevenue:         160.0,
				MostPopularProduct:   "ProductA", // Total Qty 12
				MaxQuantitySoldUnits: 12,
			},
		},
		{
			name:         "EdgeCase_EmptyData",
			inputRecords: []analyzer.SaleRecord{},
			expected: analyzer.AnalysisResult{
				TotalTransactions:    0,
				TotalRevenue:         0.0,
				MostPopularProduct:   "",
				MaxQuantitySoldUnits: 0,
			},
		},
		{
			name: "Success_SingleRecord",
			inputRecords: []analyzer.SaleRecord{
				{"2023", "ProductZ", 1, 99.99},
			},
			expected: analyzer.AnalysisResult{
				TotalTransactions:    1,
				TotalRevenue:         99.99,
				MostPopularProduct:   "ProductZ",
				MaxQuantitySoldUnits: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.AnalyzeData(tt.inputRecords)

			// 1. Check Total Transactions
			if result.TotalTransactions != tt.expected.TotalTransactions {
				t.Errorf("TotalTransactions got %d, want %d", result.TotalTransactions, tt.expected.TotalTransactions)
			}

			// 2. Check Total Revenue (using tolerance for float comparison)
			if abs(result.TotalRevenue-tt.expected.TotalRevenue) > floatTolerance {
				t.Errorf("TotalRevenue got %.2f, want %.2f", result.TotalRevenue, tt.expected.TotalRevenue)
			}

			// 3. Check Most Popular Product
			if result.MostPopularProduct != tt.expected.MostPopularProduct {
				t.Errorf("MostPopularProduct got %s, want %s", result.MostPopularProduct, tt.expected.MostPopularProduct)
			}

			// 4. Check Max Quantity Sold
			if result.MaxQuantitySoldUnits != tt.expected.MaxQuantitySoldUnits {
				t.Errorf("MaxQuantitySoldUnits got %d, want %d", result.MaxQuantitySoldUnits, tt.expected.MaxQuantitySoldUnits)
			}
		})
	}
}

// Simple absolute value function for float comparison
func abs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}
