package users

import (
	"database/sql"
	"log"
	"testing"

	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var saveTestCases = []*User{
	{
		FirstName: "Gian",
		LastName:  "Biondi",
		Email:     "gian@namely.com",
	},
	{
		FirstName: "Andrew",
		LastName:  "Danker",
		Email:     "adank@hs.com",
	},
	{
		FirstName: "Blach",
		LastName:  "Bach",
		Email:     "test@test.com",
	},
}

const DBOPEN = "user=ubuntu dbname=circle_test sslmode=disable"

//const DBOPEN = "user=manuel dbname=test_developers sslmode=disable"

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func initDb(config string) *gorp.DbMap {
	db, err := sql.Open("postgres", config)
	if err != nil {
		panic(err)
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbmap.AddTableWithName(User{}, "users").SetKeys(true, "Id")
	return dbmap
}

func TestSave(t *testing.T) {
	db, derr := sql.Open("postgres", DBOPEN)
	if derr != nil {
		log.Fatalln("Could not truncate db")
	}
	db.Exec("Truncate users cascade;")
	db.Close()
	dbmap := initDb(DBOPEN)
	defer dbmap.Db.Close()
	for _, test := range saveTestCases {
		err := test.Save(dbmap)
		if err != nil {
			log.Fatalln("SaveTestCases: ", err)
		}
		assert.NotEqual(t, test.Id, "", "ID should be set.")
	}
	var newNames = []string{"Timmy", "Tommy", "Jonny"}
	for i, test := range saveTestCases {
		test.FirstName = newNames[i]
		err := test.Save(dbmap)
		if err != nil {
			log.Fatalln("SaveTestCases: ", err)
		}
		require.Equal(t, test.FirstName, newNames[i], "First name should have changed.")
		saveTestCases[i] = test
	}
}

func TestGetAll(t *testing.T) {
	dbmap := initDb(DBOPEN)
	defer dbmap.Db.Close()
	u, err := GetAll(dbmap)
	if err != nil {
		log.Fatalln("TestGetall: ", err)
	}
	assert.NotEmpty(t, u, "Should have gotten all Users back")
	for i, user := range u {
		assert.Equal(t, user.Id, saveTestCases[i].Id, "GetAll does not match expected")
	}
}

func TestGet(t *testing.T) {
	dbmap := initDb(DBOPEN)
	defer dbmap.Db.Close()
	for _, user := range saveTestCases {
		u, err := Get(dbmap, user.Id)
		if err != nil {
			log.Fatalln("TestGet: ", err)
		}
		if assert.NotNil(t, u) {
			assert.Equal(t, u.FirstName, user.FirstName, "Get value is not as expected")
		}
	}
	//Unknown Id's for Get
	var badIds = []string{"624f6dd0-91f2-4026-a684-01924da4be84", "624f6dd0-91f2-4026-a684-01924da4be25"}

	for _, id := range badIds {
		_, err := Get(dbmap, id)
		if err != sql.ErrNoRows {
			assert.Equal(t, true, false, "Nothing should be returned here")
		}
		assert.NotNil(t, err, "Should throw an error when an invalid UUID is passed")
	}

	//Invalid Id's for Get
	var invalidIds = []string{"Timmy", ""}

	for _, id := range invalidIds {
		_, err := Get(dbmap, id)
		//assert.Nil(t, u, "Nothing should be returned here")
		assert.NotNil(t, err, "Should throw an error when an invalid UUID is passed")
	}
}

func TestEmailExists(t *testing.T) {
	dbmap := initDb(DBOPEN)
	defer dbmap.Db.Close()
	exists, err := EmailExists(dbmap, saveTestCases[0].Email)
	if err != nil {
		log.Fatalln("TestEmailExists: ", err)
	}
	assert.Equal(t, true, exists, "The emails SHOULD exist..")
	exists, err = EmailExists(dbmap, "NotRealEmail@real.not")
	if err != nil {
		log.Fatalln("TestEmailExists: ", err)
	}
	assert.Equal(t, false, exists, "This email should not exist..")
}

func TestDeleteById(t *testing.T) {
	dbmap := initDb(DBOPEN)
	defer dbmap.Db.Close()
	exists, err := IdExists(dbmap, saveTestCases[0].Id)
	if err != nil {
		log.Fatalln("TestEmailExists: ", err)
	}
	err = DeleteById(dbmap, saveTestCases[0].Id)
	exists, err = IdExists(dbmap, saveTestCases[0].Id)
	assert.Equal(t, false, exists, "The user should not exist..")
}
