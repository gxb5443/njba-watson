package companies

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type Company struct {
	Id       string    `json:",omitempty" db:"id"`
	Name     string    `json:",omitempty" db:"name"`
	Address1 string    `json:",omitempty" db:"address1"`
	Address2 string    `json:",omitempty" db:"address2"`
	Zip      string    `json:",omitempty" db:"zip"`
	State    string    `json:",omitempty" db:"state"`
	Country  string    `json:",omitempty" db:"country"`
	Created  time.Time `json:",omitempty" db:"created"`
}

type nullCo struct {
	Id       sql.NullString `json:",omitempty" db:"id"`
	Name     sql.NullString `json:",omitempty" db:"name"`
	Address1 sql.NullString `json:",omitempty" db:"address1"`
	Address2 sql.NullString `json:",omitempty" db:"address2"`
	Zip      sql.NullString `json:",omitempty" db:"zip"`
	State    sql.NullString `json:",omitempty" db:"state"`
	Country  sql.NullString `json:",omitempty" db:"country"`
	Created  time.Time      `json:",omitempty" db:"created"`
}

func (co *Company) Save(db *sqlx.DB) error {
	tx := db.MustBegin()
	defer tx.Commit()
	if co.Id == "" {
		//New Record
		tx.MustExec("INSERT INTO companies (name, address1, address2, zip, state, country) VALUES ($1,$2,$3,$4,$5,$6)", co.Name, co.Address1, co.Address2, co.Zip, co.State, co.Country)
		return nil
	}
	tx.MustExec("UPDATE companies SET name=$1, address1=$2, address2=$3, zip=$4, state=$5, country=$6 WHERE id=$7", co.Name, co.Address1, co.Address2, co.Zip, co.State, co.Country, co.Id)
	return nil
}

func GetAll(db *sqlx.DB) ([]*Company, error) {
	query := ` SELECT id, name FROM companies;	`
	var c []*Company
	err := db.Select(&c, query)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func GetById(db *sqlx.DB, id string) (*Company, error) {
	query := `SELECT id, name, address1, address2, zip, state, country, created FROM companies WHERE id=$1 LIMIT 1`
	var c nullCo
	err := db.Get(&c, query, id)
	if err != nil {
		return nil, err
	}
	out := new(Company)
	if c.Id.Valid {
		out.Id = c.Id.String
	}
	if c.Name.Valid {
		out.Name = c.Name.String
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
