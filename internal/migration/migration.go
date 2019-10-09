package migration

import (
	_ "github.com/lib/pq" // package init.
	"gitlab.com/mikrowezel/backend/config"
	"gitlab.com/mikrowezel/backend/migration"
)

// TODO: Implement a more ergonomic way to add migrations.

const (
	devDb  = "granica"
	testDb = "granica_test"
	prodDb = "granica_prod"
)

// Init to explicitly start the migrator.
func Init(cfg *config.Config) *migration.Migrator {
	m := migration.Init(cfg)

	// Migrations
	// Enable Postgis
	mg := &mig{}
	mg.Config(mg.EnablePostgis, mg.DropPostgis)
	m.AddMigration(mg)

	// CreateUsersTable
	mg = &mig{}
	mg.Config(mg.CreateUsersTable, mg.DropUsersTable)
	m.AddMigration(mg)

	return m
}
