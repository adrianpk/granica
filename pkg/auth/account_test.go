package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	//"github.com/davecgh/go-spew/spew"

	"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx"
	"gitlab.com/mikrowezel/backend/config"
	"gitlab.com/mikrowezel/backend/db"
	"gitlab.com/mikrowezel/backend/log"
	mwmig "gitlab.com/mikrowezel/backend/migration"
	svc "gitlab.com/mikrowezel/backend/service"
	"gitlab.com/mikrowezel/backend/granica/internal/migration"
	"gitlab.com/mikrowezel/backend/granica/internal/model"
	"gitlab.com/mikrowezel/backend/granica/internal/repo"
	"gitlab.com/mikrowezel/backend/granica/pkg/auth"
)

var (
	accountDataValid = map[string]string{
		"tenantID":    "localhost",
		"name":        "name",
		"ownerID":     "ba3b11b3-947b-4536-8958-8c77185c06a7",
		"parentID":    "24f696d1-453b-4d32-bdfe-8b0261c3cb16",
		"accountType": "user",
		"email":       "username@mail.com",
		"shownName":   "shownname",
	}

	accountUpdateDataValid = map[string]string{
		"tenantID":    "localhost",
		"name":        "nameUpd",
		"ownerID":     "ba3b11b3-947b-4536-8958-8c77185c06a7",
		"parentID":    "24f696d1-453b-4d32-bdfe-8b0261c3cb16",
		"accountType": "userUpd",
		"email":       "usernameUpd@mail.com",
		"shownName":   "shownnameUpd",
	}

	accountSample1 = map[string]string{
		"tenantID":    "localhost",
		"name":        "name1",
		"ownerID":     "d5882c7b-9838-429b-98ec-5025c238a91f",
		"parentID":    "4fa05de8-b91c-4b3d-8461-07bedab4738a",
		"parentId":    "-",
		"accountType": "user1",
		"email":       "username1@mail.com",
		"shownName":   "shownname1",
	}

	accountSample2 = map[string]string{
		"tenantID":    "localhost",
		"name":        "name2",
		"ownerID":     "4f7612d2-88aa-43a6-8f8b-065e1c782664",
		"parentID":    "517abf68-3dd7-4b3f-989f-fcf4aa26a103",
		"accountType": "user2",
		"email":       "username2@mail.com",
		"shownName":   "shownname2",
	}
)

func TestMain(m *testing.M) {
	mgr := setup()
	code := m.Run()
	teardown(mgr)
	os.Exit(code)
}

// TestCreateAccount tests account creation.
func TestCreateAccount(t *testing.T) {
	// Setup
	req := auth.CreateAccountReq{
		auth.Account{
			TenantID:    accountDataValid["tenantId"],
			Name:        accountDataValid["name"],
			OwnerID:     accountDataValid["ownerID"],
			ParentID:    accountDataValid["parentID"],
			AccountType: accountDataValid["accountType"],
			Email:       accountDataValid["email"],
			ShownName:   accountDataValid["shownName"],
		},
	}

	var res auth.CreateAccountRes

	ctx := context.Background()
	cfg := testConfig()
	log := testLogger()

	// Repo
	userRepo, err := testRepo(ctx, cfg, log, "repo-handler")
	if err != nil {
		t.Error(err.Error())
	}
	// Auth
	a := testAuth(ctx, cfg, log, "auth-handler", userRepo)

	// Test
	err = a.CreateAccount(req, &res)
	if err != nil {
		t.Errorf("create account error: %s", err.Error())
	}

	// Verify
	account := res.Account
	accountVerify, err := getAccountBySlug(account.Slug, cfg)
	if err != nil {
		t.Errorf("cannot get account from database: %s", err.Error())
	}

	if accountVerify == nil {
		t.Errorf("cannot get account from database")
	}

	if !isSameAccount(account, *accountVerify) {
		t.Logf("%+v\n", spew.Sdump(account))
		t.Logf("%+v\n", spew.Sdump(accountVerify))

		t.Error("Account data and its verification does not match.")
	}
}

// TestGetAccounts tests get all accounts.
func TestGetAccounts(t *testing.T) {
	// Prerequisites
	_, err := createSampleAccounts()
	if err != nil {
		t.Errorf("error creating sample accounts: %s", err.Error())
	}

	// Setup
	req := auth.GetAccountsReq{}

	var res auth.GetAccountsRes

	ctx := context.Background()
	cfg := testConfig()
	log := testLogger()

	// Repo
	userRepo, err := testRepo(ctx, cfg, log, "repo-handler")
	if err != nil {
		t.Error(err.Error())
	}
	// Auth
	a := testAuth(ctx, cfg, log, "auth-handler", userRepo)

	// Test
	err = a.GetAccounts(req, &res)
	if err != nil {
		t.Errorf("get accounts error: %s", err.Error())
	}

	// Verify
	vAccounts := res.Accounts
	if vAccounts == nil {
		t.Error("no response")
	}

	if res.Error != "" {
		t.Errorf("Response error: %s", res.Error)
	}

	qty := len(vAccounts)
	if qty != 2 {
		t.Errorf("expecting two accounts got %d", qty)
	}

	if vAccounts[0].Slug != accountSample1["slug"] || vAccounts[1].Slug != accountSample2["slug"] {
		t.Error("obtained values do not match expected ones")
	}
}

// TestGetAccount tests get accounts by slug.
func TestGetAccount(t *testing.T) {
	// Prerequisites
	accounts, err := createSampleAccounts()
	if err != nil {
		t.Errorf("error creating sample accounts: %s", err.Error())
	}

	// Setup
	req := auth.GetAccountReq{
		auth.Identifier{
			Slug: accounts[0].Slug.String,
		},
	}

	var res auth.GetAccountRes

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
	// Auth
	a := testAuth(ctx, cfg, log, "auth-handler", userRepo)

	// Test
	err = a.GetAccount(req, &res)
	if err != nil {
		t.Errorf("get account error: %s", err.Error())
	}

	// Verify
	if res.Error != "" {
		t.Errorf("Response error: %s", res.Error)
	}

	accountRes := res.Account
	if accountRes.Name != accountSample1["name"] {
		t.Error("obtained values do not match expected ones")
	}
}

// TestUpdateAccount account repo update.
func TestUpdateAccount(t *testing.T) {
	// Prerequisites
	accounts, err := createSampleAccounts()

	if err != nil {
		t.Errorf("error creating sample accounts: %s", err.Error())
	}

	// Setup
	account := accounts[0]
	req := auth.UpdateAccountReq{
		auth.Identifier{
			Slug: account.Slug.String,
		},
		auth.Account{
			TenantID:    accountUpdateDataValid["tenantId"],
			Name:        accountUpdateDataValid["name"],
			OwnerID:     accountUpdateDataValid["ownerID"],
			ParentID:    accountUpdateDataValid["parentID"],
			AccountType: accountUpdateDataValid["accountType"],
			Email:       accountUpdateDataValid["email"],
			ShownName:   accountUpdateDataValid["shownName"],
		},
	}

	var res auth.UpdateAccountRes

	ctx := context.Background()
	cfg := testConfig()
	log := testLogger()

	// Repo
	userRepo, err := testRepo(ctx, cfg, log, "repo-handler")
	if err != nil {
		t.Error(err.Error())
	}
	// Auth
	a := testAuth(ctx, cfg, log, "auth-handler", userRepo)

	// Test
	err = a.UpdateAccount(req, &res)
	if err != nil {
		t.Errorf("update account error: %s", err.Error())
	}

	// Verify
	accountRes := res.Account
	accountVerify, err := getAccountBySlug(accountRes.Slug, cfg)
	if err != nil {
		t.Errorf("cannot get account from database: %s", err.Error())
	}

	if accountVerify == nil {
		t.Errorf("cannot get account from database")
	}

	// TODO: Add accurate check of all updated fields.
	if accountVerify.Email.String != accountUpdateDataValid["email"] {
		t.Error("obtained values do not match expected ones")
	}
}

// TestDeleteAccount tests delete accounts from repo.
func TestDeleteAccount(t *testing.T) {
	// Prerequisites
	accounts, err := createSampleAccounts()
	if err != nil {
		t.Errorf("error creating sample accounts: %s", err.Error())
	}

	// Setup
	account := accounts[0]
	req := auth.DeleteAccountReq{
		auth.Identifier{
			Slug: account.Slug.String,
		},
	}

	var res auth.DeleteAccountRes

	ctx := context.Background()
	cfg := testConfig()
	log := testLogger()

	// Repo
	userRepo, err := testRepo(ctx, cfg, log, "repo-handler")
	if err != nil {
		t.Error(err.Error())
	}
	// Auth
	a := testAuth(ctx, cfg, log, "auth-handler", userRepo)

	// Test
	err = a.DeleteAccount(req, &res)
	if err != nil {
		t.Errorf("delete account error: %s", err.Error())
	}

	// Verify
	vAccount, err := getAccountBySlug(account.Slug.String, cfg)
	if err != nil {
		return
	}

	if vAccount == nil {
		return
	}

	if vAccount.Slug.String == account.Slug.String {
		t.Error("account was not deleted from database")
	}
}

func getAccountBySlug(slug string, cfg *config.Config) (*model.Account, error) {
	conn, err := getConn()
	if err != nil {
		return nil, err
	}

	schema := cfg.ValOrDef("pg.schema", "public")

	st := `SELECT * FROM %s.accounts WHERE slug='%s';`
	st = fmt.Sprintf(st, schema, slug)

	u := &model.Account{}
	err = conn.Get(u, st)
	if err != nil {
		msg := fmt.Sprintf("cannot get account: %s", err.Error())
		return nil, errors.New(msg)
	}

	return u, nil
}

func isSameAccount(account auth.Account, toCompare model.Account) bool {
	return account.TenantID == toCompare.TenantID.String &&
		account.Slug == toCompare.Slug.String &&
		account.Name == toCompare.Name.String &&
		account.OwnerID == toCompare.OwnerID.String &&
		account.ParentID == toCompare.ParentID.String &&
		account.AccountType == toCompare.AccountType.String &&
		account.Email == toCompare.Email.String &&
		account.ShownName == toCompare.ShownName.String
}

func createSampleAccounts() (accounts []*model.Account, err error) {
	ctx := context.Background()
	cfg := testConfig()
	log := testLogger()

	r, err := repo.NewHandler(ctx, cfg, log, "repo-handler")
	if err != nil {
		return accounts, err
	}
	r.Connect()

	account1 := &model.Account{
		Name:        db.ToNullString(accountSample1["name"]),
		OwnerID:     db.ToNullString(accountSample1["ownerID"]),
		ParentID:    db.ToNullString(accountSample1["parentID"]),
		AccountType: db.ToNullString(accountSample1["accountType"]),
		Email:       db.ToNullString(accountSample1["email"]),
		ShownName:   db.ToNullString(accountSample1["shownName"]),
	}

	err = createAccount(r, account1)
	if err != nil {
		return accounts, err
	}

	accounts = append(accounts, account1)

	account2 := &model.Account{
		Name:        db.ToNullString(accountSample2["name"]),
		OwnerID:     db.ToNullString(accountSample2["ownerID"]),
		ParentID:    db.ToNullString(accountSample2["parentID"]),
		AccountType: db.ToNullString(accountSample2["accountType"]),
		Email:       db.ToNullString(accountSample2["email"]),
		ShownName:   db.ToNullString(accountSample1["shownName"]),
	}

	err = createAccount(r, account2)
	if err != nil {
		return accounts, err
	}

	accounts = append(accounts, account2)

	return accounts, nil
}

func createAccount(r *repo.Repo, account *model.Account) error {
	accountRepo, err := r.AccountRepoNewTx()
	if err != nil {
		return err
	}

	account.SetCreateValues()
	err = accountRepo.Create(account)
	if err != nil {
		return err
	}

	err = accountRepo.Commit()
	if err != nil {
		return err
	}

	return nil
}

func setup() *mwmig.Migrator {
	m := migration.GetMigrator(testConfig())
	m.Reset()
	//m.RollbackAll()
	//m.Migrate()
	return m
}

func teardown(m *mwmig.Migrator) {
	// m.RollbackAll()
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

func testAuth(ctx context.Context, cfg *config.Config, log *log.Logger, name string, rh *repo.Repo) *auth.Auth {
	a := auth.NewWorker(ctx, cfg, log, name)
	hs := map[string]svc.Handler{
		rh.Name(): rh,
	}
	a.SetHandlers(hs)
	return a
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
