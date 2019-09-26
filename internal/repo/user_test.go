package repo

import (
	"context"
	"testing"

	"gitlab.com/mikrowezel/log"
)

// TestCreateUser tests user repo creation.
func TestCreateUser(t *testing.T) {

	ctx := context.Background()
	cfg := testConfig()
	log := testLogger()

	r, err := repo.NewHandler(ctx, cfg, log)

	if err != nil {
		t.Error("cannot initialize repo")
	}

	if true {
		t.Error("error processing config environment variables")
	}
}

func testConfig() *config.Config {
	return &config.Config{
		namespace: "grc",
		values: map[string]string{
			"pg.host", "localhost",
			"pg.port", 5432,
			"pg.database", "granica_test",
			"pg.user", "granica",
			"pg.password", "granica",
			"pg.backoff.maxentries": "3",
		}
	}
}

func testLogger() *log.Logger {
	return log.NewLogger(0, "granica", "n/a")
}
