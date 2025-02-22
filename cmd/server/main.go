package main

import (
	"log"

	"github.com/hsrvms/autoparts/internal/server"
	"github.com/hsrvms/autoparts/pkg/config"
	"github.com/hsrvms/autoparts/pkg/db"
)

func main() {
	cfg := config.New()

	database, err := db.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	srv := server.New(cfg, database)
	srv.Start()
}
