package dashboardmodels

import "time"

type Stats struct {
	LowStockCount int     `json:"low_stock_count"`
	TodaySales    float64 `json:"today_sales"`
	ActiveItems   int     `json:"active_items"`
	SupplierCount int     `json:"supplier_count"`
}

type Activity struct {
	Type      string    `json:"type"` // "sale", "purchase", "inventory"
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type LowStockItem struct {
	ItemID       int    `json:"item_id"`
	PartNumber   string `json:"part_number"`
	Description  string `json:"description"`
	CurrentStock int    `json:"current_stock"`
	MinimumStock int    `json:"minimum_stock"`
	Category     string `json:"category"`
}
