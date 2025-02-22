package server

import (
	"net/http"

	"github.com/hsrvms/autoparts/internal/modules/categories"
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

	categories.RegisterRoutes(api, s.DB)
	vehicles.RegisterRoutes(api, s.DB)

}
