package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"gitlab.com/mikrowezel/backend/config"
	"gitlab.com/mikrowezel/backend/granica/internal/model"
	logger "gitlab.com/mikrowezel/backend/log"
)

type (
	ProfileRepo struct {
		ctx context.Context
		cfg *config.Config
		log *logger.Logger
		Tx  *sqlx.Tx
	}
)

func makeProfileRepo(ctx context.Context, cfg *config.Config, log *logger.Logger, tx *sqlx.Tx) *ProfileRepo {
	return &ProfileRepo{
		ctx: ctx,
		cfg: cfg,
		log: log,
		Tx:  tx,
	}
}

// Create a profile in repo.
func (ur *ProfileRepo) Create(profile *model.Profile) error {
	profile.SetCreateValues()

	st := `INSERT INTO profiles (id, tenant_id, slug, owner_id, profile_type, name, email, description, location, bio, moto,website, aniversary_date, avatar_path, header_path,geolocation, locale, base_tz, current_tz, starts_at, ends_at, is_active, is_deleted, created_by_id, updated_by_id, created_at, updated_at)
VALUES (:id, tenant_id, :slug, :owner_id, :profile_type, :name, :email, :description, :location, :bio, :moto, :website, :aniversary_date, :avatar_path, :header_path, :geolocation, :locale, :base_tz, :current_tz, :starts_at, :ends_at, :is_active, :is_deleted, :created_by_id, :updated_by_id, :created_at, :updated_at)`

	_, err := ur.Tx.NamedExec(st, profile)

	return err
}

// GetAll profiles from repo.
func (ur *ProfileRepo) GetAll() (profiles []model.Profile, err error) {
	st := `SELECT * FROM profiles;`

	err = ur.Tx.Select(&profiles, st)

	return profiles, err
}

// Get profile by ID.
func (ur *ProfileRepo) Get(id interface{}) (model.Profile, error) {
	var profile model.Profile

	st := `SELECT * FROM PROFILES WHERE id = '%s' LIMIT 1;`
	st = fmt.Sprintf(st, id.(string))

	err := ur.Tx.Get(&profile, st)

	return profile, err
}

// GetBySlug profile from repo by slug.
func (ur *ProfileRepo) GetBySlug(slug string) (model.Profile, error) {
	var profile model.Profile

	st := `SELECT * FROM PROFILES WHERE slug = '%s' LIMIT 1;`
	st = fmt.Sprintf(st, slug)

	err := ur.Tx.Get(&profile, st)

	return profile, err
}

// Update profile data in repo.
func (ur *ProfileRepo) Update(profile *model.Profile) error {
	ref, err := ur.Get(profile.ID.String())
	if err != nil {
		return fmt.Errorf("cannot retrieve reference profile: %s", err.Error())
	}

	profile.SetUpdateValues()

	var st strings.Builder
	pcu := false // previous column updated?

	st.WriteString("UPDATE profiles SET ")

	if profile.OwnerID.String != ref.OwnerID.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("owner_id", "owner_id"))
		pcu = true
	}

	if profile.ProfileType.String != ref.ProfileType.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("profile_type", "profile_type"))
		pcu = true
	}

	if profile.Name.String != ref.Name.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("name", "name"))
		pcu = true
	}

	if profile.Email.String != ref.Email.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("email", "email"))
		pcu = true
	}

	if profile.Description.String != ref.Description.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("description", "description"))
		pcu = true
	}

	if profile.Location.String != ref.Location.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("location", "location"))
		pcu = true
	}

	if profile.Bio.String != ref.Bio.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("bio", "bio"))
		pcu = true
	}

	if profile.Moto.String != ref.Moto.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("moto", "moto"))
		pcu = true
	}

	if profile.Website.String != ref.Website.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("website", "website"))
		pcu = true
	}

	if profile.AniversaryDate.Time != ref.AniversaryDate.Time {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("aniversary_date", "aniversary_date"))
		pcu = true
	}

	if profile.AvatarPath.String != ref.AvatarPath.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("avatar_path", "avatar_path"))
		pcu = true
	}

	if profile.HeaderPath.String != ref.HeaderPath.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("header_path", "header_path"))
		pcu = true
	}

	if profile.Geolocation.Point != ref.Geolocation.Point {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("geolocation", "geolocation"))
		pcu = true
	}

	if profile.Locale.String != ref.Locale.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("locale", "locale"))
		pcu = true
	}

	if profile.BaseTZ.String != ref.BaseTZ.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("base_tz", "base_tz"))
		pcu = true
	}

	if profile.CurrentTZ.String != ref.CurrentTZ.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("current_tz", "current_tz"))
		pcu = true
	}

	if profile.IsActive.Bool != ref.IsActive.Bool {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("is_active", "is_active"))
		pcu = true
	}

	if profile.IsDeleted.Bool != ref.IsDeleted.Bool {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("is_deleted", "is_deleted"))
		pcu = true
	}

	if profile.CreatedByID.String != ref.CreatedByID.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("created_by_id", "created_by_id"))
		pcu = true
	}

	if profile.UpdatedByID.String != ref.UpdatedByID.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("updated_by_id", "updated_by_id"))
		pcu = true
	}

	if profile.CreatedAt.Time != ref.CreatedAt.Time {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("created_at", "created_at"))
		pcu = true
	}

	if profile.UpdatedAt.Time != ref.UpdatedAt.Time {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("updated_at", "updated_at"))
		pcu = true
	}

	st.WriteString(" ")
	st.WriteString(whereID(ref.ID.String()))
	st.WriteString(";")

	//fmt.Println(st.String())

	_, err = ur.Tx.NamedExec(st.String(), profile)

	return err
}

// Delete profile from repo by ID.
func (ur *ProfileRepo) Delete(id string) error {
	st := `DELETE FROM PROFILES WHERE id = '%s';`
	st = fmt.Sprintf(st, id)

	_, err := ur.Tx.Exec(st)

	return err
}

// DeleteBySlug profile from repo by slug.
func (ur *ProfileRepo) DeleteBySlug(slug string) error {
	st := `DELETE FROM PROFILES WHERE slug = '%s';`
	st = fmt.Sprintf(st, slug)

	_, err := ur.Tx.Exec(st)

	return err
}

// Commit transaction
func (ur *ProfileRepo) Commit() error {
	return ur.Tx.Commit()
}
