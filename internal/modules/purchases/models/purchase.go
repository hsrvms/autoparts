package purchasemodels

import "time"

type Purchase struct {
	PurchaseID    int       `json:"purchase_id" db:"purchase_id"`
	Date          time.Time `json:"date" db:"date"`
	SupplierID    int       `json:"supplier_id" db:"supplier_id"`
	ItemID        int       `json:"item_id" db:"item_id"`
	Quantity      int       `json:"quantity" db:"quantity"`
	CostPerUnit   float64   `json:"cost_per_unit" db:"cost_per_unit"`
	TotalCost     float64   `json:"total_cost" db:"total_cost"`
	InvoiceNumber *string   `json:"invoice_number,omitempty" db:"invoice_number"`
	ReceivedBy    *string   `json:"received_by,omitempty" db:"received_by"`
	Notes         *string   `json:"notes,omitempty" db:"notes"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`

	// Additional fields for API responses
	SupplierName    string `json:"supplier_name,omitempty" db:"supplier_name"`
	ItemPartNumber  string `json:"item_part_number,omitempty" db:"item_part_number"`
	ItemDescription string `json:"item_description,omitempty" db:"item_description"`
}

type PurchaseFilter struct {
	SupplierID    *int       `query:"supplier_id"`
	ItemID        *int       `query:"item_id"`
	StartDate     *time.Time `query:"start_date"`
	EndDate       *time.Time `query:"end_date"`
	InvoiceNumber *string    `query:"invoice_number"`
}
