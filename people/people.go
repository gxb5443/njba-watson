package people

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type Person struct {
	Id           string    `json:",omitempty"`
	FirstName    string    `json:"first_name,omitempty" db:"first_name"`
	LastName     string    `json:"last_name,omitempty" db:"last_name"`
	MiddleName   string    `json:"middle_name,omitempty" db:"middle_name"`
	Suffix       string    `json:",omitempty"`
	Prefix       string    `json:",omitempty"`
	Title        string    `json:",omitempty"`
	Address1     string    `json:",omitempty"`
	Address2     string    `json:",omitempty"`
	Zip          string    `json:",omitempty"`
	State        string    `json:",omitempty"`
	Country      string    `json:",omitempty"`
	HomePhone    string    `json:"home_phone,omitempty"`
	CellPhone    string    `json:"cell_phone,omitempty"`
	EmailAddress string    `json:"email_address,omitempty"`
	Source       string    `json:"source,omitempty"`
	Created      time.Time `json:"source,omitempty"`
}

type nullPerson struct {
	Id           sql.NullString `json:",omitempty"`
	FirstName    sql.NullString `json:"first_name,omitempty" db:"first_name"`
	LastName     sql.NullString `json:"last_name,omitempty" db:"last_name"`
	MiddleName   sql.NullString `json:"middle_name,omitempty" db:"middle_name"`
	Suffix       sql.NullString `json:",omitempty"`
	Prefix       sql.NullString `json:",omitempty"`
	Title        sql.NullString `json:",omitempty"`
	Address1     sql.NullString `json:",omitempty"`
	Address2     sql.NullString `json:",omitempty"`
	Zip          sql.NullString `json:",omitempty"`
	State        sql.NullString `json:",omitempty"`
	Country      sql.NullString `json:",omitempty"`
	HomePhone    sql.NullString `json:"home_phone,omitempty" db:"home_phone"`
	CellPhone    sql.NullString `json:"cell_phone,omitempty" db:"cell_phone"`
	EmailAddress sql.NullString `json:"email_address,omitempty" db:"email_address"`
	Source       sql.NullString `json:"source,omitempty"`
	Created      time.Time      `json:"source,omitempty"`
}

func GetAll(db *sqlx.DB) ([]*Person, error) {
	query := ` SELECT id, first_name, last_name FROM people;	`
	var p []*Person
	err := db.Select(&p, query)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func GetById(db *sqlx.DB, id string) (*Person, error) {
	query := `SELECT id, first_name, last_name, middle_name, suffix, prefix, title, address1, address2, zip, state, country, home_phone, cell_phone, email_address FROM people WHERE id=$1 LIMIT 1`
	var c nullPerson
	err := db.Get(&c, query, id)
	if err != nil {
		return nil, err
	}
	out := new(Person)
	if c.Id.Valid {
		out.Id = c.Id.String
	}
	if c.FirstName.Valid {
		out.FirstName = c.FirstName.String
	}
	if c.MiddleName.Valid {
		out.MiddleName = c.MiddleName.String
	}
	if c.LastName.Valid {
		out.LastName = c.LastName.String
	}
	if c.HomePhone.Valid {
		out.HomePhone = c.HomePhone.String
	}
	if c.CellPhone.Valid {
		out.CellPhone = c.CellPhone.String
	}
	if c.Suffix.Valid {
		out.Suffix = c.Suffix.String
	}
	if c.Prefix.Valid {
		out.Prefix = c.Prefix.String
	}
	if c.Title.Valid {
		out.Title = c.Title.String
	}
	if c.Address1.Valid {
		out.Address1 = c.Address1.String
	}
	if c.Address2.Valid {
		out.Address2 = c.Address2.String
	}
	if c.Zip.Valid {
		out.Zip = c.Zip.String
	}
	if c.State.Valid {
		out.State = c.State.String
	}
	if c.Country.Valid {
		out.Country = c.Country.String
	}
	out.Created = c.Created
	return out, nil
}
