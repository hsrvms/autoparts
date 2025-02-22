package vehiclemodels

import "time"

// Make represents a vehicle manufacturer
type Make struct {
	MakeID    int       `json:"make_id" db:"make_id"`
	MakeName  string    `json:"make_name" db:"make_name"`
	Country   string    `json:"country" db:"country"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Model represents a vehicle model (like A3, 418, etc.)
type Model struct {
	ModelID   int       `json:"model_id" db:"model_id"`
	MakeID    int       `json:"make_id" db:"make_id"`
	ModelName string    `json:"model_name" db:"model_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// Additional fields for API responses
	MakeName string `json:"make_name,omitempty" db:"-"`
}

// Submodel represents specific variants of a model
type Submodel struct {
	SubmodelID         int       `json:"submodel_id" db:"submodel_id"`
	ModelID            int       `json:"model_id" db:"model_id"`
	SubmodelName       string    `json:"submodel_name" db:"submodel_name"`
	YearFrom           int       `json:"year_from" db:"year_from"`
	YearTo             *int      `json:"year_to,omitempty" db:"year_to"`
	EngineType         string    `json:"engine_type" db:"engine_type"`
	EngineDisplacement float64   `json:"engine_displacement" db:"engine_displacement"`
	FuelType           string    `json:"fuel_type" db:"fuel_type"`
	TransmissionType   string    `json:"transmission_type" db:"transmission_type"`
	BodyType           string    `json:"body_type" db:"body_type"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`

	// Additional fields for API responses
	ModelName string `json:"model_name,omitempty" db:"-"`
	MakeName  string `json:"make_name,omitempty" db:"-"`
}
