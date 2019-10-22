package model

import (
	"database/sql"

	"github.com/lib/pq"
	"gitlab.com/mikrowezel/backend/db"
	m "gitlab.com/mikrowezel/backend/model"
	"golang.org/x/crypto/bcrypt"
)

//// Identification -------------------------------------------------------------
//type (
//// Identification model
//Identification struct {
//ID       uuid.UUID      `db:"id" json:"id"`
//TenantID sql.NullString `db:"tenant_id" json:"tenantID"`
//Slug     sql.NullString `db:"slug" json:"slug"`
//}
//)

//// GetID representation.
//func (i *Identification) GetID() interface{} {
//return i.ID.Value
//}

//// SetID for user.
//func (i *Identification) SetID(id uuid.UUID) {
//i.ID = id
//}

//// GenID for user.
//func (i *Identification) GenID() {
//if i.ID == uuid.Nil {
//i.ID = uuid.NewV4()
//}
//}

//// UpdateSlug if it was not set.
//func (i *Identification) UpateSlug(prefix string) (slug string, err error) {
//if strings.Trim(i.Slug.String, " ") == "" {
//s, err := i.genSlug(prefix)
//if err != nil {
//return "", err
//}
//i.Slug = db.ToNullString(s)
//}
//return i.Slug.String, nil
//}

//// genSlug if it was not set.
//func (i *Identification) genSlug(prefix string) (slug string, err error) {
//if strings.TrimSpace(prefix) == "" {
//return "", errors.New("no slug prefix defined")
//}

////--
//prefix = strings.Replace(prefix, "-", "", -1)
//prefix = strings.Replace(prefix, "_", "", -1)

//if !utf8.ValidString(prefix) {
//v := make([]rune, 0, len(prefix))
//for i, r := range prefix {
//if r == utf8.RuneError {
//_, size := utf8.DecodeRuneInString(prefix[i:])
//if size == 1 {
//continue
//}
//}
//v = append(v, r)
//}
//prefix = string(v)
//}

//prefix = strings.ToLower(prefix)
////----

//s := strings.Split(uuid.NewV4().String(), "-")
//l := s[len(s)-1]

//return strings.ToLower(fmt.Sprintf("%s-%s", prefix, l)), nil
//}

//// SetCreateValues sets de ID and slug.
//func (i *Identification) SetCreateValues(slugPrefix string) error {
//i.GenID()
//_, err := i.UpateSlug(slugPrefix)
//if err != nil {
//return err
//}
//return nil
//}

//// Audit ----------------------------------------------------------------------
//type Audit struct {
//CreatedByID sql.NullString `db:"created_by_id" json:"createdByID"`
//UpdatedByID sql.NullString `db:"updated_by_id" json:"updatedByID"`
//CreatedAt   pq.NullTime    `db:"created_at" json:"createdAt"`
//UpdatedAt   pq.NullTime    `db:"updated_at" json:"updatedAt"`
//}

//// SetCreateValues sets de ID and slug.
//func (a *Audit) SetCreateValues() error {
//now := time.Now()
//a.CreatedAt = pg.ToNullTime(now)
//a.UpdatedAt = pg.NullTime()
//return nil
//}

//// SetUpdateValues
//func (a *Audit) SetUpdateValues() error {
//now := time.Now()
//a.UpdatedAt = pg.ToNullTime(now)
//return nil
//}

//// MarshalJSON is a custom MarshalJSON function.
//func (audit *Audit) MarshalJSON() ([]byte, error) {
//type Alias Audit
//return json.Marshal(&struct {
//*Alias
//CreatedAt int64 `json:"createdAt"`
//UpdatedAt int64 `json:"updatedAt"`
//}{
//Alias:     (*Alias)(audit),
//CreatedAt: audit.CreatedAt.Time.Unix(),
//UpdatedAt: audit.UpdatedAt.Time.Unix(),
//})
//}

//// UnmarshalJSON is a custom UnmarshalJSON function.
//func (audit *Audit) UnmarshalJSON(data []byte) error {
//type Alias Audit
//aux := &struct {
//*Alias
//CreatedAt int64 `json:"createdAt"`
//UpdatedAt int64 `json:"updatedAt"`
//}{
//Alias:     (*Alias)(audit),
//CreatedAt: audit.CreatedAt.Time.Unix(),
//UpdatedAt: audit.UpdatedAt.Time.Unix(),
//}
//if err := json.Unmarshal(data, &aux); err != nil {
//return err
//}
//tc := time.Unix(aux.CreatedAt, 0)
//tu := time.Unix(aux.UpdatedAt, 0)
//audit.CreatedAt = pq.NullTime{tc, true}
//audit.UpdatedAt = pq.NullTime{tu, true}
//return nil
//}

// User -----------------------------------------------------------------------
type (
	// User model
	User struct {
		m.Identification
		Username          sql.NullString `db:"username" json:"username"`
		Password          string         `db:"password" json:"password"`
		PasswordDigest    sql.NullString `db:"password_digest" json:"-"`
		Email             sql.NullString `db:"email" json:"email"`
		EmailConfirmation sql.NullString `db:"emailConfirmation" json:"emailConfirmation"`
		GivenName         sql.NullString `db:"given_name" json:"givenName"`
		MiddleNames       sql.NullString `db:"middle_names" json:"middleNames"`
		FamilyName        sql.NullString `db:"family_name" json:"familyName"`
		Geolocation       db.NullPoint   `db:"geolocation" json:"geolocation"`
		Locale            sql.NullString `db:"locale" json:"locale"`
		BaseTZ            sql.NullString `db:"base_tz" json:"baseTZ"`
		CurrentTZ         sql.NullString `db:"current_tz" json:"currentTZ"`
		StartsAt          pq.NullTime    `db:"starts_at" json:"startsAt"`
		EndsAt            pq.NullTime    `db:"ends_at" json:"endsAt"`
		IsActive          sql.NullBool   `db:"is_active" json:"isActive"`
		IsDeleted         sql.NullBool   `db:"is_deleted" json:"isDeleted"`
		m.Audit
	}
)

// UpdatePasswordDigest if password changed.
func (user *User) UpdatePasswordDigest() (digest string, err error) {
	if user.Password == "" {
		return user.PasswordDigest.String, nil
	}

	hpass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return user.PasswordDigest.String, err
	}
	user.PasswordDigest.String = string(hpass)
	return user.PasswordDigest.String, nil
}

// SetCreateValues sets de ID and slug.
func (user *User) SetCreateValues() error {
	pfx := user.Username.String
	user.Identification.SetCreateValues(pfx)
	user.Audit.SetCreateValues()
	return nil
}

// SetUpdateValues
func (user *User) SetUpdateValues() error {
	user.Audit.SetUpdateValues()
	return nil
}

// Account --------------------------------------------------------------------
type (
	// Account model
	Account struct {
		m.Identification
		Name        sql.NullString `db:"name" json:"name"`
		OwnerID     sql.NullString `db:"owner_id" json:"ownerID"`
		ParentID    sql.NullString `db:"parent_id" json:"parentID"`
		AccountType sql.NullString `db:"account_type" json:"accountType"`
		Email       sql.NullString `db:"email" json:"email"`
		ShownName   sql.NullString `db:"shown_name" json:"shownName"`
		Geolocation db.NullPoint   `db:"geolocation" json:"geolocation"`
		Locale      sql.NullString `db:"locale" json:"locale"`
		BaseTZ      sql.NullString `db:"base_tz" json:"baseTZ"`
		CurrentTZ   sql.NullString `db:"current_tz" json:"currentTZ"`
		StartsAt    pq.NullTime    `db:"starts_at" json:"startsAt"`
		EndsAt      pq.NullTime    `db:"ends_at" json:"endsAt"`
		IsActive    sql.NullBool   `db:"is_active" json:"isActive"`
		IsDeleted   sql.NullBool   `db:"is_deleted" json:"isDeleted"`
		m.Audit
	}
)

// SetCreateValues sets de ID and slug.
func (account *Account) SetCreateValues() error {
	pfx := account.Name.String
	account.Identification.SetCreateValues(pfx)
	account.Audit.SetCreateValues()
	return nil
}

// SetUpdateValues
func (account *Account) SetUpdateValues() error {
	account.Audit.SetUpdateValues()
	return nil
}
