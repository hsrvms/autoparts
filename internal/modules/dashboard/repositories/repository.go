package repositories

import (
	"context"

	dashboardmodels "github.com/hsrvms/autoparts/internal/modules/dashboard/models"
)

type DashboardRepository interface {
	GetStats(ctx context.Context) (*dashboardmodels.Stats, error)
	GetRecentActivities(ctx context.Context, limit int) ([]*dashboardmodels.Activity, error)
	GetLowStockItems(ctx context.Context, limit int) ([]*dashboardmodels.LowStockItem, error)
}
