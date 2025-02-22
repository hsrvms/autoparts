package services

import (
	"context"
	"errors"

	vehiclemodels "github.com/hsrvms/autoparts/internal/modules/vehicles/models"
	"github.com/hsrvms/autoparts/internal/modules/vehicles/repositories"
)

var (
	ErrMakeNotFound     = errors.New("make not found")
	ErrModelNotFound    = errors.New("model not found")
	ErrSubmodelNotFound = errors.New("submodel not found")
	ErrInvalidMakeID    = errors.New("invalid make ID")
	ErrInvalidModelID   = errors.New("invalid model ID")
)

type VehicleService interface {
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

type vehicleService struct {
	repo repositories.VehicleRepository
}

func NewVehicleService(repo repositories.VehicleRepository) VehicleService {
	return &vehicleService{
		repo: repo,
	}
}

func (s *vehicleService) GetAllMakes(ctx context.Context) ([]*vehiclemodels.Make, error) {
	return s.repo.GetAllMakes(ctx)
}

func (s *vehicleService) GetMakeByID(ctx context.Context, id int) (*vehiclemodels.Make, error) {
	if id <= 0 {
		return nil, ErrInvalidMakeID
	}

	make, err := s.repo.GetMakeByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if make == nil {
		return nil, ErrMakeNotFound
	}

	return make, nil
}

func (s *vehicleService) CreateMake(ctx context.Context, make *vehiclemodels.Make) (int, error) {
	// Validate required fields
	if make.MakeName == "" {
		return 0, errors.New("make name is required")
	}

	return s.repo.CreateMake(ctx, make)
}

func (s *vehicleService) UpdateMake(ctx context.Context, make *vehiclemodels.Make) error {
	if make.MakeID <= 0 {
		return ErrInvalidMakeID
	}

	if make.MakeName == "" {
		return errors.New("make name is required")
	}

	// Check if make exists
	existing, err := s.repo.GetMakeByID(ctx, make.MakeID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrMakeNotFound
	}

	return s.repo.UpdateMake(ctx, make)
}

func (s *vehicleService) DeleteMake(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrInvalidMakeID
	}

	// Check if make exists
	existing, err := s.repo.GetMakeByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrMakeNotFound
	}

	// Check if make has any models
	models, err := s.repo.GetModelsByMake(ctx, id)
	if err != nil {
		return err
	}
	if len(models) > 0 {
		return errors.New("cannot delete make with existing models")
	}

	return s.repo.DeleteMake(ctx, id)
}

// Model operations
func (s *vehicleService) GetAllModels(ctx context.Context) ([]*vehiclemodels.Model, error) {
	return s.repo.GetAllModels(ctx)
}

func (s *vehicleService) GetModelsByMake(ctx context.Context, makeID int) ([]*vehiclemodels.Model, error) {
	if makeID <= 0 {
		return nil, ErrInvalidMakeID
	}

	// Verify make exists
	make, err := s.repo.GetMakeByID(ctx, makeID)
	if err != nil {
		return nil, err
	}
	if make == nil {
		return nil, ErrMakeNotFound
	}

	return s.repo.GetModelsByMake(ctx, makeID)
}

func (s *vehicleService) GetModelByID(ctx context.Context, id int) (*vehiclemodels.Model, error) {
	if id <= 0 {
		return nil, ErrInvalidModelID
	}

	model, err := s.repo.GetModelByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, ErrModelNotFound
	}

	return model, nil
}

func (s *vehicleService) CreateModel(ctx context.Context, model *vehiclemodels.Model) (int, error) {
	// Validate required fields
	if model.ModelName == "" {
		return 0, errors.New("model name is required")
	}

	// Verify make exists
	make, err := s.repo.GetMakeByID(ctx, model.MakeID)
	if err != nil {
		return 0, err
	}
	if make == nil {
		return 0, ErrMakeNotFound
	}

	return s.repo.CreateModel(ctx, model)
}

func (s *vehicleService) UpdateModel(ctx context.Context, model *vehiclemodels.Model) error {
	if model.ModelID <= 0 {
		return ErrInvalidModelID
	}

	if model.ModelName == "" {
		return errors.New("model name is required")
	}

	// Check if model exists
	existing, err := s.repo.GetModelByID(ctx, model.ModelID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrModelNotFound
	}

	// Verify make exists if make ID is being changed
	if existing.MakeID != model.MakeID {
		make, err := s.repo.GetMakeByID(ctx, model.MakeID)
		if err != nil {
			return err
		}
		if make == nil {
			return ErrMakeNotFound
		}
	}

	return s.repo.UpdateModel(ctx, model)
}

func (s *vehicleService) DeleteModel(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrInvalidModelID
	}

	// Check if model exists
	existing, err := s.repo.GetModelByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrModelNotFound
	}

	// Check if model has any submodels
	submodels, err := s.repo.GetSubmodelsByModel(ctx, id)
	if err != nil {
		return err
	}
	if len(submodels) > 0 {
		return errors.New("cannot delete model with existing submodels")
	}

	return s.repo.DeleteModel(ctx, id)
}

// Submodel operations
func (s *vehicleService) GetAllSubmodels(ctx context.Context) ([]*vehiclemodels.Submodel, error) {
	return s.repo.GetAllSubmodels(ctx)
}

func (s *vehicleService) GetSubmodelsByModel(ctx context.Context, modelID int) ([]*vehiclemodels.Submodel, error) {
	if modelID <= 0 {
		return nil, ErrInvalidModelID
	}

	// Verify model exists
	model, err := s.repo.GetModelByID(ctx, modelID)
	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, ErrModelNotFound
	}

	return s.repo.GetSubmodelsByModel(ctx, modelID)
}

func (s *vehicleService) GetSubmodelByID(ctx context.Context, id int) (*vehiclemodels.Submodel, error) {
	if id <= 0 {
		return nil, errors.New("invalid submodel ID")
	}

	submodel, err := s.repo.GetSubmodelByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if submodel == nil {
		return nil, ErrSubmodelNotFound
	}

	return submodel, nil
}

func (s *vehicleService) CreateSubmodel(ctx context.Context, submodel *vehiclemodels.Submodel) (int, error) {
	// Validate required fields
	if err := s.validateSubmodel(submodel); err != nil {
		return 0, err
	}

	// Verify model exists
	model, err := s.repo.GetModelByID(ctx, submodel.ModelID)
	if err != nil {
		return 0, err
	}
	if model == nil {
		return 0, ErrModelNotFound
	}

	// Validate year range
	if submodel.YearTo != nil && *submodel.YearTo < submodel.YearFrom {
		return 0, errors.New("end year cannot be earlier than start year")
	}

	return s.repo.CreateSubmodel(ctx, submodel)
}

func (s *vehicleService) UpdateSubmodel(ctx context.Context, submodel *vehiclemodels.Submodel) error {
	if submodel.SubmodelID <= 0 {
		return errors.New("invalid submodel ID")
	}

	// Validate required fields
	if err := s.validateSubmodel(submodel); err != nil {
		return err
	}

	// Check if submodel exists
	existing, err := s.repo.GetSubmodelByID(ctx, submodel.SubmodelID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrSubmodelNotFound
	}

	// Verify model exists if model ID is being changed
	if existing.ModelID != submodel.ModelID {
		model, err := s.repo.GetModelByID(ctx, submodel.ModelID)
		if err != nil {
			return err
		}
		if model == nil {
			return ErrModelNotFound
		}
	}

	// Validate year range
	if submodel.YearTo != nil && *submodel.YearTo < submodel.YearFrom {
		return errors.New("end year cannot be earlier than start year")
	}

	return s.repo.UpdateSubmodel(ctx, submodel)
}

func (s *vehicleService) DeleteSubmodel(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid submodel ID")
	}

	// Check if submodel exists
	existing, err := s.repo.GetSubmodelByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrSubmodelNotFound
	}

	// You might want to check for dependencies (e.g., parts compatibility) before deletion
	// This would require additional repository methods

	return s.repo.DeleteSubmodel(ctx, id)
}

// Helper function for submodel validation
func (s *vehicleService) validateSubmodel(submodel *vehiclemodels.Submodel) error {
	if submodel.SubmodelName == "" {
		return errors.New("submodel name is required")
	}
	if submodel.YearFrom <= 0 {
		return errors.New("valid start year is required")
	}
	if submodel.EngineType == "" {
		return errors.New("engine type is required")
	}
	if submodel.EngineDisplacement <= 0 {
		return errors.New("valid engine displacement is required")
	}
	if submodel.FuelType == "" {
		return errors.New("fuel type is required")
	}
	if submodel.TransmissionType == "" {
		return errors.New("transmission type is required")
	}
	if submodel.BodyType == "" {
		return errors.New("body type is required")
	}
	return nil
}
