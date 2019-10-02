package migration

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // package init.
)

// TODO: Refactor to make it a generically usable module.
// TODO: Move to its own module (mikrowezel/migration)
// NOTE: This is a work in progress, not ready for production.

var (
	mig *migrator
)

// Init to explicitly start the migrator.
func Init() {
	mig = &migrator{}
	err := mig.Connect()
	if err != nil {
		os.Exit(1)
	}

	// Migrations
	// TODO: build a helper to create :Migration struct
	// 00000001
	mig.makeMigration(mig.Up00000001)
	mig.makeMigration(mig.Up00000002)

	// Rollbacks
}

func Migrator() *migrator {
	return mig
}

func (m *migrator) makeMigration(f func() error) {
	tx := transaction{conn: m.conn}
	tx.function = f
	m.AddUp(&migration{proc{tx: tx}})
}

func (m *migrator) getTx() *sqlx.Tx {
	return m.conn.MustBegin()
}

func (m *migrator) Connect() error {
	conn, err := sqlx.Open("postgres", m.dbURL())
	if err != nil {
		log.Printf("Connection error: %s\n", err.Error())
		return err
	}

	err = conn.Ping()
	if err != nil {
		log.Printf("Connection error: %s", err.Error())
		return err
	}

	m.conn = conn
	return nil
}

func (m *migrator) CreateDb() error {
	return nil
}

func (m *migrator) DropDb() error {
	return nil
}

func (m *migrator) Reset() error {
	err := m.DropDb()
	if err != nil {
		log.Printf("Drop database error: %s", err.Error())
		// Do't return maybe it was not created before.
	}

	err = m.CreateDb()
	if err != nil {
		log.Printf("Drop database error: %s", err.Error())
		return err
	}

	err = m.MigrateAll()
	if err != nil {
		log.Printf("Drop database error: %s", err.Error())
		return err
	}

	return nil
}

func (m *migrator) AddUp(mg *migration) {
	m.up = append(m.up, mg)
}

func (m *migrator) AddDown(rb *rollback) {

}

func (m *migrator) MigrateAll() error {
	for i, _ := range m.up {
		// FIX: quick and dirty formatter just fot testing.
		// Does properly work only for i < 10.
		fn := fmt.Sprintf("Up0000000%d", i+1)
		reflect.ValueOf(m).MethodByName(fn).Call([]reflect.Value{})
	}
	return nil
}

func (m *migrator) RollbackAll() error {
	return nil
}

func (m *migrator) MigrateThis(mg migration) error {
	return nil
}

func (m *migrator) RollbackThis(r rollback) error {
	return nil
}

func (m *migrator) dbURL() string {
	// TODO: make these values configurable
	host := "localhost"
	port := "5432"
	db := "granica_test"
	user := "granica"
	pass := "granica"
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, db)
}
