package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	salesmodels "github.com/hsrvms/autoparts/internal/modules/sales/models"
	"github.com/hsrvms/autoparts/pkg/db"
	"github.com/jackc/pgx/v5"
)

type PostgresSaleRepository struct {
    db *db.Database
}

func NewPostgresSaleRepository(database *db.Database) SaleRepository {
    return &PostgresSaleRepository{
        db: database,
    }
}

func (r *PostgresSaleRepository) GetAll(ctx context.Context, filter *salesmodels.SaleFilter) ([]*salesmodels.Sale, error) {
    query := `
        SELECT
            s.sale_id, s.date, s.item_id, s.quantity,
            s.price_per_unit, s.total_price, s.transaction_number,
            s.customer_name, s.customer_phone, s.customer_email,
            s.sold_by, s.notes, s.created_at, s.updated_at,
            i.part_number as item_part_number,
            i.description as item_description,
            c.category_name
        FROM sales s
        JOIN items i ON s.item_id = i.item_id
        LEFT JOIN categories c ON i.category_id = c.category_id
        WHERE 1=1
    `

    var conditions []string
    var params []interface{}
    paramCount := 1

    if filter != nil {
        if filter.ItemID != nil {
            conditions = append(conditions, fmt.Sprintf("s.item_id = $%d", paramCount))
            params = append(params, *filter.ItemID)
            paramCount++
        }

        if filter.StartDate != nil {
            conditions = append(conditions, fmt.Sprintf("s.date >= $%d", paramCount))
            params = append(params, *filter.StartDate)
            paramCount++
        }

        if filter.EndDate != nil {
            conditions = append(conditions, fmt.Sprintf("s.date <= $%d", paramCount))
            params = append(params, *filter.EndDate)
            paramCount++
        }

        if filter.CustomerName != nil {
            conditions = append(conditions, fmt.Sprintf("s.customer_name ILIKE $%d", paramCount))
            params = append(params, "%"+*filter.CustomerName+"%")
            paramCount++
        }

        if filter.CustomerPhone != nil {
            conditions = append(conditions, fmt.Sprintf("s.customer_phone = $%d", paramCount))
            params = append(params, *filter.CustomerPhone)
            paramCount++
        }

        if filter.CustomerEmail != nil {
            conditions = append(conditions, fmt.Sprintf("s.customer_email = $%d", paramCount))
            params = append(params, *filter.CustomerEmail)
            paramCount++
        }

        if filter.TransactionNumber != nil {
            conditions = append(conditions, fmt.Sprintf("s.transaction_number = $%d", paramCount))
            params = append(params, *filter.TransactionNumber)
            paramCount++
        }

        if filter.SoldBy != nil {
            conditions = append(conditions, fmt.Sprintf("s.sold_by = $%d", paramCount))
            params = append(params, *filter.SoldBy)
            paramCount++
        }
    }

    if len(conditions) > 0 {
        query += " AND " + strings.Join(conditions, " AND ")
    }

    query += " ORDER BY s.date DESC"

    rows, err := r.db.Pool.Query(ctx, query, params...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var sales []*salesmodels.Sale
    for rows.Next() {
        sale := &salesmodels.Sale{}
        err := rows.Scan(
            &sale.SaleID,
            &sale.Date,
            &sale.ItemID,
            &sale.Quantity,
            &sale.PricePerUnit,
            &sale.TotalPrice,
            &sale.TransactionNumber,
            &sale.CustomerName,
            &sale.CustomerPhone,
            &sale.CustomerEmail,
            &sale.SoldBy,
            &sale.Notes,
            &sale.CreatedAt,
            &sale.UpdatedAt,
            &sale.ItemPartNumber,
            &sale.ItemDescription,
            &sale.CategoryName,
        )
        if err != nil {
            return nil, err
        }
        sales = append(sales, sale)
    }

    return sales, rows.Err()
}

func (r *PostgresSaleRepository) GetByID(ctx context.Context, id int) (*salesmodels.Sale, error) {
    query := `
        SELECT
            s.sale_id, s.date, s.item_id, s.quantity,
            s.price_per_unit, s.total_price, s.transaction_number,
            s.customer_name, s.customer_phone, s.customer_email,
            s.sold_by, s.notes, s.created_at, s.updated_at,
            i.part_number as item_part_number,
            i.description as item_description,
            c.category_name
        FROM sales s
        JOIN items i ON s.item_id = i.item_id
        LEFT JOIN categories c ON i.category_id = c.category_id
        WHERE s.sale_id = $1
    `

    sale := &salesmodels.Sale{}
    err := r.db.Pool.QueryRow(ctx, query, id).Scan(
        &sale.SaleID,
        &sale.Date,
        &sale.ItemID,
        &sale.Quantity,
        &sale.PricePerUnit,
        &sale.TotalPrice,
        &sale.TransactionNumber,
        &sale.CustomerName,
        &sale.CustomerPhone,
        &sale.CustomerEmail,
        &sale.SoldBy,
        &sale.Notes,
        &sale.CreatedAt,
        &sale.UpdatedAt,
        &sale.ItemPartNumber,
        &sale.ItemDescription,
        &sale.CategoryName,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }

    return sale, nil
}

func (r *PostgresSaleRepository) Create(ctx context.Context, sale *salesmodels.Sale) (int, error) {
    tx, err := r.db.Pool.Begin(ctx)
    if err != nil {
        return 0, err
    }
    defer tx.Rollback(ctx)

    // Insert the sale
    query := `
        INSERT INTO sales (
            date, item_id, quantity, price_per_unit,
            total_price, transaction_number, customer_name,
            customer_phone, customer_email, sold_by, notes
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING sale_id
    `

    var id int
    err = tx.QueryRow(
        ctx, query,
        sale.Date,
        sale.ItemID,
        sale.Quantity,
        sale.PricePerUnit,
        sale.TotalPrice,
        sale.TransactionNumber,
        sale.CustomerName,
        sale.CustomerPhone,
        sale.CustomerEmail,
        sale.SoldBy,
        sale.Notes,
    ).Scan(&id)

    if err != nil {
        return 0, err
    }

    // Commit the transaction
    if err = tx.Commit(ctx); err != nil {
        return 0, err
    }

    return id, nil
}

func (r *PostgresSaleRepository) Update(ctx context.Context, sale *salesmodels.Sale) error {
    query := `
        UPDATE sales SET
            date = $2,
            item_id = $3,
            quantity = $4,
            price_per_unit = $5,
            total_price = $6,
            transaction_number = $7,
            customer_name = $8,
            customer_phone = $9,
            customer_email = $10,
            sold_by = $11,
            notes = $12
        WHERE sale_id = $1
    `

    result, err := r.db.Pool.Exec(
        ctx, query,
        sale.SaleID,
        sale.Date,
        sale.ItemID,
        sale.Quantity,
        sale.PricePerUnit,
        sale.TotalPrice,
        sale.TransactionNumber,
        sale.CustomerName,
        sale.CustomerPhone,
        sale.CustomerEmail,
        sale.SoldBy,
        sale.Notes,
    )

    if err != nil {
        return err
    }

    if result.RowsAffected() == 0 {
        return errors.New("sale not found")
    }

    return nil
}

func (r *PostgresSaleRepository) Delete(ctx context.Context, id int) error {
    query := `DELETE FROM sales WHERE sale_id = $1`

    result, err := r.db.Pool.Exec(ctx, query, id)
    if err != nil {
        return err
    }

    if result.RowsAffected() == 0 {
        return errors.New("sale not found")
    }

    return nil
}

func (r *PostgresSaleRepository) GetByTransactionNumber(ctx context.Context, transactionNumber string) (*salesmodels.Sale, error) {
    query := `
        SELECT
            s.sale_id, s.date, s.item_id, s.quantity,
            s.price_per_unit, s.total_price, s.transaction_number,
            s.customer_name, s.customer_phone, s.customer_email,
            s.sold_by, s.notes, s.created_at, s.updated_at,
            i.part_number as item_part_number,
            i.description as item_description,
            c.category_name
        FROM sales s
        JOIN items i ON s.item_id = i.item_id
        LEFT JOIN categories c ON i.category_id = c.category_id
        WHERE s.transaction_number = $1
    `

    sale := &salesmodels.Sale{}
    err := r.db.Pool.QueryRow(ctx, query, transactionNumber).Scan(
        &sale.SaleID,
        &sale.Date,
        &sale.ItemID,
        &sale.Quantity,
        &sale.PricePerUnit,
        &sale.TotalPrice,
        &sale.TransactionNumber,
        &sale.CustomerName,
        &sale.CustomerPhone,
        &sale.CustomerEmail,
        &sale.SoldBy,
        &sale.Notes,
        &sale.CreatedAt,
        &sale.UpdatedAt,
        &sale.ItemPartNumber,
        &sale.ItemDescription,
        &sale.CategoryName,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }

    return sale, nil
}

func (r *PostgresSaleRepository) GetItemSales(ctx context.Context, itemID int) ([]*salesmodels.Sale, error) {
    filter := &salesmodels.SaleFilter{
        ItemID: &itemID,
    }
    return r.GetAll(ctx, filter)
}

func (r *PostgresSaleRepository) GetCustomerSales(ctx context.Context, customerEmail string) ([]*salesmodels.Sale, error) {
    filter := &salesmodels.SaleFilter{
        CustomerEmail: &customerEmail,
    }
    return r.GetAll(ctx, filter)
}
