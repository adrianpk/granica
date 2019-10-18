package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"gitlab.com/mikrowezel/backend/config"
	logger "gitlab.com/mikrowezel/backend/log"
	"gitlab.com/mikrowezel/granica/internal/model"
)

type (
	AccountRepo struct {
		ctx context.Context
		cfg *config.Config
		log *logger.Logger
		Tx  *sqlx.Tx
	}
)

func makeAccountRepo(ctx context.Context, cfg *config.Config, log *logger.Logger, tx *sqlx.Tx) *AccountRepo {
	return &AccountRepo{
		ctx: ctx,
		cfg: cfg,
		log: log,
		Tx:  tx,
	}
}

// Create a account in repo.
func (ur *AccountRepo) Create(account *model.Account) error {
	account.SetCreateValues()

	st := `INSERT INTO accounts (id, slug, name, account_type, owner_id, parent_id, email, shown_name, geolocation, locale, base_tz, current_tz, starts_at, ends_at, is_active, is_deleted, created_by_id, updated_by_id, created_at, updated_at)
VALUES (:id, :slug, :name, :account_type, :owner_id, :parent_id, :email, :shown_name, :geolocation, :locale, :base_tz, :current_tz, :starts_at, :ends_at, :is_active, :is_deleted, :created_by_id, :updated_by_id, :created_at, :updated_at)`

	_, err := ur.Tx.NamedExec(st, account)

	return err
}

// GetAll accounts from repo.
func (ur *AccountRepo) GetAll() (accounts []model.Account, err error) {
	st := `SELECT * FROM accounts;`

	err = ur.Tx.Select(&accounts, st)

	return accounts, err
}

// Get account by ID.
func (ur *AccountRepo) Get(id interface{}) (model.Account, error) {
	var account model.Account

	st := `SELECT * FROM ACCOUNTS WHERE id = '%s' LIMIT 1;`
	st = fmt.Sprintf(st, id.(string))

	err := ur.Tx.Get(&account, st)

	return account, err
}

// GetBySlug account from repo by slug.
func (ur *AccountRepo) GetBySlug(slug string) (model.Account, error) {
	var account model.Account

	st := `SELECT * FROM ACCOUNTS WHERE slug = '%s' LIMIT 1;`
	st = fmt.Sprintf(st, slug)

	err := ur.Tx.Get(&account, st)

	return account, err
}

// Update account data in repo.
func (ur *AccountRepo) Update(account *model.Account) error {
	ref, err := ur.Get(account.ID.String())
	if err != nil {
		return fmt.Errorf("cannot retrieve reference account: %s", err.Error())
	}

	account.SetUpdateValues()

	var st strings.Builder
	pcu := false // previous column updated?

	st.WriteString("UPDATE accounts SET ")

	if account.Name.String != ref.Name.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("name", "name"))
		pcu = true
	}

	if account.OwnerID.String != ref.OwnerID.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("owner_id", "owner_id"))
		pcu = true
	}

	if account.ParentID.String != ref.ParentID.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("parent_id", "parent_id"))
		pcu = true
	}

	if account.AccountType.String != ref.AccountType.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("account_type", "account_type"))
		pcu = true
	}

	if account.Email.String != ref.Email.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("email", "email"))
		pcu = true
	}

	if account.ShownName.String != ref.ShownName.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("shown_name", "shown_name"))
		pcu = true
	}

	st.WriteString(" ")
	st.WriteString(whereID(ref.ID.String()))
	st.WriteString(";")

	//fmt.Println(st.String())

	_, err = ur.Tx.NamedExec(st.String(), account)

	return err
}

// Delete account from repo by ID.
func (ur *AccountRepo) Delete(id string) error {
	st := `DELETE FROM ACCOUNTS WHERE id = '%s';`
	st = fmt.Sprintf(st, id)

	_, err := ur.Tx.Exec(st)

	return err
}

// DeleteBySlug account from repo by slug.
func (ur *AccountRepo) DeleteBySlug(slug string) error {
	st := `DELETE FROM ACCOUNTS WHERE slug = '%s';`
	st = fmt.Sprintf(st, slug)

	_, err := ur.Tx.Exec(st)

	return err
}

// Commit transaction
func (ur *AccountRepo) Commit() error {
	return ur.Tx.Commit()
}
