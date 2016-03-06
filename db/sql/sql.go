package sql

// Embed migration SQL
//go:generate go-bindata -pkg sql -o migrations.go migrations/

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/rubenv/sql-migrate"
)

// Initialize the database and applies migration scripts.
func Initialize(db *sqlx.DB) error {
	// Apply migrations
	steps := &migrate.AssetMigrationSource{
		Asset:    Asset,
		AssetDir: AssetDir,
		Dir:      "migrations",
	}

	n, err := migrate.Exec(db.DB, db.DriverName(), steps, migrate.Up)
	if n == 0 && err == nil {
		log.Printf("Database scheme is up to date.\n")
		return nil
	}
	log.Printf("Applied %d migrations to the database.\n", n)
	if err != nil {
		log.Printf("Failed to apply migration steps to database (Error: %s).\n", err)
		return err
	}
	return nil
}
