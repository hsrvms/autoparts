package purchases

import (
	"github.com/hsrvms/autoparts/internal/modules/purchases/handlers"
	"github.com/hsrvms/autoparts/internal/modules/purchases/repositories"
	"github.com/hsrvms/autoparts/internal/modules/purchases/services"
	"github.com/hsrvms/autoparts/pkg/db"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(api *echo.Group, database *db.Database) {
    // Initialize repository
    repo := repositories.NewPostgresPurchaseRepository(database)

    // Initialize service
    service := services.NewPurchaseService(repo)

    // Initialize handler
    handler := handlers.NewPurchaseHandler(service)

    // Register routes
    purchases := api.Group("/purchases")
    purchases.GET("", handler.GetPurchases)
    purchases.GET("/:id", handler.GetPurchaseByID)
    purchases.POST("", handler.CreatePurchase)
    purchases.PUT("/:id", handler.UpdatePurchase)
    purchases.DELETE("/:id", handler.DeletePurchase)

    // Additional routes for supplier and item specific purchases
    api.GET("/suppliers/:supplierId/purchases", handler.GetSupplierPurchases)
    api.GET("/items/:itemId/purchases", handler.GetItemPurchases)
}
