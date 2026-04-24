package main

import (
	"log"

	"careerflow/backend/internal/api"
	"careerflow/backend/internal/config"
	"careerflow/backend/internal/db"
)

func main() {
	cfg := config.Load()
	database, err := db.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer database.Close()

	if err := db.RunMigrations(database); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	if err := db.SeedDemoData(database); err != nil {
		log.Fatalf("seed failed: %v", err)
	}

	r := api.NewRouter(cfg, database)
	log.Printf("backend listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
