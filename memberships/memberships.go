package memberships

import (
	"time"

	"github.com/gxb5443/Cuddy/companies"

	"github.com/jmoiron/sqlx"
)

type Membership struct {
	MemberId int    `json:",omitempty"`
	Status   string `json:",omitempty"`
	NAHBId   int    `json:", omitempty"`
	LocalId  int    `json:", omitempty"`
	*companies.Company
	Created time.Time `json:",omitempty"`
}

func genNewMemebership(db *sqlx.DB) (*Membership, error) {
	query := `INSERT INTO membership DEFAULT VALUES RETURNING id, status, created;`
	var m []*Membership
	db.Select(&m, query)
	return m[0], nil
}

func GenerateNewMembership(db *sqlx.DB) error {
	//Generates a fresh membership in the database
	return nil
}

func (m *Membership) Save(db *sqlx.DB) error {
	temp := []struct {
		id      int
		status  string
		created time.Time
	}{}
	tx := db.MustBegin()
	defer tx.Commit()
	if m.MemberId == 0 {
		//New Membership Record
		query := `INSERT INTO membership DEFAULT VALUES RETURNING id, status, created;`
		db.Select(&temp, query)
		//If Company Information is present, Connect membership to company
		query = `INSERT INTO membership DEFAULT VALUES RETURNING id, status, created;`
		return nil
	}
	return nil
}

func GetAll(db *sqlx.DB) ([]*Membership, error) {
	query := `
		SELECT companies.name as name, companies.id as id, membership_holdings.created as created, membership_holdings.member_id as MemberId from membership_holdings inner join companies on membership_holdings.company_id=companies.id;
	`
	var m []*Membership
	err := db.Select(&m, query)
	if err != nil {
		return nil, err
	}
	return m, nil
}
