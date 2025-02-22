// internal/modules/vehicles/repositories/repository.go
package repositories

import (
	"context"

	vehiclemodels "github.com/hsrvms/autoparts/internal/modules/vehicles/models"
)

// VehicleRepository defines the interface for vehicle database operations
type VehicleRepository interface {
	// Make operations
	GetAllMakes(ctx context.Context) ([]*vehiclemodels.Make, error)
	GetMakeByID(ctx context.Context, id int) (*vehiclemodels.Make, error)
	CreateMake(ctx context.Context, make *vehiclemodels.Make) (int, error)
	UpdateMake(ctx context.Context, make *vehiclemodels.Make) error
	DeleteMake(ctx context.Context, id int) error

	// Model operations
	GetAllModels(ctx context.Context) ([]*vehiclemodels.Model, error)
	GetModelsByMake(ctx context.Context, makeID int) ([]*vehiclemodels.Model, error)
	GetModelByID(ctx context.Context, id int) (*vehiclemodels.Model, error)
	CreateModel(ctx context.Context, model *vehiclemodels.Model) (int, error)
	UpdateModel(ctx context.Context, model *vehiclemodels.Model) error
	DeleteModel(ctx context.Context, id int) error

	// Submodel operations
	GetAllSubmodels(ctx context.Context) ([]*vehiclemodels.Submodel, error)
	GetSubmodelsByModel(ctx context.Context, modelID int) ([]*vehiclemodels.Submodel, error)
	GetSubmodelByID(ctx context.Context, id int) (*vehiclemodels.Submodel, error)
	CreateSubmodel(ctx context.Context, submodel *vehiclemodels.Submodel) (int, error)
	UpdateSubmodel(ctx context.Context, submodel *vehiclemodels.Submodel) error
	DeleteSubmodel(ctx context.Context, id int) error
}
