package salesmodels

import "time"

type Sale struct {
	SaleID            int       `json:"sale_id" db:"sale_id"`
	Date              time.Time `json:"date" db:"date"`
	ItemID            int       `json:"item_id" db:"item_id"`
	Quantity          int       `json:"quantity" db:"quantity"`
	PricePerUnit      float64   `json:"price_per_unit" db:"price_per_unit"`
	TotalPrice        float64   `json:"total_price" db:"total_price"`
	TransactionNumber string    `json:"transaction_number" db:"transaction_number"`
	CustomerName      *string   `json:"customer_name,omitempty" db:"customer_name"`
	CustomerPhone     *string   `json:"customer_phone,omitempty" db:"customer_phone"`
	CustomerEmail     *string   `json:"customer_email,omitempty" db:"customer_email"`
	SoldBy            *string   `json:"sold_by,omitempty" db:"sold_by"`
	Notes             *string   `json:"notes,omitempty" db:"notes"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`

	// Additional fields for API responses
	ItemPartNumber  string `json:"item_part_number,omitempty" db:"item_part_number"`
	ItemDescription string `json:"item_description,omitempty" db:"item_description"`
	CategoryName    string `json:"category_name,omitempty" db:"category_name"`
}

type SaleFilter struct {
	ItemID            *int       `query:"item_id"`
	StartDate         *time.Time `query:"start_date"`
	EndDate           *time.Time `query:"end_date"`
	CustomerName      *string    `query:"customer_name"`
	CustomerPhone     *string    `query:"customer_phone"`
	CustomerEmail     *string    `query:"customer_email"`
	TransactionNumber *string    `query:"transaction_number"`
	SoldBy            *string    `query:"sold_by"`
}
