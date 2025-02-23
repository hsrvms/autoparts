package inventorymodels

import "time"

type Compatibility struct {
	CompatID   int       `json:"compat_id" db:"compat_id"`
	ItemID     int       `json:"item_id" db:"item_id"`
	SubmodelID int       `json:"submodel_id" db:"submodel_id"`
	Notes      *string   `json:"notes,omitempty" db:"notes"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`

	// Additional fields for API responses
	ModelName    string `json:"model_name,omitempty" db:"-"`
	MakeName     string `json:"make_name,omitempty" db:"-"`
	SubmodelName string `json:"submodel_name,omitempty" db:"-"`
}
