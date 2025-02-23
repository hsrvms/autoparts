// internal/modules/suppliers/repositories/postgres_repository.go
package repositories

import (
	"context"
	"errors"
	"fmt"

	suppliermodels "github.com/hsrvms/autoparts/internal/modules/suppliers/models"
	"github.com/hsrvms/autoparts/pkg/db"
	"github.com/jackc/pgx/v5"
)

type PostgresSupplierRepository struct {
    db *db.Database
}

func NewPostgresSupplierRepository(database *db.Database) SupplierRepository {
    return &PostgresSupplierRepository{
        db: database,
    }
}

func (r *PostgresSupplierRepository) GetAll(ctx context.Context, filter *suppliermodels.SupplierFilter) ([]*suppliermodels.Supplier, error) {
    query := `
        SELECT DISTINCT s.supplier_id, s.name, s.contact_person, s.phone, s.email,
               s.address, s.tax_id, s.payment_terms, s.notes, s.created_at, s.updated_at
        FROM suppliers s
    `
    params := []interface{}{}
    paramCount := 1

    if filter != nil {
        if filter.SearchTerm != nil {
            query += fmt.Sprintf(" WHERE (s.name ILIKE $%d OR s.contact_person ILIKE $%d OR s.email ILIKE $%d)",
                paramCount, paramCount, paramCount)
            params = append(params, "%"+*filter.SearchTerm+"%")
            paramCount++
        }

        if filter.HasActiveItems != nil && *filter.HasActiveItems {
            if paramCount == 1 {
                query += " WHERE"
            } else {
                query += " AND"
            }
            query += " EXISTS (SELECT 1 FROM items i WHERE i.supplier_id = s.supplier_id AND i.is_active = true)"
        }
    }

    query += " ORDER BY s.name"

    rows, err := r.db.Pool.Query(ctx, query, params...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var suppliers []*suppliermodels.Supplier
    for rows.Next() {
        supplier := &suppliermodels.Supplier{}
        err := rows.Scan(
            &supplier.SupplierID,
            &supplier.Name,
            &supplier.ContactPerson,
            &supplier.Phone,
            &supplier.Email,
            &supplier.Address,
            &supplier.TaxID,
            &supplier.PaymentTerms,
            &supplier.Notes,
            &supplier.CreatedAt,
            &supplier.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        suppliers = append(suppliers, supplier)
    }

    return suppliers, rows.Err()
}

func (r *PostgresSupplierRepository) GetByID(ctx context.Context, id int) (*suppliermodels.Supplier, error) {
    query := `
        SELECT supplier_id, name, contact_person, phone, email,
               address, tax_id, payment_terms, notes, created_at, updated_at
        FROM suppliers
        WHERE supplier_id = $1
    `

    supplier := &suppliermodels.Supplier{}
    err := r.db.Pool.QueryRow(ctx, query, id).Scan(
        &supplier.SupplierID,
        &supplier.Name,
        &supplier.ContactPerson,
        &supplier.Phone,
        &supplier.Email,
        &supplier.Address,
        &supplier.TaxID,
        &supplier.PaymentTerms,
        &supplier.Notes,
        &supplier.CreatedAt,
        &supplier.UpdatedAt,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }

    return supplier, nil
}

func (r *PostgresSupplierRepository) Create(ctx context.Context, supplier *suppliermodels.Supplier) (int, error) {
    query := `
        INSERT INTO suppliers (
            name, contact_person, phone, email, address,
            tax_id, payment_terms, notes
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING supplier_id
    `

    var id int
    err := r.db.Pool.QueryRow(
        ctx, query,
        supplier.Name,
        supplier.ContactPerson,
        supplier.Phone,
        supplier.Email,
        supplier.Address,
        supplier.TaxID,
        supplier.PaymentTerms,
        supplier.Notes,
    ).Scan(&id)

    if err != nil {
        return 0, err
    }

    return id, nil
}

func (r *PostgresSupplierRepository) Update(ctx context.Context, supplier *suppliermodels.Supplier) error {
    query := `
        UPDATE suppliers SET
            name = $2,
            contact_person = $3,
            phone = $4,
            email = $5,
            address = $6,
            tax_id = $7,
            payment_terms = $8,
            notes = $9
        WHERE supplier_id = $1
    `

    result, err := r.db.Pool.Exec(
        ctx, query,
        supplier.SupplierID,
        supplier.Name,
        supplier.ContactPerson,
        supplier.Phone,
        supplier.Email,
        supplier.Address,
        supplier.TaxID,
        supplier.PaymentTerms,
        supplier.Notes,
    )

    if err != nil {
        return err
    }

    if result.RowsAffected() == 0 {
        return errors.New("supplier not found")
    }

    return nil
}

func (r *PostgresSupplierRepository) Delete(ctx context.Context, id int) error {
    query := `DELETE FROM suppliers WHERE supplier_id = $1`

    result, err := r.db.Pool.Exec(ctx, query, id)
    if err != nil {
        return err
    }

    if result.RowsAffected() == 0 {
        return errors.New("supplier not found")
    }

    return nil
}
