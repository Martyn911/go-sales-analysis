package analyzer

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

