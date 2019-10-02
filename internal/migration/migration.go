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
	Mig *Migrator
)

func Init() {
	Mig = &Migrator{}
	err := Mig.Connect()
	if err != nil {
		os.Exit(1)
	}

	// Migrations
	// TODO: build a helper to create :Migration struct
	// 00000001
	tx1 := transaction{conn: Mig.conn}
	tx1.function = tx1.Up00000001
	Mig.AddUp(&Migration{proc{tx: tx1}})
	// 00000002
	tx2 := transaction{conn: Mig.conn}
	tx2.function = tx2.Up00000002
	Mig.AddUp(&Migration{proc{tx: tx2}})

	// Rollbacks
}

func (t *transaction) getTx() *sqlx.Tx {
	return t.conn.MustBegin()
}

func (m *Migrator) Connect() error {
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

func (m *Migrator) CreateDb() error {
	return nil
}

func (m *Migrator) DropDb() error {
	return nil
}

func (m *Migrator) Reset() error {
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

func (m *Migrator) AddUp(mg *Migration) {
	m.up = append(m.up, mg)
}

func (m *Migrator) AddDown(rb *Rollback) {

}

func (m *Migrator) MigrateAll() error {
	log.Printf("Migration: %+v\n", m.up)
	for i, mg := range m.up {
		// FIX: quick and dirty formatter just fot testing.
		// Does properly work only for i < 10.
		fn := fmt.Sprintf("Up0000000%d", i+1)
		reflect.ValueOf(&mg.tx).MethodByName(fn).Call([]reflect.Value{})
	}
	return nil
}

func (m *Migrator) RollbackAll() error {
	return nil
}

func (m *Migrator) MigrateThis(mg Migration) error {
	return nil
}

func (m *Migrator) RollbackThis(r Rollback) error {
	return nil
}

func (m *Migrator) dbURL() string {
	// TODO: make these values configurable
	host := "localhost"
	port := "5432"
	db := "granica_test"
	user := "granica"
	pass := "granica"
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, db)
}
