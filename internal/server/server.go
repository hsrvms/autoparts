package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/hsrvms/autoparts/pkg/config"
	"github.com/hsrvms/autoparts/pkg/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Server represents our HTTP server
type Server struct {
	Echo   *echo.Echo
	DB     *db.Database
	Config *config.Config
}

// New creates a new server instance
func New(cfg *config.Config, database *db.Database) *Server {
	e := echo.New()

	// Enable middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Serve static files
	e.Static("/static", "web/static")
	e.Static("/js", "web/js")

	// Create server instance
	server := &Server{
		Echo:   e,
		DB:     database,
		Config: cfg,
	}

	// Initialize routes
	server.initRoutes()

	return server
}

// Start starts the HTTP server
func (s *Server) Start() {
	// Server address
	addr := fmt.Sprintf(":%d", s.Config.Server.Port)

	// Create a custom server
	httpServer := &http.Server{
		Addr:         addr,
		ReadTimeout:  s.Config.Server.ReadTimeout,
		WriteTimeout: s.Config.Server.WriteTimeout,
		IdleTimeout:  s.Config.Server.IdleTimeout,
	}

	// Start server
	go func() {
		if err := s.Echo.StartServer(httpServer); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on %s", addr)

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Echo.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to gracefully shutdown server: %v", err)
	}

	log.Println("Server stopped")
}
