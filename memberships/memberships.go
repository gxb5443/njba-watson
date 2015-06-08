package memberships

import (
	"time"

	"../companies"

	"github.com/jmoiron/sqlx"
)

type Membership struct {
	MemberId string `json:",omitempty"`
	*companies.Company
	Created time.Time `json:",omitempty"`
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
