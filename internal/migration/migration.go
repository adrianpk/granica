package migration

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // package init.
	"gitlab.com/mikrowezel/config"
)

// TODO: Refactor to make it a generically usable module.
// TODO: Move to its own module (mikrowezel/migration)
// NOTE: This is a work in progress, not ready for production.

const (
	devDb  = "granica"
	testDb = "granica_test"
	prodDb = "granica_prod"
)

var (
	mig *migrator
)

// Init to explicitly start the migrator.
func Init(cfg *config.Config) {
	mig = &migrator{cfg: cfg}
	err := mig.Connect()
	if err != nil {
		os.Exit(1)
	}

	// Migrations
	// TODO: build a helper to create Migration struct
	mig.makeMigration(mig.Up00000001)
	mig.makeMigration(mig.Up00000002)

	// Rollbacks
	mig.makeRollback(mig.Down00000001)
	mig.makeRollback(mig.Down00000002)
}

func Migrator() *migrator {
	return mig
}

func (m *migrator) makeMigration(f func() procResult) {
	tx := transaction{conn: m.conn}
	tx.function = f
	m.AddUp(&migration{proc{tx: tx}})
}

func (m *migrator) makeRollback(f func() procResult) {
	tx := transaction{conn: m.conn}
	tx.function = f
	m.AddDown(&rollback{proc{tx: tx}})
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
	m.down = append(m.down, rb)
}

func (m *migrator) MigrateAll() error {
	for i, _ := range m.up {
		fn := fmt.Sprintf("Up%08d", i+1)
		reflect.ValueOf(m).MethodByName(fn).Call([]reflect.Value{})
	}

	return nil
}

func (m *migrator) RollbackAll() error {
	top := len(m.down) - 1
	for i := top; i >= 0; i-- {
		fn := fmt.Sprintf("Down%08d", i+1)
		reflect.ValueOf(m).MethodByName(fn).Call([]reflect.Value{})
	}

	return nil
}

func (m *migrator) MigrateThis(mg migration) error {
	return nil
}

func (m *migrator) RollbackThis(r rollback) error {
	return nil
}

func (m *migrator) makeProcResult(tx *sqlx.Tx, name string, err error) procResult {
	return procResult{
		tx:   tx,
		name: name,
		err:  err,
	}
}

func (m *migrator) dbURL() string {
	host := m.cfg.ValOrDef("pg.host", "localhost")
	port := m.cfg.ValOrDef("pg.port", "5432")
	db := m.cfg.ValOrDef("pg.database", "granica_test_d1x89s0l")
	user := m.cfg.ValOrDef("pg.user", "granica")
	pass := m.cfg.ValOrDef("pg.password", "granica")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, db)
}
