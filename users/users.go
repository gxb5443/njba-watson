//Package users provides a model for working with Namely Developer objects
package users

import (
	"errors"
	"time"

	"github.com/coopernurse/gorp"
	"github.com/dchest/uniuri"
)

type User struct {
	Id        string    `db:"id" json:"id"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" json:"last_name"`
	Email     string    `db:"email" json:"email"`
	Active    bool      `db:"active" json:"active"`
	Admin     bool      `db:"admin" json:"admin"`
	Companies string    `db:"companies" json:"companies"`
	Created   time.Time `db:"created" json:"created"`
}

//Save saves developer to datastore.  Returns error if something went wrong.
//If it is a new user, it will insert into database and
//update User with proper ID
func (u *User) Save(dbMap *gorp.DbMap) error {
	if u.Id != "" {
		_, err := dbMap.Update(u)
		if err != nil {
			return errors.New("Could not update given User")
		}
	} else {
		u.Created = time.Now()
		err := dbMap.Insert(u)
		if err != nil {
			return err
			//return errors.New("Could not add new user")
		}
	}
	return nil
}

//DeleteById finds a user by their UUID and then deletes them.
func DeleteById(dbMap *gorp.DbMap, id string) error {
	u, err := Get(dbMap, id)
	if err != nil {
		return err
	}
	if u == nil {
		return nil
	}
	err = Delete(dbMap, u)
	return err
}

//Delete deletes a user from the datastore
func Delete(dbMap *gorp.DbMap, user *User) error {
	_, err := dbMap.Delete(user)
	return err
}

//GenerateAccessCode Generates Invite code
func GenerateAccessCode() string {
	accesscode := uniuri.NewLen(32)
	return accesscode
}

//GetAlL fetches all Users from database
func GetAll(dbMap *gorp.DbMap) ([]*User, error) {
	objs, err := dbMap.Select(User{}, "select * from users")
	users := make([]*User, len(objs))
	for i, u := range objs {
		users[i] = u.(*User)
	}
	return users, err
}

//Get Fetches a single User from database by Id
func Get(dbMap *gorp.DbMap, id string) (*User, error) {
	//obj, err := dbMap.Get(User{}, id)
	user := new(User)
	err := dbMap.SelectOne(&user, "select * from users where id=$1 LIMIT 1", id)
	return user, err
}

//EmailExists Checks if an email address is available
func EmailExists(dbMap *gorp.DbMap, email string) (bool, error) {
	var exists = false
	err := dbMap.SelectOne(&exists, "select exists(select email from users where email=$1)", email)
	return exists, err
}

//IdExists Checks if a user id exists
func IdExists(dbMap *gorp.DbMap, id string) (bool, error) {
	var exists = false
	err := dbMap.SelectOne(&exists, "select exists(select id from users where id=$1)", id)
	return exists, err
}
