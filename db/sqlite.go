package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

const databaseName string = "furnace-monitoring-system.db"

var Instance *sql.DB

func CreateDatabase(runMigrations bool) error {
	db, err := sql.Open("sqlite", databaseName)
	if err != nil {
		return fmt.Errorf("db: failed to open the database: %w", err)
	}

	Instance = db

	if runMigrations {
		if err = PerformMeasurementTableMigration(Instance); err != nil {
			return fmt.Errorf("db: failed to perform migrations: %w", err)
		}
	}

	return nil
}
