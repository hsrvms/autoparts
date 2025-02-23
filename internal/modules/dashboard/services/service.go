package services

import (
	"context"

	dashboardmodels "github.com/hsrvms/autoparts/internal/modules/dashboard/models"
	"github.com/hsrvms/autoparts/internal/modules/dashboard/repositories"
)

type DashboardService interface {
    GetStats(ctx context.Context) (*dashboardmodels.Stats, error)
    GetRecentActivities(ctx context.Context) ([]*dashboardmodels.Activity, error)
    GetLowStockItems(ctx context.Context) ([]*dashboardmodels.LowStockItem, error)
}

type dashboardService struct {
    repo repositories.DashboardRepository
}

func NewDashboardService(repo repositories.DashboardRepository) DashboardService {
    return &dashboardService{
        repo: repo,
    }
}

func (s *dashboardService) GetStats(ctx context.Context) (*dashboardmodels.Stats, error) {
    return s.repo.GetStats(ctx)
}

func (s *dashboardService) GetRecentActivities(ctx context.Context) ([]*dashboardmodels.Activity, error) {
    return s.repo.GetRecentActivities(ctx, 10) // Show last 10 activities
}

func (s *dashboardService) GetLowStockItems(ctx context.Context) ([]*dashboardmodels.LowStockItem, error) {
    return s.repo.GetLowStockItems(ctx, 5) // Show top 5 low stock items
}
