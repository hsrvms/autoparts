package repositories

import (
	"context"
	"time"

	dashboardmodels "github.com/hsrvms/autoparts/internal/modules/dashboard/models"
	"github.com/hsrvms/autoparts/pkg/db"
)

type PostgresDashboardRepository struct {
	db *db.Database
}

func NewPostgresDashboardRepository(database *db.Database) DashboardRepository {
	return &PostgresDashboardRepository{
		db: database,
	}
}

func (r *PostgresDashboardRepository) GetStats(ctx context.Context) (*dashboardmodels.Stats, error) {
	stats := &dashboardmodels.Stats{}

	// Get low stock count
	err := r.db.Pool.QueryRow(ctx, `
        SELECT COUNT(*)
        FROM items
        WHERE current_stock <= minimum_stock AND is_active = true
    `).Scan(&stats.LowStockCount)
	if err != nil {
		return nil, err
	}

	// Get today's sales
	today := time.Now().Format("2006-01-02")
	err = r.db.Pool.QueryRow(ctx, `
        SELECT COALESCE(SUM(total_price), 0)
        FROM sales
        WHERE DATE(date) = $1
    `, today).Scan(&stats.TodaySales)
	if err != nil {
		return nil, err
	}

	// Get active items count
	err = r.db.Pool.QueryRow(ctx, `
        SELECT COUNT(*)
        FROM items
        WHERE is_active = true
    `).Scan(&stats.ActiveItems)
	if err != nil {
		return nil, err
	}

	// Get active suppliers count
	err = r.db.Pool.QueryRow(ctx, `
        SELECT COUNT(DISTINCT supplier_id)
        FROM items
        WHERE is_active = true AND supplier_id IS NOT NULL
    `).Scan(&stats.SupplierCount)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *PostgresDashboardRepository) GetRecentActivities(ctx context.Context, limit int) ([]*dashboardmodels.Activity, error) {
	query := `
        (SELECT
            'sale' as type,
            CONCAT('Sale: ', i.part_number, ' (', s.quantity, ' units)') as message,
            s.created_at as timestamp
        FROM sales s
        JOIN items i ON s.item_id = i.item_id)
        UNION ALL
        (SELECT
            'purchase' as type,
            CONCAT('Purchase: ', i.part_number, ' (', p.quantity, ' units)') as message,
            p.created_at as timestamp
        FROM purchases p
        JOIN items i ON p.item_id = i.item_id)
        ORDER BY timestamp DESC
        LIMIT $1
    `

	rows, err := r.db.Pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []*dashboardmodels.Activity
	for rows.Next() {
		activity := &dashboardmodels.Activity{}
		err := rows.Scan(
			&activity.Type,
			&activity.Message,
			&activity.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	return activities, rows.Err()
}

func (r *PostgresDashboardRepository) GetLowStockItems(ctx context.Context, limit int) ([]*dashboardmodels.LowStockItem, error) {
	query := `
        SELECT
            i.item_id,
            i.part_number,
            i.description,
            i.current_stock,
            i.minimum_stock,
            c.category_name as category
        FROM items i
        LEFT JOIN categories c ON i.category_id = c.category_id
        WHERE i.current_stock <= i.minimum_stock
            AND i.is_active = true
        ORDER BY (i.current_stock::float / i.minimum_stock::float)
        LIMIT $1
    `

	rows, err := r.db.Pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*dashboardmodels.LowStockItem
	for rows.Next() {
		item := &dashboardmodels.LowStockItem{}
		err := rows.Scan(
			&item.ItemID,
			&item.PartNumber,
			&item.Description,
			&item.CurrentStock,
			&item.MinimumStock,
			&item.Category,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}
