package db

import (
	"database/sql"
	"fmt"

	"github.com/Krzysztofz01/furnace-monitoring-system/log"
	_ "modernc.org/sqlite"
)

const databaseName string = "db/furnace-monitoring-system.db"

var Instance *sql.DB

func CreateDatabase(runMigrations bool) error {
	db, err := sql.Open("sqlite", databaseName)
	if err != nil {
		return fmt.Errorf("db: failed to open the database: %w", err)
	}

	Instance = db
	log.Instance.Debugln("Database driver instance initialzied successful.")

	if runMigrations {
		log.Instance.Debugln("Starting the database migration process.")

		if err = PerformMeasurementTableMigration(Instance); err != nil {
			return fmt.Errorf("db: failed to perform migrations: %w", err)
		}

		log.Instance.Debugln("Database migrations finished successful.")
	}

	return nil
}
