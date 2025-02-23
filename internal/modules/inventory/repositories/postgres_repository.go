package repositories

import (
	"context"
	"errors"
	"fmt"

	inventorymodels "github.com/hsrvms/autoparts/internal/modules/inventory/models"
	"github.com/hsrvms/autoparts/pkg/db"
	"github.com/jackc/pgx/v5"
)

type PostgresInventoryRepository struct {
	db *db.Database
}

func NewPostgresInventoryRepository(database *db.Database) InventoryRepository {
	return &PostgresInventoryRepository{
		db: database,
	}
}

func (r *PostgresInventoryRepository) GetItems(ctx context.Context, filter *inventorymodels.ItemFilter) ([]*inventorymodels.Item, error) {
	query := `
		SELECT
			i.item_id, i.part_number, i.description, i.category_id, i.buy_price,
			i.sell_price, i.current_stock, i.minimum_stock, i.barcode, i.supplier_id,
			i.location_aisle, i.location_shelf, i.location_bin, i.weight_kg,
			i.dimensions_cm, i.warranty_period, i.image_url, i.is_active, i.notes,
			i.created_at, i.updated_at,
			c.category_name, s.name as supplier_name
		FROM items i
		LEFT JOIN categories c ON i.category_id = c.category_id
		LEFT JOIN suppliers s ON i.supplier_id = s.supplier_id
		WHERE 1=1
	`

	params := []interface{}{}
	paramCount := 1

	// Build query based on filters
	if filter != nil {
		if filter.CategoryID != nil {
			query += fmt.Sprintf(" AND i.category_id = $%d", paramCount)
			params = append(params, *filter.CategoryID)
			paramCount++
		}

		if filter.SupplierID != nil {
			query += fmt.Sprintf(" AND i.supplier_id = $%d", paramCount)
			params = append(params, *filter.SupplierID)
			paramCount++
		}

		if filter.PartNumber != nil {
			query += fmt.Sprintf(" AND i.part_number ILIKE $%d", paramCount)
			params = append(params, "%"+*filter.PartNumber+"%")
			paramCount++
		}

		if filter.SearchTerm != nil {
			query += fmt.Sprintf(" AND (i.part_number ILIKE $%d OR i.description ILIKE $%d)", paramCount, paramCount)
			params = append(params, "%"+*filter.SearchTerm+"%")
			paramCount++
		}

		if filter.LowStock != nil && *filter.LowStock {
			query += " AND i.current_stock <= i.minimum_stock"
		}

		if filter.IsActive != nil {
			query += fmt.Sprintf(" AND i.is_active = $%d", paramCount)
			params = append(params, *filter.IsActive)
			paramCount++
		}
	}

	query += " ORDER BY i.part_number"

	rows, err := r.db.Pool.Query(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*inventorymodels.Item
	for rows.Next() {
		item := &inventorymodels.Item{}
		err := rows.Scan(
			&item.ItemID, &item.PartNumber, &item.Description, &item.CategoryID,
			&item.BuyPrice, &item.SellPrice, &item.CurrentStock, &item.MinimumStock,
			&item.Barcode, &item.SupplierID, &item.LocationAisle, &item.LocationShelf,
			&item.LocationBin, &item.WeightKg, &item.DimensionsCm, &item.WarrantyPeriod,
			&item.ImageURL, &item.IsActive, &item.Notes, &item.CreatedAt, &item.UpdatedAt,
			&item.CategoryName, &item.SupplierName,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *PostgresInventoryRepository) GetItemByID(ctx context.Context, id int) (*inventorymodels.Item, error) {
	query := `
		SELECT
			i.item_id, i.part_number, i.description, i.category_id, i.buy_price,
			i.sell_price, i.current_stock, i.minimum_stock, i.barcode, i.supplier_id,
			i.location_aisle, i.location_shelf, i.location_bin, i.weight_kg,
			i.dimensions_cm, i.warranty_period, i.image_url, i.is_active, i.notes,
			i.created_at, i.updated_at,
			c.category_name, s.name as supplier_name
		FROM items i
		LEFT JOIN categories c ON i.category_id = c.category_id
		LEFT JOIN suppliers s ON i.supplier_id = s.supplier_id
		WHERE i.item_id = $1
	`

	item := &inventorymodels.Item{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&item.ItemID, &item.PartNumber, &item.Description, &item.CategoryID,
		&item.BuyPrice, &item.SellPrice, &item.CurrentStock, &item.MinimumStock,
		&item.Barcode, &item.SupplierID, &item.LocationAisle, &item.LocationShelf,
		&item.LocationBin, &item.WeightKg, &item.DimensionsCm, &item.WarrantyPeriod,
		&item.ImageURL, &item.IsActive, &item.Notes, &item.CreatedAt, &item.UpdatedAt,
		&item.CategoryName, &item.SupplierName,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return item, nil
}

func (r *PostgresInventoryRepository) GetItemByPartNumber(ctx context.Context, partNumber string) (*inventorymodels.Item, error) {
	query := `
		SELECT
			i.item_id, i.part_number, i.description, i.category_id, i.buy_price,
			i.sell_price, i.current_stock, i.minimum_stock, i.barcode, i.supplier_id,
			i.location_aisle, i.location_shelf, i.location_bin, i.weight_kg,
			i.dimensions_cm, i.warranty_period, i.image_url, i.is_active, i.notes,
			i.created_at, i.updated_at,
			c.category_name, s.name as supplier_name
		FROM items i
		LEFT JOIN categories c ON i.category_id = c.category_id
		LEFT JOIN suppliers s ON i.supplier_id = s.supplier_id
		WHERE i.part_number = $1
	`

	item := &inventorymodels.Item{}
	err := r.db.Pool.QueryRow(ctx, query, partNumber).Scan(
		&item.ItemID, &item.PartNumber, &item.Description, &item.CategoryID,
		&item.BuyPrice, &item.SellPrice, &item.CurrentStock, &item.MinimumStock,
		&item.Barcode, &item.SupplierID, &item.LocationAisle, &item.LocationShelf,
		&item.LocationBin, &item.WeightKg, &item.DimensionsCm, &item.WarrantyPeriod,
		&item.ImageURL, &item.IsActive, &item.Notes, &item.CreatedAt, &item.UpdatedAt,
		&item.CategoryName, &item.SupplierName,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return item, nil
}

func (r *PostgresInventoryRepository) GetItemByBarcode(ctx context.Context, barcode string) (*inventorymodels.Item, error) {
	query := `
		SELECT
			i.item_id, i.part_number, i.description, i.category_id, i.buy_price,
			i.sell_price, i.current_stock, i.minimum_stock, i.barcode, i.supplier_id,
			i.location_aisle, i.location_shelf, i.location_bin, i.weight_kg,
			i.dimensions_cm, i.warranty_period, i.image_url, i.is_active, i.notes,
			i.created_at, i.updated_at,
			c.category_name, s.name as supplier_name
		FROM items i
		LEFT JOIN categories c ON i.category_id = c.category_id
		LEFT JOIN suppliers s ON i.supplier_id = s.supplier_id
		WHERE i.barcode = $1
	`

	item := &inventorymodels.Item{}
	err := r.db.Pool.QueryRow(ctx, query, barcode).Scan(
		&item.ItemID, &item.PartNumber, &item.Description, &item.CategoryID,
		&item.BuyPrice, &item.SellPrice, &item.CurrentStock, &item.MinimumStock,
		&item.Barcode, &item.SupplierID, &item.LocationAisle, &item.LocationShelf,
		&item.LocationBin, &item.WeightKg, &item.DimensionsCm, &item.WarrantyPeriod,
		&item.ImageURL, &item.IsActive, &item.Notes, &item.CreatedAt, &item.UpdatedAt,
		&item.CategoryName, &item.SupplierName,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return item, nil
}

func (r *PostgresInventoryRepository) CreateItem(ctx context.Context, item *inventorymodels.Item) (int, error) {
	query := `
		INSERT INTO items (
			part_number, description, category_id, buy_price, sell_price,
			current_stock, minimum_stock, barcode, supplier_id, location_aisle,
			location_shelf, location_bin, weight_kg, dimensions_cm,
			warranty_period, image_url, is_active, notes
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18
		)
		RETURNING item_id
	`

	var id int
	err := r.db.Pool.QueryRow(
		ctx, query,
		item.PartNumber, item.Description, item.CategoryID, item.BuyPrice,
		item.SellPrice, item.CurrentStock, item.MinimumStock, item.Barcode,
		item.SupplierID, item.LocationAisle, item.LocationShelf, item.LocationBin,
		item.WeightKg, item.DimensionsCm, item.WarrantyPeriod, item.ImageURL,
		item.IsActive, item.Notes,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PostgresInventoryRepository) UpdateItem(ctx context.Context, item *inventorymodels.Item) error {
	query := `
		UPDATE items SET
			part_number = $2, description = $3, category_id = $4,
			buy_price = $5, sell_price = $6, current_stock = $7,
			minimum_stock = $8, barcode = $9, supplier_id = $10,
			location_aisle = $11, location_shelf = $12, location_bin = $13,
			weight_kg = $14, dimensions_cm = $15, warranty_period = $16,
			image_url = $17, is_active = $18, notes = $19
		WHERE item_id = $1
	`

	result, err := r.db.Pool.Exec(
		ctx, query,
		item.ItemID, item.PartNumber, item.Description, item.CategoryID,
		item.BuyPrice, item.SellPrice, item.CurrentStock, item.MinimumStock,
		item.Barcode, item.SupplierID, item.LocationAisle, item.LocationShelf,
		item.LocationBin, item.WeightKg, item.DimensionsCm, item.WarrantyPeriod,
		item.ImageURL, item.IsActive, item.Notes,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("item not found")
	}

	return nil
}

func (r *PostgresInventoryRepository) DeleteItem(ctx context.Context, id int) error {
	query := `DELETE FROM items WHERE item_id = $1`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("item not found")
	}

	return nil
}

func (r *PostgresInventoryRepository) GetCompatibilities(ctx context.Context, itemID int) ([]*inventorymodels.Compatibility, error) {
	query := `
        SELECT
            c.compat_id, c.item_id, c.submodel_id, c.notes, c.created_at,
            m.model_name, mk.make_name, s.submodel_name
        FROM compatibility c
        JOIN vehicle_submodels s ON c.submodel_id = s.submodel_id
        JOIN vehicle_models m ON s.model_id = m.model_id
        JOIN vehicle_makes mk ON m.make_id = mk.make_id
        WHERE c.item_id = $1
        ORDER BY mk.make_name, m.model_name, s.submodel_name
    `

	rows, err := r.db.Pool.Query(ctx, query, itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var compatibilities []*inventorymodels.Compatibility
	for rows.Next() {
		compatibility := &inventorymodels.Compatibility{}
		err := rows.Scan(
			&compatibility.CompatID,
			&compatibility.ItemID,
			&compatibility.SubmodelID,
			&compatibility.Notes,
			&compatibility.CreatedAt,
			&compatibility.ModelName,
			&compatibility.MakeName,
			&compatibility.SubmodelName,
		)
		if err != nil {
			return nil, err
		}
		compatibilities = append(compatibilities, compatibility)
	}

	return compatibilities, rows.Err()
}

func (r *PostgresInventoryRepository) AddCompatibility(ctx context.Context, compatibility *inventorymodels.Compatibility) (int, error) {
	query := `
        INSERT INTO compatibility (item_id, submodel_id, notes)
        VALUES ($1, $2, $3)
        RETURNING compat_id
    `

	var id int
	err := r.db.Pool.QueryRow(
		ctx, query,
		compatibility.ItemID,
		compatibility.SubmodelID,
		compatibility.Notes,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PostgresInventoryRepository) RemoveCompatibility(ctx context.Context, itemID, submodelID int) error {
	query := `
        DELETE FROM compatibility
        WHERE item_id = $1 AND submodel_id = $2
    `

	result, err := r.db.Pool.Exec(ctx, query, itemID, submodelID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("compatibility not found")
	}

	return nil
}

func (r *PostgresInventoryRepository) GetCompatibleItems(ctx context.Context, submodelID int) ([]*inventorymodels.Item, error) {
	query := `
        SELECT
            i.item_id, i.part_number, i.description, i.category_id, i.buy_price,
            i.sell_price, i.current_stock, i.minimum_stock, i.barcode, i.supplier_id,
            i.location_aisle, i.location_shelf, i.location_bin, i.weight_kg,
            i.dimensions_cm, i.warranty_period, i.image_url, i.is_active, i.notes,
            i.created_at, i.updated_at,
            c.category_name, s.name as supplier_name
        FROM items i
        LEFT JOIN categories c ON i.category_id = c.category_id
        LEFT JOIN suppliers s ON i.supplier_id = s.supplier_id
        JOIN compatibility comp ON i.item_id = comp.item_id
        WHERE comp.submodel_id = $1 AND i.is_active = true
        ORDER BY i.part_number
    `

	rows, err := r.db.Pool.Query(ctx, query, submodelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*inventorymodels.Item
	for rows.Next() {
		item := &inventorymodels.Item{}
		err := rows.Scan(
			&item.ItemID, &item.PartNumber, &item.Description, &item.CategoryID,
			&item.BuyPrice, &item.SellPrice, &item.CurrentStock, &item.MinimumStock,
			&item.Barcode, &item.SupplierID, &item.LocationAisle, &item.LocationShelf,
			&item.LocationBin, &item.WeightKg, &item.DimensionsCm, &item.WarrantyPeriod,
			&item.ImageURL, &item.IsActive, &item.Notes, &item.CreatedAt, &item.UpdatedAt,
			&item.CategoryName, &item.SupplierName,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *PostgresInventoryRepository) GetLowStockItems(ctx context.Context) ([]*inventorymodels.Item, error) {
	query := `
        SELECT
            i.item_id, i.part_number, i.description, i.category_id, i.buy_price,
            i.sell_price, i.current_stock, i.minimum_stock, i.barcode, i.supplier_id,
            i.location_aisle, i.location_shelf, i.location_bin, i.weight_kg,
            i.dimensions_cm, i.warranty_period, i.image_url, i.is_active, i.notes,
            i.created_at, i.updated_at,
            c.category_name, s.name as supplier_name
        FROM items i
        LEFT JOIN categories c ON i.category_id = c.category_id
        LEFT JOIN suppliers s ON i.supplier_id = s.supplier_id
        WHERE i.current_stock <= i.minimum_stock AND i.is_active = true
        ORDER BY i.current_stock ASC, i.part_number
    `

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*inventorymodels.Item
	for rows.Next() {
		item := &inventorymodels.Item{}
		err := rows.Scan(
			&item.ItemID, &item.PartNumber, &item.Description, &item.CategoryID,
			&item.BuyPrice, &item.SellPrice, &item.CurrentStock, &item.MinimumStock,
			&item.Barcode, &item.SupplierID, &item.LocationAisle, &item.LocationShelf,
			&item.LocationBin, &item.WeightKg, &item.DimensionsCm, &item.WarrantyPeriod,
			&item.ImageURL, &item.IsActive, &item.Notes, &item.CreatedAt, &item.UpdatedAt,
			&item.CategoryName, &item.SupplierName,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}
