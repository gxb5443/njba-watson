package tokens

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

var ErrTokenNotValid = errors.New("Token not valid")
var ErrNoTokensFound = errors.New("No Tokens found")

type RefreshToken struct {
	Token   string    `db:"id" json:"id"`
	UserId  string    `db:"user_id" json:"user_id"`
	Active  bool      `db:"active" "json:"active"`
	Created time.Time `db:"created" json:"created"`
}

//New accepts a UserId and generates a RefreshToken
func New(db *sqlx.DB, uid string) (*RefreshToken, error) {
	rt := new(RefreshToken)
	rt.Active = true
	rt.UserId = uid
	err := rt.Save(db)
	if err != nil {
		return nil, err
	}
	return rt, nil
}

//Save saves developer to datastore.  Returns error if something went wrong.
//If it is a new user, it will insert into database and
//update User with proper ID
func (u *RefreshToken) Save(db *sqlx.DB) error {
	tx := db.MustBegin()
	defer tx.Commit()
	if u.Token == "" {
		u.Created = time.Now()
		tx.MustExec("INSERT INTO refresh_tokens(id, user_id) VALUES ($2, $2)", u.Token, u.UserId)
		return nil
	}
	tx.MustExec("UPDATE refresh_tokens SET user_id=$1, active=$2 WHERE id=$3", u.UserId, u.Active, u.Token)
	return nil
}
