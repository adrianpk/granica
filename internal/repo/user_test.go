package repo

import (
	"context"
	"testing"

	"gitlab.com/mikrowezel/db"
	"gitlab.com/mikrowezel/granica/internal/migration"
	"gitlab.com/mikrowezel/granica/internal/model"
	"gitlab.com/mikrowezel/granica/internal/repo"

	"gitlab.com/mikrowezel/config"
	"gitlab.com/mikrowezel/log"
)

// TestCreateUser tests user repo creation.
func TestCreateUser(t *testing.T) {
	// TODO: move to test setup and teardown function.
	migration.Init()
	m := migration.Migrator()
	m.MigrateAll()

	// Valid user data
	user := &model.User{
		Username:          db.ToNullString("username"),
		Password:          "password",
		Email:             db.ToNullString("username@mail.com"),
		EmailConfirmation: db.ToNullString("username@mail.com"),
		GivenName:         db.ToNullString("name"),
		MiddleNames:       db.ToNullString("middles"),
		FamilyName:        db.ToNullString("family"),
	}

	ctx := context.Background()
	cfg := testConfig()
	log := testLogger()

	r, err := repo.NewHandler(ctx, cfg, log, "repo-handler")
	if err != nil {
		t.Error("cannot initialize repo handler")
	}
	r.Connect()

	userRepo, err := r.UserRepoNewTx()
	if err != nil {
		t.Error("cannot initialize user repo")
	}

	err = userRepo.Create(user)
	if err != nil {
		t.Log(err)
		t.Error("create user error")
	}

	err = userRepo.Commit()
	if err != nil {
		t.Log(err)
		t.Error("create user commit error")
	}

}

func testConfig() *config.Config {
	cfg := &config.Config{}
	values := map[string]string{
		"pg.host":               "localhost",
		"pg.port":               "5432",
		"pg.database":           "granica_test",
		"pg.user":               "granica",
		"pg.password":           "granica",
		"pg.backoff.maxentries": "3",
	}

	cfg.SetNamespace("grc")
	cfg.SetValues(values)
	return cfg
}

func testLogger() *log.Logger {
	return log.NewDevLogger(0, "granica", "n/a")
}
