package vehicles

import (
	"github.com/hsrvms/autoparts/internal/modules/vehicles/handlers"
	"github.com/hsrvms/autoparts/internal/modules/vehicles/repositories"
	"github.com/hsrvms/autoparts/internal/modules/vehicles/services"
	"github.com/hsrvms/autoparts/pkg/db"
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all vehicle-related routes
func RegisterRoutes(api *echo.Group, database *db.Database) {
	// Initialize repository
	repo := repositories.NewPostgresVehicleRepository(database)

	// Initialize service
	service := services.NewVehicleService(repo)

	// Initialize handler
	handler := handlers.NewVehicleHandler(service)

	// Vehicle makes routes
	makes := api.Group("/makes")
	makes.GET("", handler.GetAllMakes)
	makes.GET("/:id", handler.GetMakeByID)
	makes.POST("", handler.CreateMake)
	makes.PUT("/:id", handler.UpdateMake)
	makes.DELETE("/:id", handler.DeleteMake)
	makes.GET("/:makeId/models", handler.GetModelsByMake) // Get models for a specific make

	// Vehicle models routes
	models := api.Group("/models")
	models.GET("", handler.GetAllModels)
	models.GET("/:id", handler.GetModelByID)
	models.POST("", handler.CreateModel)
	models.PUT("/:id", handler.UpdateModel)
	models.DELETE("/:id", handler.DeleteModel)
	models.GET("/:modelId/submodels", handler.GetSubmodelsByModel) // Get submodels for a specific model

	// Vehicle submodels routes
	submodels := api.Group("/submodels")
	submodels.GET("", handler.GetAllSubmodels)
	submodels.GET("/:id", handler.GetSubmodelByID)
	submodels.POST("", handler.CreateSubmodel)
	submodels.PUT("/:id", handler.UpdateSubmodel)
	submodels.DELETE("/:id", handler.DeleteSubmodel)
}
