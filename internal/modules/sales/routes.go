package sales

import (
	"github.com/hsrvms/autoparts/internal/modules/sales/handlers"
	"github.com/hsrvms/autoparts/internal/modules/sales/repositories"
	"github.com/hsrvms/autoparts/internal/modules/sales/services"
	"github.com/hsrvms/autoparts/pkg/db"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(api *echo.Group, database *db.Database) {
    // Initialize repository
    repo := repositories.NewPostgresSaleRepository(database)

    // Initialize service
    service := services.NewSaleService(repo)

    // Initialize handler
    handler := handlers.NewSaleHandler(service)

    // Register routes
    sales := api.Group("/sales")
    sales.GET("", handler.GetSales)
    sales.GET("/:id", handler.GetSaleByID)
    sales.POST("", handler.CreateSale)
    sales.PUT("/:id", handler.UpdateSale)
    sales.DELETE("/:id", handler.DeleteSale)
    sales.GET("/transaction/:transactionNumber", handler.GetByTransactionNumber)
    sales.GET("/customer/:customerEmail", handler.GetCustomerSales)
}
