package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	purchasemodels "github.com/hsrvms/autoparts/internal/modules/purchases/models"
	"github.com/hsrvms/autoparts/pkg/db"
	"github.com/jackc/pgx/v5"
)

type PostgresPurchaseRepository struct {
    db *db.Database
}

func NewPostgresPurchaseRepository(database *db.Database) PurchaseRepository {
    return &PostgresPurchaseRepository{
        db: database,
    }
}

func (r *PostgresPurchaseRepository) GetAll(ctx context.Context, filter *purchasemodels.PurchaseFilter) ([]*purchasemodels.Purchase, error) {
    query := `
        SELECT
            p.purchase_id, p.date, p.supplier_id, p.item_id,
            p.quantity, p.cost_per_unit, p.total_cost,
            p.invoice_number, p.received_by, p.notes,
            p.created_at, p.updated_at,
            s.name as supplier_name,
            i.part_number as item_part_number,
            i.description as item_description
        FROM purchases p
        JOIN suppliers s ON p.supplier_id = s.supplier_id
        JOIN items i ON p.item_id = i.item_id
        WHERE 1=1
    `

    var conditions []string
    var params []interface{}
    paramCount := 1

    if filter != nil {
        if filter.SupplierID != nil {
            conditions = append(conditions, fmt.Sprintf("p.supplier_id = $%d", paramCount))
            params = append(params, *filter.SupplierID)
            paramCount++
        }

        if filter.ItemID != nil {
            conditions = append(conditions, fmt.Sprintf("p.item_id = $%d", paramCount))
            params = append(params, *filter.ItemID)
            paramCount++
        }

        if filter.StartDate != nil {
            conditions = append(conditions, fmt.Sprintf("p.date >= $%d", paramCount))
            params = append(params, *filter.StartDate)
            paramCount++
        }

        if filter.EndDate != nil {
            conditions = append(conditions, fmt.Sprintf("p.date <= $%d", paramCount))
            params = append(params, *filter.EndDate)
            paramCount++
        }

        if filter.InvoiceNumber != nil {
            conditions = append(conditions, fmt.Sprintf("p.invoice_number ILIKE $%d", paramCount))
            params = append(params, "%"+*filter.InvoiceNumber+"%")
            paramCount++
        }
    }

    if len(conditions) > 0 {
        query += " AND " + strings.Join(conditions, " AND ")
    }

    query += " ORDER BY p.date DESC"

    rows, err := r.db.Pool.Query(ctx, query, params...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var purchases []*purchasemodels.Purchase
    for rows.Next() {
        purchase := &purchasemodels.Purchase{}
        err := rows.Scan(
            &purchase.PurchaseID,
            &purchase.Date,
            &purchase.SupplierID,
            &purchase.ItemID,
            &purchase.Quantity,
            &purchase.CostPerUnit,
            &purchase.TotalCost,
            &purchase.InvoiceNumber,
            &purchase.ReceivedBy,
            &purchase.Notes,
            &purchase.CreatedAt,
            &purchase.UpdatedAt,
            &purchase.SupplierName,
            &purchase.ItemPartNumber,
            &purchase.ItemDescription,
        )
        if err != nil {
            return nil, err
        }
        purchases = append(purchases, purchase)
    }

    return purchases, rows.Err()
}

func (r *PostgresPurchaseRepository) GetByID(ctx context.Context, id int) (*purchasemodels.Purchase, error) {
    query := `
        SELECT
            p.purchase_id, p.date, p.supplier_id, p.item_id,
            p.quantity, p.cost_per_unit, p.total_cost,
            p.invoice_number, p.received_by, p.notes,
            p.created_at, p.updated_at,
            s.name as supplier_name,
            i.part_number as item_part_number,
            i.description as item_description
        FROM purchases p
        JOIN suppliers s ON p.supplier_id = s.supplier_id
        JOIN items i ON p.item_id = i.item_id
        WHERE p.purchase_id = $1
    `

    purchase := &purchasemodels.Purchase{}
    err := r.db.Pool.QueryRow(ctx, query, id).Scan(
        &purchase.PurchaseID,
        &purchase.Date,
        &purchase.SupplierID,
        &purchase.ItemID,
        &purchase.Quantity,
        &purchase.CostPerUnit,
        &purchase.TotalCost,
        &purchase.InvoiceNumber,
        &purchase.ReceivedBy,
        &purchase.Notes,
        &purchase.CreatedAt,
        &purchase.UpdatedAt,
        &purchase.SupplierName,
        &purchase.ItemPartNumber,
        &purchase.ItemDescription,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }

    return purchase, nil
}

func (r *PostgresPurchaseRepository) Create(ctx context.Context, purchase *purchasemodels.Purchase) (int, error) {
    query := `
        INSERT INTO purchases (
            date, supplier_id, item_id, quantity,
            cost_per_unit, total_cost, invoice_number,
            received_by, notes
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING purchase_id
    `

    var id int
    err := r.db.Pool.QueryRow(
        ctx, query,
        purchase.Date,
        purchase.SupplierID,
        purchase.ItemID,
        purchase.Quantity,
        purchase.CostPerUnit,
        purchase.TotalCost,
        purchase.InvoiceNumber,
        purchase.ReceivedBy,
        purchase.Notes,
    ).Scan(&id)

    if err != nil {
        return 0, err
    }

    return id, nil
}

func (r *PostgresPurchaseRepository) Update(ctx context.Context, purchase *purchasemodels.Purchase) error {
    query := `
        UPDATE purchases SET
            date = $2,
            supplier_id = $3,
            item_id = $4,
            quantity = $5,
            cost_per_unit = $6,
            total_cost = $7,
            invoice_number = $8,
            received_by = $9,
            notes = $10
        WHERE purchase_id = $1
    `

    result, err := r.db.Pool.Exec(
        ctx, query,
        purchase.PurchaseID,
        purchase.Date,
        purchase.SupplierID,
        purchase.ItemID,
        purchase.Quantity,
        purchase.CostPerUnit,
        purchase.TotalCost,
        purchase.InvoiceNumber,
        purchase.ReceivedBy,
        purchase.Notes,
    )

    if err != nil {
        return err
    }

    if result.RowsAffected() == 0 {
        return errors.New("purchase not found")
    }

    return nil
}

func (r *PostgresPurchaseRepository) Delete(ctx context.Context, id int) error {
    query := `DELETE FROM purchases WHERE purchase_id = $1`

    result, err := r.db.Pool.Exec(ctx, query, id)
    if err != nil {
        return err
    }

    if result.RowsAffected() == 0 {
        return errors.New("purchase not found")
    }

    return nil
}

func (r *PostgresPurchaseRepository) GetByInvoiceNumber(ctx context.Context, invoiceNumber string) (*purchasemodels.Purchase, error) {
    query := `
        SELECT
            p.purchase_id, p.date, p.supplier_id, p.item_id,
            p.quantity, p.cost_per_unit, p.total_cost,
            p.invoice_number, p.received_by, p.notes,
            p.created_at, p.updated_at,
            s.name as supplier_name,
            i.part_number as item_part_number,
            i.description as item_description
        FROM purchases p
        JOIN suppliers s ON p.supplier_id = s.supplier_id
        JOIN items i ON p.item_id = i.item_id
        WHERE p.invoice_number = $1
    `

    purchase := &purchasemodels.Purchase{}
    err := r.db.Pool.QueryRow(ctx, query, invoiceNumber).Scan(
        &purchase.PurchaseID,
        &purchase.Date,
        &purchase.SupplierID,
        &purchase.ItemID,
        &purchase.Quantity,
        &purchase.CostPerUnit,
        &purchase.TotalCost,
        &purchase.InvoiceNumber,
        &purchase.ReceivedBy,
        &purchase.Notes,
        &purchase.CreatedAt,
        &purchase.UpdatedAt,
        &purchase.SupplierName,
        &purchase.ItemPartNumber,
        &purchase.ItemDescription,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }

    return purchase, nil
}

func (r *PostgresPurchaseRepository) GetSupplierPurchases(ctx context.Context, supplierID int) ([]*purchasemodels.Purchase, error) {
    filter := &purchasemodels.PurchaseFilter{
        SupplierID: &supplierID,
    }
    return r.GetAll(ctx, filter)
}

func (r *PostgresPurchaseRepository) GetItemPurchases(ctx context.Context, itemID int) ([]*purchasemodels.Purchase, error) {
    filter := &purchasemodels.PurchaseFilter{
        ItemID: &itemID,
    }
    return r.GetAll(ctx, filter)
}
