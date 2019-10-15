package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	//"github.com/davecgh/go-spew/spew"

	"github.com/jmoiron/sqlx"
	"gitlab.com/mikrowezel/backend/config"
	"gitlab.com/mikrowezel/backend/db"
	"gitlab.com/mikrowezel/backend/log"
	mwmig "gitlab.com/mikrowezel/backend/migration"
	svc "gitlab.com/mikrowezel/backend/service"
	"gitlab.com/mikrowezel/granica/internal/migration"
	"gitlab.com/mikrowezel/granica/internal/model"
	"gitlab.com/mikrowezel/granica/internal/repo"
	"gitlab.com/mikrowezel/granica/pkg/auth"
)

var (
	userDataValid = map[string]string{
		"username":          "username",
		"password":          "password",
		"email":             "username@mail.com",
		"emailConfirmation": "username@mail.com",
		"givenName":         "name",
		"middleNames":       "middles",
		"familyName":        "family",
	}

	userUpdateDataValid = map[string]string{
		"username":          "usernameUpd",
		"password":          "passwordUpd",
		"email":             "usernameUpd@mail.com",
		"emailConfirmation": "usernameUpd@mail.com",
		"givenName":         "nameUpd",
		"middleNames":       "middlesUpd",
		"familyName":        "familyUpd",
	}

	userSample1 = map[string]string{
		"username":          "username1",
		"password":          "password1",
		"email":             "username1@mail.com",
		"emailConfirmation": "username1@mail.com",
		"givenName":         "name1",
		"middleNames":       "middles1",
		"familyName":        "family1",
	}

	userSample2 = map[string]string{
		"username":          "username2",
		"password":          "password2",
		"email":             "username2@mail.com",
		"emailConfirmation": "username2@mail.com",
		"givenName":         "name2",
		"middleNames":       "middles2",
		"familyName":        "family2",
	}
)

func TestMain(m *testing.M) {
	mgr := setup()
	code := m.Run()
	teardown(mgr)
	os.Exit(code)
}

////TestCreateUser tests user repo creation.
//func TestCreateUser(t *testing.T) {
//t.Log("Mock test")
//}

// TestCreateUser tests user repo creation.
func TestCreateUser(t *testing.T) {
	// Setup
	cuReq := auth.CreateUserReq{
		auth.User{
			Username:          userDataValid["username"],
			Password:          userDataValid["password"],
			Email:             userDataValid["email"],
			EmailConfirmation: userDataValid["emailConfirmation"],
			GivenName:         userDataValid["givenName"],
			MiddleNames:       userDataValid["middleNames"],
			FamilyName:        userDataValid["familyName"],
		},
	}

	var cuRes auth.CreateUserRes

	ctx := context.Background()
	cfg := testConfig()
	log := testLogger()

	// Repo
	rhName := "repo-handler"
	rh, err := repo.NewHandler(ctx, cfg, log, rhName)
	rh.Connect()

	// Set service worker
	a := auth.NewWorker(ctx, cfg, log, "auth-worker")
	hs := map[string]svc.Handler{
		rhName: rh,
	}
	a.SetHandlers(hs)

	// Test
	err = a.CreateUser(cuReq, &cuRes)
	if err != nil {
		t.Errorf("create user error: %s", err.Error())
	}

	// Verify
	userVerify, err := getUserByUsername(userDataValid["username"], cfg)
	if err != nil {
		t.Errorf("cannot get user from database: %s", err.Error())
	}

	if userVerify == nil {
		t.Errorf("cannot get user from database")
	}

	// t.Logf("%+v\n", spew.Sdump(user))
	// t.Logf("%+v\n", spew.Sdump(userVerify))

	user := cuRes.User
	if !compareUsers(user, *userVerify) {
		t.Error("User data and its verification does not match.")
	}
}

// Helpers
func getUserByUsername(username string, cfg *config.Config) (*model.User, error) {
	conn, err := getConn()
	if err != nil {
		return nil, err
	}

	schema := cfg.ValOrDef("pg.schema", "public")

	st := `SELECT * FROM %s.users WHERE username='%s';`
	st = fmt.Sprintf(st, schema, username)

	u := &model.User{}
	err = conn.Get(u, st)
	if err != nil {
		msg := fmt.Sprintf("cannot get user: %s", err.Error())
		return nil, errors.New(msg)
	}

	return u, nil
}

func getUserBySlug(username string, cfg *config.Config) (*model.User, error) {
	conn, err := getConn()
	if err != nil {
		return nil, err
	}

	schema := cfg.ValOrDef("pg.schema", "public")

	st := `SELECT * FROM %s.users WHERE slug='%s';`
	st = fmt.Sprintf(st, schema, username)

	u := &model.User{}
	err = conn.Get(u, st)
	if err != nil {
		msg := fmt.Sprintf("cannot get user: %s", err.Error())
		return nil, errors.New(msg)
	}

	return u, nil
}

func compareUsers(user auth.User, toCompare model.User) (areEqual bool) {
	return user.Username == toCompare.Username.String &&
		user.Email == toCompare.Email.String &&
		user.GivenName == toCompare.GivenName.String &&
		user.MiddleNames == toCompare.MiddleNames.String &&
		user.FamilyName == toCompare.FamilyName.String
}

func createSampleUsers() (users []*model.User, err error) {
	ctx := context.Background()
	cfg := testConfig()
	log := testLogger()

	r, err := repo.NewHandler(ctx, cfg, log, "repo-handler")
	if err != nil {
		return users, err
	}
	r.Connect()

	user1 := &model.User{
		Username:          db.ToNullString(userSample1["username"]),
		Password:          userSample2["password1"],
		Email:             db.ToNullString(userSample1["email"]),
		EmailConfirmation: db.ToNullString(userSample1["emailConfirmation"]),
		GivenName:         db.ToNullString(userSample1["givenName"]),
		MiddleNames:       db.ToNullString(userSample1["middleNames"]),
		FamilyName:        db.ToNullString(userSample1["familyName"]),
	}

	err = createUser(r, user1)
	if err != nil {
		return users, err
	}

	users = append(users, user1)

	user2 := &model.User{
		Username:          db.ToNullString(userSample2["username"]),
		Password:          userSample2["password"],
		Email:             db.ToNullString(userSample2["email"]),
		EmailConfirmation: db.ToNullString(userSample2["emailConfirmation"]),
		GivenName:         db.ToNullString(userSample2["givenName"]),
		MiddleNames:       db.ToNullString(userSample2["middleNames"]),
		FamilyName:        db.ToNullString(userSample2["familyName"]),
	}

	err = createUser(r, user2)
	if err != nil {
		return users, err
	}

	users = append(users, user2)

	return users, nil
}

func createUser(r *repo.Repo, user *model.User) error {
	userRepo, err := r.UserRepoNewTx()
	if err != nil {
		return err
	}

	err = userRepo.Create(user)
	if err != nil {
		return err
	}

	err = userRepo.Commit()
	if err != nil {
		return err
	}

	return nil
}

func setup() *mwmig.Migrator {
	m := migration.GetMigrator(testConfig())
	m.RollbackAll()
	m.Migrate()
	return m
}

func teardown(m *mwmig.Migrator) {
	m.RollbackAll()
}

func testConfig() *config.Config {
	cfg := &config.Config{}
	values := map[string]string{
		"pg.host":               "localhost",
		"pg.port":               "5432",
		"pg.schema":             "public",
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

// getConn returns a connection used to
// verify repo insert and update operations.
func getConn() (*sqlx.DB, error) {
	cfg := testConfig()
	conn, err := sqlx.Open("postgres", dbURL(cfg))
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

//
func dbURL(cfg *config.Config) string {
	host := cfg.ValOrDef("pg.host", "localhost")
	port := cfg.ValOrDef("pg.port", "5432")
	schema := cfg.ValOrDef("pg.schema", "public")
	db := cfg.ValOrDef("pg.database", "granica_test")
	user := cfg.ValOrDef("pg.user", "granica")
	pass := cfg.ValOrDef("pg.password", "granica")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=%s", host, port, user, pass, db, schema)
}
