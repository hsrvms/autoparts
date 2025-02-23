package suppliers

import (
	"github.com/hsrvms/autoparts/internal/modules/suppliers/handlers"
	"github.com/hsrvms/autoparts/internal/modules/suppliers/repositories"
	"github.com/hsrvms/autoparts/internal/modules/suppliers/services"
	"github.com/hsrvms/autoparts/pkg/db"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(api *echo.Group, database *db.Database) {
    // Initialize repository
    repo := repositories.NewPostgresSupplierRepository(database)

    // Initialize service
    service := services.NewSupplierService(repo)

    // Initialize handler
    handler := handlers.NewSupplierHandler(service)

    // Register routes
    suppliers := api.Group("/suppliers")
    suppliers.GET("", handler.GetSuppliers)
    suppliers.GET("/:id", handler.GetSupplierByID)
    suppliers.POST("", handler.CreateSupplier)
    suppliers.PUT("/:id", handler.UpdateSupplier)
    suppliers.DELETE("/:id", handler.DeleteSupplier)
}
