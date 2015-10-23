package people

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type Person struct {
	Id           string    `json:",omitempty" db:"id"`
	FirstName    string    `json:"first_name,omitempty" db:"first_name"`
	LastName     string    `json:"last_name,omitempty" db:"last_name"`
	MiddleName   string    `json:"middle_name,omitempty" db:"middle_name"`
	Suffix       string    `json:",omitempty" db:"suffix"`
	Prefix       string    `json:",omitempty" db:"prefix"`
	Title        string    `json:",omitempty" db:"title"`
	Address1     string    `json:",omitempty" db:"address1"`
	Address2     string    `json:",omitempty" db:"address2"`
	Zip          string    `json:",omitempty" db:"zip`
	State        string    `json:",omitempty" db:"state"`
	Country      string    `json:",omitempty" db:"country"`
	HomePhone    string    `json:"home_phone,omitempty" db:"home_phone"`
	CellPhone    string    `json:"cell_phone,omitempty" db:"cell_phone"`
	EmailAddress string    `json:"email_address,omitempty" db:"email_address"`
	Source       string    `json:"source,omitempty" db:"source"`
	Created      time.Time `json:"created,omitempty" db:"created"`
}

type nullPerson struct {
	Id           sql.NullString `json:",omitempty" db:"id"`
	FirstName    sql.NullString `json:"first_name,omitempty" db:"first_name"`
	LastName     sql.NullString `json:"last_name,omitempty" db:"last_name"`
	MiddleName   sql.NullString `json:"middle_name,omitempty" db:"middle_name"`
	Suffix       sql.NullString `json:",omitempty" db:"suffix"`
	Prefix       sql.NullString `json:",omitempty" db:"prefix"`
	Title        sql.NullString `json:",omitempty" db:"title"`
	Address1     sql.NullString `json:",omitempty" db:"address1"`
	Address2     sql.NullString `json:",omitempty" db:"address2"`
	Zip          sql.NullString `json:",omitempty" db:"zip`
	State        sql.NullString `json:",omitempty" db:"state"`
	Country      sql.NullString `json:",omitempty" db:"country"`
	HomePhone    sql.NullString `json:"home_phone,omitempty" db:"home_phone"`
	CellPhone    sql.NullString `json:"cell_phone,omitempty" db:"cell_phone"`
	EmailAddress sql.NullString `json:"email_address,omitempty" db:"email_address"`
	Source       sql.NullString `json:"source,omitempty" db:"source"`
	Created      time.Time      `json:"created,omitempty" db:"created"`
}

func (p *Person) Save(db *sqlx.DB) error {
	if p.Id == "" {
		//New Record
		//SCAN FOR ID
		rows, err := db.NamedQuery("INSERT INTO people (first_name, last_name, middle_name, suffix, prefix, title, home_phone, cell_phone, source, email_address, address1, address2, zip, state, country) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)", p)
		if err != nil {
			return err
		}
		defer rows.Close()
		if rows.Next() {
			rows.Scan(&p.Id)
		}
	}
	_, err := db.NamedQuery("UPDATE people SET first_name=:first_name, last_name=:last_name, middle_name=:middle_name, suffix=:suffix, prefix=:prefix, title=:title, home_phone=:home_phone, cell_phone=:cell_phone, source=:source, email_address=:email_address, address1=:address1, address2=:address2, zip=$13, state=$14, country=$15 WHERE id=$16", p)
	if err != nil {
		return err
	}
	return nil
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
	out := translateNulls(c)
	out.Created = c.Created
	return out, nil
}

func translateNulls(c nullPerson) *Person {
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
	return out
}
