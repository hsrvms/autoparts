package suppliermodels

import "time"

type Supplier struct {
	SupplierID    int       `json:"supplier_id" db:"supplier_id"`
	Name          string    `json:"name" db:"name"`
	ContactPerson *string   `json:"contact_person,omitempty" db:"contact_person"`
	Phone         *string   `json:"phone,omitempty" db:"phone"`
	Email         *string   `json:"email,omitempty" db:"email"`
	Address       *string   `json:"address,omitempty" db:"address"`
	TaxID         *string   `json:"tax_id,omitempty" db:"tax_id"`
	PaymentTerms  *string   `json:"payment_terms,omitempty" db:"payment_terms"`
	Notes         *string   `json:"notes,omitempty" db:"notes"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// Filter represents the search criteria for suppliers
type SupplierFilter struct {
	SearchTerm     *string `query:"search"`
	HasActiveItems *bool   `query:"has_active_items"`
}
