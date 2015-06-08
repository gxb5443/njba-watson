package companies

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Company struct {
	Id       string    `json:",omitempty"`
	Name     string    `json:",omitempty" db:"name"`
	Address1 string    `json:",omitempty"`
	Address2 string    `json:",omitempty"`
	Zip      string    `json:",omitempty"`
	State    string    `json:",omitempty"`
	Country  string    `json:",omitempty"`
	Created  time.Time `json:"source,omitempty"`
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
