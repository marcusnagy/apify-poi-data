package main

import (
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Import the postgres driver
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Import the file source driver
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib" // Import the pgx driver

	"apify-poi-data/config"
	"apify-poi-data/db"
)

func main() {
	cfg := config.LoadDatabaseConfig()

	err := db.RunMigrations(
		cfg.Database.URL,
		cfg.Database.MigrationsPath,
		cfg.Database.DatabaseName,
		cfg.Database.DatabaseVersion,
	)
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	log.Println("Migrations applied successfully")
}
