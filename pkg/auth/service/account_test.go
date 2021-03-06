package service

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
	"gitlab.com/mikrowezel/backend/granica/internal/migration"
	"gitlab.com/mikrowezel/backend/granica/internal/model"
	"gitlab.com/mikrowezel/backend/granica/internal/repo"
	"gitlab.com/mikrowezel/backend/granica/pkg/auth/service"
	tp "gitlab.com/mikrowezel/backend/granica/pkg/auth/transport"
	"gitlab.com/mikrowezel/backend/log"
	mig "gitlab.com/mikrowezel/backend/migration"
)

var (
	accountDataValid = map[string]string{
		"tenantID":    "localhost",
		"name":        "name",
		"ownerID":     "ba3b11b3-947b-4536-8958-8c77185c06a7",
		"parentID":    "24f696d1-453b-4d32-bdfe-8b0261c3cb16",
		"accountType": "user",
		"email":       "username@mail.com",
	}

	accountUpdateDataValid = map[string]string{
		"tenantID":    "localhost",
		"name":        "nameUpd",
		"ownerID":     "ba3b11b3-947b-4536-8958-8c77185c06a7",
		"parentID":    "24f696d1-453b-4d32-bdfe-8b0261c3cb16",
		"accountType": "userUpd",
		"email":       "usernameUpd@mail.com",
	}

	accountSample1 = map[string]string{
		"tenantID":    "localhost",
		"name":        "name1",
		"ownerID":     "d5882c7b-9838-429b-98ec-5025c238a91f",
		"parentID":    "4fa05de8-b91c-4b3d-8461-07bedab4738a",
		"parentId":    "-",
		"accountType": "user1",
		"email":       "username1@mail.com",
	}

	accountSample2 = map[string]string{
		"tenantID":    "localhost",
		"name":        "name2",
		"ownerID":     "4f7612d2-88aa-43a6-8f8b-065e1c782664",
		"parentID":    "517abf68-3dd7-4b3f-989f-fcf4aa26a103",
		"accountType": "user2",
		"email":       "username2@mail.com",
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

// TestCreateAccount tests account creation.
func TestCreateAccount(t *testing.T) {
	// Prerequisites
	users, err := createSampleUsers()
	if err != nil {
		t.Errorf("error creating sample users: %s", err.Error())
	}

	// Setup
	req := tp.CreateAccountReq{
		tp.Account{
			TenantID:    accountDataValid["tenantId"],
			Name:        accountDataValid["name"],
			OwnerID:     users[0].ID.String(),
			ParentID:    accountDataValid["parentID"],
			AccountType: accountDataValid["accountType"],
			Email:       accountDataValid["email"],
		},
	}

	var res tp.CreateAccountRes

	ctx := context.Background()
	cfg := testConfig()
	log := testLogger()

	// Repo
	accountRepo, err := testRepo(ctx, cfg, log, "repo-handler")
	if err != nil {
		t.Error(err.Error())
	}

	// Service
	s := testService(ctx, cfg, log, accountRepo)

	// Test
	err = s.CreateAccount(req, &res)

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

// TestGetAllAccounts tests get all accounts.
func TestGetAllAccounts(t *testing.T) {
	// Prerequisites
	_, err := createSampleAccounts()
	if err != nil {
		t.Errorf("error creating sample accounts: %s", err.Error())
	}

	// Setup
	req := tp.GetAccountsReq{}

	var res tp.GetAccountsRes

	ctx := context.Background()
	cfg := testConfig()
	log := testLogger()

	// Repo
	accountRepo, err := testRepo(ctx, cfg, log, "repo-handler")
	if err != nil {
		t.Error(err.Error())
	}

	// Service
	s := testService(ctx, cfg, log, accountRepo)

	// Test
	err = s.GetAccounts(req, &res)
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
	req := tp.GetAccountReq{
		tp.Identifier{
			Slug: accounts[0].Slug.String,
		},
	}

	var res tp.GetAccountRes

	ctx := context.Background()
	cfg := testConfig()
	log := testLogger()

	r, err := repo.NewHandler(ctx, cfg, log, "repo-handler")
	if err != nil {
		t.Errorf("cannot initialize repo handler: %s", err.Error())
	}
	r.Connect()

	// Repo
	accountRepo, err := testRepo(ctx, cfg, log, "repo-handler")
	if err != nil {
		t.Error(err.Error())
	}

	// Service
	s := testService(ctx, cfg, log, accountRepo)

	// Test
	err = s.GetAccount(req, &res)
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
	req := tp.UpdateAccountReq{
		tp.Identifier{
			Slug: account.Slug.String,
		},
		tp.Account{
			TenantID:    accountUpdateDataValid["tenantId"],
			Name:        accountUpdateDataValid["name"],
			AccountType: accountUpdateDataValid["accountType"],
			Email:       accountUpdateDataValid["email"],
		},
	}

	var res tp.UpdateAccountRes

	ctx := context.Background()
	cfg := testConfig()
	log := testLogger()

	// Repo
	accountRepo, err := testRepo(ctx, cfg, log, "repo-handler")
	if err != nil {
		t.Error(err.Error())
	}

	// Service
	s := testService(ctx, cfg, log, accountRepo)

	// Test
	err = s.UpdateAccount(req, &res)
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
	req := tp.DeleteAccountReq{
		tp.Identifier{
			Slug: account.Slug.String,
		},
	}

	var res tp.DeleteAccountRes

	ctx := context.Background()
	cfg := testConfig()
	log := testLogger()

	// Repo
	accountRepo, err := testRepo(ctx, cfg, log, "repo-handler")
	if err != nil {
		t.Error(err.Error())
	}

	// Service
	s := testService(ctx, cfg, log, accountRepo)

	// Test
	err = s.DeleteAccount(req, &res)
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

func isSameAccount(account tp.Account, toCompare model.Account) bool {
	return account.TenantID == toCompare.TenantID.String &&
		account.Slug == toCompare.Slug.String &&
		account.Name == toCompare.Name.String &&
		account.OwnerID == toCompare.OwnerID.String &&
		account.ParentID == toCompare.ParentID.String &&
		account.AccountType == toCompare.AccountType.String &&
		account.Email == toCompare.Email.String
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

func createSampleAccounts() (accounts []*model.Account, err error) {
	// Prerequisites
	users, err := createSampleUsers()
	if err != nil {
		return nil, err
	}

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
		OwnerID:     db.ToNullString(users[0].ID.String()),
		ParentID:    db.ToNullString(accountSample1["parentID"]),
		AccountType: db.ToNullString(accountSample1["accountType"]),
		Email:       db.ToNullString(accountSample1["email"]),
	}

	err = createAccount(r, account1)
	if err != nil {
		return accounts, err
	}

	accounts = append(accounts, account1)

	account2 := &model.Account{
		Name:        db.ToNullString(accountSample2["name"]),
		OwnerID:     db.ToNullString(users[1].ID.String()),
		ParentID:    db.ToNullString(accountSample2["parentID"]),
		AccountType: db.ToNullString(accountSample2["accountType"]),
		Email:       db.ToNullString(accountSample2["email"]),
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

func setup() *mig.Migrator {
	m := migration.GetMigrator(testConfig())
	//m.Reset()
	m.RollbackAll()
	m.Migrate()
	return m
}

func teardown(m *mig.Migrator) {
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
