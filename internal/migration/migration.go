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
	// 00000001
	mg := &mig{}
	mg.SetFx(mg.Up00000001)
	m.AddMigration(mg)
	// 00000002
	mg = &mig{}
	mg.SetFx(mg.Up00000001)
	m.AddMigration(mg)

	// Rollbacks
	// 00000001
	mg = &mig{}
	mg.SetFx(mg.Down00000001)
	m.AddRollback(mg)
	// 00000002
	mg = &mig{}
	mg.SetFx(mg.Down00000002)
	m.AddRollback(mg)

	return m
}
