//Package users provides a model for working with Namely Developer objects
package credentials

import (
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	"devportal/users"

	"github.com/dchest/uniuri"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var ErrTokenNotFound = errors.New("Token not found")

type Credential struct {
	Id                 string    `db:"id" json:"id"`
	Username           string    `db:"username" json:"username"`
	PasswordHash       string    `db:"password_hash" json:"password"`
	PasswordResetToken string    `db:"password_reset_token" json:"password_reset_token"`
	UserId             string    `db:"user_id" "json:"-"`
	Active             bool      `db:"active" "json:"active"`
	Created            time.Time `db:"created" json:"created"`
	LastUpdate         time.Time `db:"last_updated" json:"last_updated"`
}

type nullCredential struct {
	Id                 sql.NullString `db:"id" json:"id"`
	Username           sql.NullString `db:"username" json:"username"`
	PasswordHash       sql.NullString `db:"password_hash" json:"password"`
	PasswordResetToken sql.NullString `db:"password_reset_token" json:"password_reset_token"`
	UserId             sql.NullString `db:"user_id" "json:"-"`
	Active             bool           `db:"active" "json:"active"`
	Created            time.Time      `db:"created" json:"created"`
	LastUpdate         time.Time      `db:"last_updated" json:"last_updated"`
}

//Save saves developer to datastore.  Returns error if something went wrong.
//If it is a new user, it will insert into database and
//update User with proper ID
func (u *Credential) Save(db *sqlx.DB) error {
	u.Username = strings.ToLower(u.Username)
	if u.Id == "" {
		rows, err := db.NamedQuery("INSERT INTO credentials (username, password_hash, password_reset_token, user_id) VALUES (:username, :password_hash, :password_reset_token, :user_id) RETURNING id", u)
		if err != nil {
			return err
		}
		defer rows.Close()
		if rows.Next() {
			rows.Scan(&u.Id)
		}
	}
	_, err := db.NamedExec("UPDATE credentials SET username=:username, password_hash=:password_hash, password_reset_token=:password_reset_token, user_id=:user_id WHERE id=:id", u)
	if err != nil {
		return err
	}
	return nil
}

//DeleteById finds a user by their UUID and then deletes them.
func DeleteById(db *sqlx.DB, id string) error {
	_, err := db.NamedExec("DELETE FROM credentials where id=:id", id)
	return err
}

//Delete deletes a user from the datastore
func (u *Credential) Delete(db *sqlx.DB) error {
	return DeleteById(db, u.Id)
}

//HashPassword Generates hashed password based on supplied password string
func (u *Credential) HashPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

//Login Checks the database to check if user exists and if the supplied plaintext matches the hashed
//password.  Returns valid user object based on credentials
func Login(db *sqlx.DB, username string, password string) (*users.User, error) {
	var user *users.User
	login, err := GetByUsername(db, strings.ToLower(username))
	if err != nil {
		return nil, err
	}
	if !login.Active {
		return nil, errors.New("Account Disabled. Contact Administrator.")
	}
	if verr := bcrypt.CompareHashAndPassword([]byte(login.PasswordHash), []byte(password)); verr != nil {
		return nil, errors.New("Incorrect Password")
	}
	err = db.Select(&user, "SELECT * FROM users where id=$1", login.UserId)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Could not access users database")
	}
	if user == nil {
		return nil, errors.New("No associated user")
	}
	return user, nil
}

//GenerateAccessCode Generates Invite code
func GenerateAccessCode() string {
	accesscode := uniuri.NewLen(32)
	return accesscode
}

//Get Fetches a single Credential from database by Id
func Get(db *sqlx.DB, id string) (*Credential, error) {
	credential := new(Credential)
	err := db.Select(&credential, "select * from login_credentials where id=$1 LIMIT 1", id)
	return credential, err
}

//GetByUsername Fetches a single Credential from database by Id
func GetByUsername(db *sqlx.DB, username string) (*Credential, error) {
	var login *Credential
	err := db.Select(&login, "select * from credentials where username=$1", strings.ToLower(username))
	if err == sql.ErrNoRows {
		return nil, errors.New("Username not found")
	}
	if err != nil {
		log.Println(err)
		return nil, errors.New("Could not access credential database")
	}
	return login, nil
}

//GetByResetToken Fetches a single Credential from database by token
func GetByResetToken(db *sqlx.DB, token string) (*Credential, error) {
	var login *Credential
	err := db.Select(&login, "select * from credentials where password_reset_token=$1", token)
	if err == sql.ErrNoRows {
		return nil, ErrTokenNotFound
	}
	if err != nil {
		log.Println(err)
		return nil, errors.New("Could not access credential database")
	}
	return login, nil
}

//EmailExists Checks if an email address is available
func UsernameExists(db *sqlx.DB, username string) (bool, error) {
	var exists = false
	err := db.Select(&exists, "select exists(select username from login_credentials where username=$1)", strings.ToLower(username))
	return exists, err
}

//IdExists Checks if a user id exists
func IdExists(db *sqlx.DB, id string) (bool, error) {
	var exists = false
	err := db.Select(&exists, "select exists(select id from users where id=$1)", id)
	return exists, err
}
