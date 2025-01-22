package db

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
)

// RunMigrations automatically applies all pending migrations
func RunMigrations(databaseURL string, migrationPath string, databasename string, dbversion uint) error {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return err
	}
	defer db.Close()

	driver, err := migratePostgres.WithInstance(db, &migratePostgres.Config{})
	if err != nil {
		log.Fatalf("Error creating migration driver: %v", err)
		return err
	}

	absMigrationPath, err := filepath.Abs(migrationPath)
	if err != nil {
		log.Fatalf("Error getting absolute path: %v", err)
		return err
	}

	sourceURL := "file://" + absMigrationPath

	m, err := migrate.NewWithDatabaseInstance(sourceURL, databasename, driver)
	if err != nil {
		log.Fatalf("Error creating migration instance: %v", err)
		return err
	}

	version, _, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatalf("Error getting migration version: %v", err)
		return err
	}

	if err == migrate.ErrNilVersion {
		version = 0
	}

	if version == dbversion {
		fmt.Printf("Database is already at version %d\n", version)
		return nil
	}

	err = m.Migrate(dbversion)

	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Error applying migrations: %v", err)
		return err
	}

	return nil
}
