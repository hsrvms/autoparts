package inventory

import (
	"github.com/hsrvms/autoparts/internal/modules/inventory/handlers"
	"github.com/hsrvms/autoparts/internal/modules/inventory/repositories"
	"github.com/hsrvms/autoparts/internal/modules/inventory/services"
	"github.com/hsrvms/autoparts/pkg/db"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(api *echo.Group, database *db.Database) {
	// Initialize repository
	repo := repositories.NewPostgresInventoryRepository(database)

	// Initialize service
	service := services.NewInventoryService(repo)

	// Initialize handler
	handler := handlers.NewInventoryHandler(service)

	// Item routes
	items := api.Group("/items")
	items.GET("", handler.GetItems)
	items.GET("/low-stock", handler.GetLowStockItems)
	items.GET("/:id", handler.GetItemByID)
	items.GET("/barcode/:barcode", handler.GetItemByBarcode)
	items.POST("", handler.CreateItem)
	items.PUT("/:id", handler.UpdateItem)
	items.DELETE("/:id", handler.DeleteItem)
	items.GET("/barcode/:barcode/image", handler.GetBarcodeImage)

	// Compatibility routes
	items.GET("/:itemId/compatibilities", handler.GetCompatibilities)
	items.POST("/:itemId/compatibilities", handler.AddCompatibility)
	items.DELETE("/:itemId/compatibilities/:submodelId", handler.RemoveCompatibility)
	api.GET("/submodels/:submodelId/compatible-items", handler.GetCompatibleItems)
}
