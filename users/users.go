//Package users provides a model for working with Namely Developer objects
package users

import (
	"database/sql"
	"errors"
	"time"

	"github.com/dchest/uniuri"
	"github.com/jmoiron/sqlx"
)

type User struct {
	Id        string    `db:"id" json:"id"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" json:"last_name"`
	Email     string    `db:"email" json:"email"`
	Active    bool      `db:"active" json:"active"`
	Admin     bool      `db:"admin" json:"admin"`
	Created   time.Time `db:"created" json:"created"`
}

type nullUser struct {
	Id        sql.NullString `db:"id" json:"id"`
	FirstName sql.NullString `db:"first_name" json:"first_name"`
	LastName  sql.NullString `db:"last_name" json:"last_name"`
	Email     sql.NullString `db:"email" json:"email"`
	Active    bool           `db:"active" json:"active"`
	Admin     bool           `db:"admin" json:"admin"`
	Created   time.Time      `db:"created" json:"created"`
}

//Save saves developer to datastore.  Returns error if something went wrong.
//If it is a new user, it will insert into database and
//update User with proper ID
func (u *User) Save(db *sqlx.DB) error {
	if u.Id == "" {
		rows, err := db.NamedQuery("INSERT INTO users (first_name, last_name, email, admin) VALUES (:first_name, :last_name, :email, :admin) RETURNING id", u)
		if err != nil {
			return errors.New("Could not update given User")
		}
		defer rows.Close()
		if rows.Next() {
			rows.Scan(&u.Id)
		}
		return nil
	}
	_, err := db.NamedExec("UPDATE users SET first_name=:first_name, last_name=:last_name, email=:email, admin=:admin, active=:active WHERE id=:id", u)
	if err != nil {
		return err
	}
	return nil
}

//DeleteById finds a user by their UUID and then deletes them.
func DeleteById(db *sqlx.DB, id string) error {
	_, err := db.NamedExec("DELETE FROM users where id=:id", id)
	return err
}

//Delete deletes a user from the datastore
func Delete(db *sqlx.DB, user *User) error {
	return DeleteById(db, user.Id)
}

//GenerateAccessCode Generates Invite code
func GenerateAccessCode() string {
	accesscode := uniuri.NewLen(32)
	return accesscode
}

//GetAlL fetches all Users from database
func GetAll(db *sqlx.DB) ([]*User, error) {
	query := `SELECT * FROM users;`
	var u []*User
	err := db.Select(&u, query)
	if err != nil {
		return nil, err
	}
	return u, nil
}

//Get Fetches a single User from database by Id
func Get(db *sqlx.DB, id string) (*User, error) {
	query := `SELECT * FROM users WHERE id=$1 LIMIT 1`
	var c nullUser
	err := db.Get(&c, query, id)
	if err != nil {
		return nil, err
	}
	out := translateNulls(c)
	out.Created = c.Created
	return out, nil
}

//EmailExists Checks if an email address is available
func EmailExists(db *sqlx.DB, email string) (bool, error) {
	var exists = false
	err := db.Select(&exists, "select exists(select email from users where email=$1)", email)
	return exists, err
}

//IdExists Checks if a user id exists
func IdExists(db *sqlx.DB, id string) (bool, error) {
	var exists = false
	err := db.Select(&exists, "select exists(select id from users where id=$1)", id)
	return exists, err
}

func translateNulls(c nullUser) *User {
	out := new(User)
	out.Created = c.Created
	out.Active = c.Active
	out.Admin = c.Admin
	if c.Id.Valid {
		out.Id = c.Id.String
	}
	if c.FirstName.Valid {
		out.FirstName = c.FirstName.String
	}
	if c.LastName.Valid {
		out.LastName = c.LastName.String
	}
	if c.Email.Valid {
		out.Email = c.Email.String
	}
	return out
}
