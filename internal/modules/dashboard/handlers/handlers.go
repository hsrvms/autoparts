package handlers

import (
	"net/http"

	"github.com/hsrvms/autoparts/internal/modules/dashboard/services"
	"github.com/labstack/echo/v4"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

type DashboardHandler struct {
	service services.DashboardService
}

func NewDashboardHandler(service services.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		service: service,
	}
}

// RenderDashboard renders the main dashboard page
func (h *DashboardHandler) RenderDashboard(c echo.Context) error {
	return c.Render(http.StatusOK, "base.html", map[string]interface{}{
		"Title":  "Dashboard",
		"Active": "dashboard",
	})
}

// GetStats handles the HTMX request for dashboard stats
func (h *DashboardHandler) GetStats(c echo.Context) error {
	ctx := c.Request().Context()
	stats, err := h.service.GetStats(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// If it's an HTMX request for a specific stat
	if stat := c.QueryParam("stat"); stat != "" {
		switch stat {
		case "low-stock-count":
			return c.String(http.StatusOK, formatNumber(stats.LowStockCount))
		case "today-sales":
			return c.String(http.StatusOK, formatCurrency(stats.TodaySales))
		case "active-items":
			return c.String(http.StatusOK, formatNumber(stats.ActiveItems))
		case "supplier-count":
			return c.String(http.StatusOK, formatNumber(stats.SupplierCount))
		}
	}

	return c.JSON(http.StatusOK, stats)
}

// GetRecentActivities handles the HTMX request for recent activities
func (h *DashboardHandler) GetRecentActivities(c echo.Context) error {
	ctx := c.Request().Context()
	activities, err := h.service.GetRecentActivities(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.Render(http.StatusOK, "dashboard/partials/activities.html", map[string]interface{}{
		"Activities": activities,
	})
}

// GetLowStockItems handles the HTMX request for low stock items
func (h *DashboardHandler) GetLowStockItems(c echo.Context) error {
	ctx := c.Request().Context()
	items, err := h.service.GetLowStockItems(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.Render(http.StatusOK, "dashboard/partials/low-stock.html", map[string]interface{}{
		"Items": items,
	})
}

// Add these helper functions
var printer = message.NewPrinter(language.English)

func formatNumber(n int) string {
	return printer.Sprint(number.Decimal(n))
}

func formatCurrency(amount float64) string {
	return printer.Sprintf("$%.2f", amount)
}
