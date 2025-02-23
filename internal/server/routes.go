package server

import (
	"net/http"

	"github.com/hsrvms/autoparts/internal/modules/categories"
	"github.com/hsrvms/autoparts/internal/modules/dashboard"
	"github.com/hsrvms/autoparts/internal/modules/inventory"
	"github.com/hsrvms/autoparts/internal/modules/purchases"
	"github.com/hsrvms/autoparts/internal/modules/sales"
	"github.com/hsrvms/autoparts/internal/modules/suppliers"
	"github.com/hsrvms/autoparts/internal/modules/vehicles"
	"github.com/labstack/echo/v4"
)

func (s *Server) initRoutes() {
	api := s.Echo.Group("/api")

	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	api.GET("/version", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"version": "1.0.0"})
	})

	dashboard.RegisterRoutes(s.Echo, api, s.DB)
	categories.RegisterRoutes(api, s.DB)
	vehicles.RegisterRoutes(api, s.DB)
	inventory.RegisterRoutes(api, s.DB)
	suppliers.RegisterRoutes(api, s.DB)
	purchases.RegisterRoutes(api, s.DB)
	sales.RegisterRoutes(api, s.DB)
}
