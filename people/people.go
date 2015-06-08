package people

import (
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

func GetAll(db *sqlx.DB) ([]*Person, error) {
	query := ` SELECT id, first_name, last_name FROM people;	`
	var p []*Person
	err := db.Select(&p, query)
	if err != nil {
		return nil, err
	}
	return p, nil
}
