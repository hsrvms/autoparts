// internal/modules/vehicles/repositories/postgres_repository.go
package repositories

import (
	"context"
	"errors"

	vehiclemodels "github.com/hsrvms/autoparts/internal/modules/vehicles/models"
	"github.com/hsrvms/autoparts/pkg/db"
	"github.com/jackc/pgx/v5"
)

type PostgresVehicleRepository struct {
	db *db.Database
}

func NewPostgresVehicleRepository(database *db.Database) VehicleRepository {
	return &PostgresVehicleRepository{
		db: database,
	}
}

// Make operations
func (r *PostgresVehicleRepository) GetAllMakes(ctx context.Context) ([]*vehiclemodels.Make, error) {
	query := `
		SELECT make_id, make_name, country, created_at, updated_at
		FROM vehicle_makes
		ORDER BY make_name
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var makes []*vehiclemodels.Make
	for rows.Next() {
		make := &vehiclemodels.Make{}
		err := rows.Scan(
			&make.MakeID,
			&make.MakeName,
			&make.Country,
			&make.CreatedAt,
			&make.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		makes = append(makes, make)
	}

	return makes, rows.Err()
}

func (r *PostgresVehicleRepository) GetMakeByID(ctx context.Context, id int) (*vehiclemodels.Make, error) {
	query := `
		SELECT make_id, make_name, country, created_at, updated_at
		FROM vehicle_makes
		WHERE make_id = $1
	`

	make := &vehiclemodels.Make{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&make.MakeID,
		&make.MakeName,
		&make.Country,
		&make.CreatedAt,
		&make.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return make, nil
}

func (r *PostgresVehicleRepository) CreateMake(ctx context.Context, make *vehiclemodels.Make) (int, error) {
	query := `
		INSERT INTO vehicle_makes (make_name, country)
		VALUES ($1, $2)
		RETURNING make_id
	`

	var id int
	err := r.db.Pool.QueryRow(ctx, query, make.MakeName, make.Country).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PostgresVehicleRepository) UpdateMake(ctx context.Context, make *vehiclemodels.Make) error {
	query := `
		UPDATE vehicle_makes
		SET make_name = $2, country = $3
		WHERE make_id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query, make.MakeID, make.MakeName, make.Country)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("make not found")
	}

	return nil
}

func (r *PostgresVehicleRepository) DeleteMake(ctx context.Context, id int) error {
	query := `DELETE FROM vehicle_makes WHERE make_id = $1`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("make not found")
	}

	return nil
}

// Model operations
func (r *PostgresVehicleRepository) GetAllModels(ctx context.Context) ([]*vehiclemodels.Model, error) {
	query := `
		SELECT m.model_id, m.make_id, m.model_name, m.created_at, m.updated_at,
			   mk.make_name
		FROM vehicle_models m
		JOIN vehicle_makes mk ON m.make_id = mk.make_id
		ORDER BY mk.make_name, m.model_name
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []*vehiclemodels.Model
	for rows.Next() {
		model := &vehiclemodels.Model{}
		err := rows.Scan(
			&model.ModelID,
			&model.MakeID,
			&model.ModelName,
			&model.CreatedAt,
			&model.UpdatedAt,
			&model.MakeName,
		)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}

	return models, rows.Err()
}

func (r *PostgresVehicleRepository) GetModelsByMake(ctx context.Context, makeID int) ([]*vehiclemodels.Model, error) {
	query := `
		SELECT m.model_id, m.make_id, m.model_name, m.created_at, m.updated_at,
			   mk.make_name
		FROM vehicle_models m
		JOIN vehicle_makes mk ON m.make_id = mk.make_id
		WHERE m.make_id = $1
		ORDER BY m.model_name
	`

	rows, err := r.db.Pool.Query(ctx, query, makeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []*vehiclemodels.Model
	for rows.Next() {
		model := &vehiclemodels.Model{}
		err := rows.Scan(
			&model.ModelID,
			&model.MakeID,
			&model.ModelName,
			&model.CreatedAt,
			&model.UpdatedAt,
			&model.MakeName,
		)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}

	return models, rows.Err()
}

func (r *PostgresVehicleRepository) GetModelByID(ctx context.Context, id int) (*vehiclemodels.Model, error) {
	query := `
		SELECT m.model_id, m.make_id, m.model_name, m.created_at, m.updated_at,
			   mk.make_name
		FROM vehicle_models m
		JOIN vehicle_makes mk ON m.make_id = mk.make_id
		WHERE m.model_id = $1
	`

	model := &vehiclemodels.Model{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&model.ModelID,
		&model.MakeID,
		&model.ModelName,
		&model.CreatedAt,
		&model.UpdatedAt,
		&model.MakeName,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return model, nil
}

func (r *PostgresVehicleRepository) CreateModel(ctx context.Context, model *vehiclemodels.Model) (int, error) {
	query := `
		INSERT INTO vehicle_models (make_id, model_name)
		VALUES ($1, $2)
		RETURNING model_id
	`

	var id int
	err := r.db.Pool.QueryRow(ctx, query, model.MakeID, model.ModelName).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PostgresVehicleRepository) UpdateModel(ctx context.Context, model *vehiclemodels.Model) error {
	query := `
		UPDATE vehicle_models
		SET make_id = $2, model_name = $3
		WHERE model_id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query, model.ModelID, model.MakeID, model.ModelName)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("model not found")
	}

	return nil
}

func (r *PostgresVehicleRepository) DeleteModel(ctx context.Context, id int) error {
	query := `DELETE FROM vehicle_models WHERE model_id = $1`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("model not found")
	}

	return nil
}

// Submodel operations
func (r *PostgresVehicleRepository) GetAllSubmodels(ctx context.Context) ([]*vehiclemodels.Submodel, error) {
	query := `
		SELECT s.submodel_id, s.model_id, s.submodel_name, s.year_from, s.year_to,
			   s.engine_type, s.engine_displacement, s.fuel_type, s.transmission_type,
			   s.body_type, s.created_at, s.updated_at,
			   m.model_name, mk.make_name
		FROM vehicle_submodels s
		JOIN vehicle_models m ON s.model_id = m.model_id
		JOIN vehicle_makes mk ON m.make_id = mk.make_id
		ORDER BY mk.make_name, m.model_name, s.submodel_name
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submodels []*vehiclemodels.Submodel
	for rows.Next() {
		submodel := &vehiclemodels.Submodel{}
		err := rows.Scan(
			&submodel.SubmodelID,
			&submodel.ModelID,
			&submodel.SubmodelName,
			&submodel.YearFrom,
			&submodel.YearTo,
			&submodel.EngineType,
			&submodel.EngineDisplacement,
			&submodel.FuelType,
			&submodel.TransmissionType,
			&submodel.BodyType,
			&submodel.CreatedAt,
			&submodel.UpdatedAt,
			&submodel.ModelName,
			&submodel.MakeName,
		)
		if err != nil {
			return nil, err
		}
		submodels = append(submodels, submodel)
	}

	return submodels, rows.Err()
}

func (r *PostgresVehicleRepository) GetSubmodelsByModel(ctx context.Context, modelID int) ([]*vehiclemodels.Submodel, error) {
	query := `
		SELECT s.submodel_id, s.model_id, s.submodel_name, s.year_from, s.year_to,
			   s.engine_type, s.engine_displacement, s.fuel_type, s.transmission_type,
			   s.body_type, s.created_at, s.updated_at,
			   m.model_name, mk.make_name
		FROM vehicle_submodels s
		JOIN vehicle_models m ON s.model_id = m.model_id
		JOIN vehicle_makes mk ON m.make_id = mk.make_id
		WHERE s.model_id = $1
		ORDER BY s.year_from DESC, s.submodel_name
	`

	rows, err := r.db.Pool.Query(ctx, query, modelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submodels []*vehiclemodels.Submodel
	for rows.Next() {
		submodel := &vehiclemodels.Submodel{}
		err := rows.Scan(
			&submodel.SubmodelID,
			&submodel.ModelID,
			&submodel.SubmodelName,
			&submodel.YearFrom,
			&submodel.YearTo,
			&submodel.EngineType,
			&submodel.EngineDisplacement,
			&submodel.FuelType,
			&submodel.TransmissionType,
			&submodel.BodyType,
			&submodel.CreatedAt,
			&submodel.UpdatedAt,
			&submodel.ModelName,
			&submodel.MakeName,
		)
		if err != nil {
			return nil, err
		}
		submodels = append(submodels, submodel)
	}

	return submodels, rows.Err()
}

func (r *PostgresVehicleRepository) GetSubmodelByID(ctx context.Context, id int) (*vehiclemodels.Submodel, error) {
	query := `
		SELECT s.submodel_id, s.model_id, s.submodel_name, s.year_from, s.year_to,
			   s.engine_type, s.engine_displacement, s.fuel_type, s.transmission_type,
			   s.body_type, s.created_at, s.updated_at,
			   m.model_name, mk.make_name
		FROM vehicle_submodels s
		JOIN vehicle_models m ON s.model_id = m.model_id
		JOIN vehicle_makes mk ON m.make_id = mk.make_id
		WHERE s.submodel_id = $1
	`

	submodel := &vehiclemodels.Submodel{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&submodel.SubmodelID,
		&submodel.ModelID,
		&submodel.SubmodelName,
		&submodel.YearFrom,
		&submodel.YearTo,
		&submodel.EngineType,
		&submodel.EngineDisplacement,
		&submodel.FuelType,
		&submodel.TransmissionType,
		&submodel.BodyType,
		&submodel.CreatedAt,
		&submodel.UpdatedAt,
		&submodel.ModelName,
		&submodel.MakeName,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return submodel, nil
}

func (r *PostgresVehicleRepository) CreateSubmodel(ctx context.Context, submodel *vehiclemodels.Submodel) (int, error) {
	query := `
		INSERT INTO vehicle_submodels (
			model_id, submodel_name, year_from, year_to,
			engine_type, engine_displacement, fuel_type,
			transmission_type, body_type
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING submodel_id
	`

	var id int
	err := r.db.Pool.QueryRow(
		ctx,
		query,
		submodel.ModelID,
		submodel.SubmodelName,
		submodel.YearFrom,
		submodel.YearTo,
		submodel.EngineType,
		submodel.EngineDisplacement,
		submodel.FuelType,
		submodel.TransmissionType,
		submodel.BodyType,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PostgresVehicleRepository) UpdateSubmodel(ctx context.Context, submodel *vehiclemodels.Submodel) error {
	query := `
		UPDATE vehicle_submodels
		SET model_id = $2, submodel_name = $3, year_from = $4, year_to = $5,
			engine_type = $6, engine_displacement = $7, fuel_type = $8,
			transmission_type = $9, body_type = $10
		WHERE submodel_id = $1
	`

	result, err := r.db.Pool.Exec(
		ctx,
		query,
		submodel.SubmodelID,
		submodel.ModelID,
		submodel.SubmodelName,
		submodel.YearFrom,
		submodel.YearTo,
		submodel.EngineType,
		submodel.EngineDisplacement,
		submodel.FuelType,
		submodel.TransmissionType,
		submodel.BodyType,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("submodel not found")
	}

	return nil
}

func (r *PostgresVehicleRepository) DeleteSubmodel(ctx context.Context, id int) error {
	query := `DELETE FROM vehicle_submodels WHERE submodel_id = $1`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("submodel not found")
	}

	return nil
}
