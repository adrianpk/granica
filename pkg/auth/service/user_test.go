package service

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
	"gitlab.com/mikrowezel/backend/granica/internal/migration"
	"gitlab.com/mikrowezel/backend/granica/internal/model"
	"gitlab.com/mikrowezel/backend/granica/internal/repo"
	"gitlab.com/mikrowezel/backend/granica/pkg/auth/service"
	tp "gitlab.com/mikrowezel/backend/granica/pkg/auth/transport"
	"gitlab.com/mikrowezel/backend/log"
	mwmig "gitlab.com/mikrowezel/backend/migration"
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

// TestMock does not test anything.
func TestMock(t *testing.T) {
	// Only used to test migration and rollback operations
	// on setup and teardown respectivelly.
}

// TestCreateUser tests user creation.
func TestCreateUser(t *testing.T) {
	// Setup
	req := tp.CreateUserReq{
		tp.User{
			Username:          userDataValid["username"],
			Password:          userDataValid["password"],
			Email:             userDataValid["email"],
			EmailConfirmation: userDataValid["emailConfirmation"],
			GivenName:         userDataValid["givenName"],
			MiddleNames:       userDataValid["middleNames"],
			FamilyName:        userDataValid["familyName"],
		},
	}

	var res tp.CreateUserRes

	ctx := context.Background()
	cfg := testConfig()
	log := testLogger()

	// Repo
	userRepo, err := testRepo(ctx, cfg, log, "repo-handler")
	if err != nil {
		t.Error(err.Error())
	}

	// Service
	s := testService(ctx, cfg, log, userRepo)

	// Test
	err = s.CreateUser(req, &res)
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

	user := res.User
	if !isSameUser(user, *userVerify) {
		t.Error("User data and its verification does not match.")
	}
}

// TestAllIndexUsers tests get all users.
func TestAllIndexUsers(t *testing.T) {
	// Prerequisites
	_, err := createSampleUsers()
	if err != nil {
		t.Errorf("error creating sample users: %s", err.Error())
	}

	// Setup
	req := tp.IndexUsersReq{}

	var res tp.IndexUsersRes

	ctx := context.Background()
	cfg := testConfig()
	log := testLogger()

	// Repo
	userRepo, err := testRepo(ctx, cfg, log, "repo-handler")
	if err != nil {
		t.Error(err.Error())
	}

	// Service
	s := testService(ctx, cfg, log, userRepo)

	// Test
	err = s.IndexUsers(req, &res)
	if err != nil {
		t.Errorf("get users error: %s", err.Error())
	}

	// Verify
	vUsers := res.Users
	if vUsers == nil {
		t.Error("no response")
	}

	if res.Error != "" {
		t.Errorf("Response error: %s", res.Error)
	}

	qty := len(vUsers)
	if qty != 2 {
		t.Errorf("expecting two users got %d", qty)
	}

	if vUsers[0].Username != userSample1["username"] || vUsers[1].Username != userSample2["username"] {
		t.Error("obtained values do not match expected ones")
	}
}

// TestGetUser tests get users by username.
func TestGetUser(t *testing.T) {
	// Prerequisites
	users, err := createSampleUsers()
	if err != nil {
		t.Errorf("error creating sample users: %s", err.Error())
	}

	// Setup
	req := tp.GetUserReq{
		tp.Identifier{
			Username: users[0].Username.String,
		},
	}

	var res tp.GetUserRes

	ctx := context.Background()
	cfg := testConfig()
	log := testLogger()

	r, err := repo.NewHandler(ctx, cfg, log, "repo-handler")
	if err != nil {
		t.Errorf("cannot initialize repo handler: %s", err.Error())
	}
	r.Connect()

	// Repo
	userRepo, err := testRepo(ctx, cfg, log, "repo-handler")
	if err != nil {
		t.Error(err.Error())
	}

	// Service
	s := testService(ctx, cfg, log, userRepo)

	// Test
	err = s.GetUser(req, &res)
	if err != nil {
		t.Errorf("get user error: %s", err.Error())
	}

	// Verify
	if res.Error != "" {
		t.Errorf("Response error: %s", res.Error)
	}

	user := res.User
	if user.Username != userSample1["username"] {
		t.Error("obtained values do not match expected ones")
	}
}

// TestUpdateUser user repo update.
func TestUpdateUser(t *testing.T) {
	// Prerequisites
	users, err := createSampleUsers()
	if err != nil {
		t.Errorf("error creating sample users: %s", err.Error())
	}

	// Setup
	user := users[0]
	req := tp.UpdateUserReq{
		tp.Identifier{
			Username: user.Username.String,
		},
		tp.User{
			Username:          userUpdateDataValid["username"],
			Password:          userUpdateDataValid["password"],
			Email:             userUpdateDataValid["email"],
			EmailConfirmation: userUpdateDataValid["emailConfirmation"],
			GivenName:         userUpdateDataValid["givenName"],
			MiddleNames:       userUpdateDataValid["middleNames"],
			FamilyName:        userUpdateDataValid["familyName"],
		},
	}

	var res tp.UpdateUserRes

	ctx := context.Background()
	cfg := testConfig()
	log := testLogger()

	// Repo
	userRepo, err := testRepo(ctx, cfg, log, "repo-handler")
	if err != nil {
		t.Error(err.Error())
	}

	// Service
	s := testService(ctx, cfg, log, userRepo)

	// Test
	err = s.UpdateUser(req, &res)
	if err != nil {
		t.Errorf("update user error: %s", err.Error())
	}

	userVerify, err := getUserByUsername(userSample1["username"], cfg)
	if err != nil {
		t.Errorf("cannot get user from database: %s", err.Error())
	}

	if userVerify == nil {
		t.Errorf("cannot get user from database")
	}

	// TODO: Add accurate check of all updated fields.
	if userVerify.Email.String != userUpdateDataValid["email"] {
		t.Error("obtained values do not match expected ones")
	}
}

// TestDeleteUser tests delete users from repo.
func TestDeleteUser(t *testing.T) {
	// Prerequisites
	users, err := createSampleUsers()
	if err != nil {
		t.Errorf("error creating sample users: %s", err.Error())
	}

	// Setup
	user := users[0]
	req := tp.DeleteUserReq{
		Identifier: tp.Identifier{
			Username: user.Username.String,
		},
	}

	var res tp.DeleteUserRes

	ctx := context.Background()
	cfg := testConfig()
	log := testLogger()

	// Repo
	userRepo, err := testRepo(ctx, cfg, log, "repo-handler")
	if err != nil {
		t.Error(err.Error())
	}

	// Service
	s := testService(ctx, cfg, log, userRepo)

	// Test
	err = s.DeleteUser(req, &res)
	if err != nil {
		t.Errorf("delete user error: %s", err.Error())
	}

	// Verify
	vUser, err := getUserBySlug(user.Slug.String, cfg)
	if err != nil {
		return
	}

	if vUser == nil {
		return
	}

	if vUser.Username.String == user.Username.String {
		t.Error("user was not deleted from database")
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

func isSameUser(user tp.User, toCompare model.User) bool {
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
		Password:          userSample1["password"],
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
	// m.Reset()
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

func testRepo(ctx context.Context, cfg *config.Config, log *log.Logger, name string) (*repo.Repo, error) {
	rh, err := repo.NewHandler(ctx, cfg, log, name)
	if err != nil {
		return nil, err
	}
	rh.Connect()
	if err != nil {
		return nil, err
	}
	return rh, nil
}

func testService(ctx context.Context, cfg *config.Config, log *log.Logger, r *repo.Repo) *service.Service {
	s := service.MakeService(ctx, cfg, log)
	s.SetRepo(r)
	return s
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

// dbURL returns a Postgres connection string.
func dbURL(cfg *config.Config) string {
	host := cfg.ValOrDef("pg.host", "localhost")
	port := cfg.ValOrDef("pg.port", "5432")
	schema := cfg.ValOrDef("pg.schema", "public")
	db := cfg.ValOrDef("pg.database", "granica_test")
	user := cfg.ValOrDef("pg.user", "granica")
	pass := cfg.ValOrDef("pg.password", "granica")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=%s", host, port, user, pass, db, schema)
}
