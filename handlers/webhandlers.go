package handlers

import (
	"errors"
	"log"

	"github.com/gxb5443/Cuddy/credentials"
	"github.com/gxb5443/Cuddy/utils"

	"github.com/gxb5443/Cuddy/companies"
	"github.com/gxb5443/Cuddy/locals"
	"github.com/gxb5443/Cuddy/memberships"
	"github.com/gxb5443/Cuddy/people"
	"github.com/gxb5443/Cuddy/users"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

const (
	APP_SECRET_LENGTH   = 30
	MIN_PASSWORD_LENGTH = 10
)

func GetUsers(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)
	users, err := users.GetAll(db)
	if err != nil {
		log.Printf("GetUsers: %s", err)
		return
	}
	c.JSON(200, users)
}

func GetPeople(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)
	people, err := people.GetAll(db)
	if err != nil {
		log.Printf("GetPeople: %s", err)
		return
	}
	c.JSON(200, people)
}

func GetCompanies(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)
	cos, err := companies.GetAll(db)
	if err != nil {
		log.Printf("GetCompanies: %s", err)
		return
	}
	c.JSON(200, cos)
}

func GetLocals(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)
	los, err := locals.GetAll(db)
	if err != nil {
		log.Printf("GetCompanies: %s", err)
		return
	}
	c.JSON(200, los)
}

func GetMemberships(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)
	mem, err := memberships.GetAll(db)
	if err != nil {
		log.Printf("GetMemberships: %s", err)
		return
	}
	c.JSON(200, mem)
}

func GetCompanyById(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)
	cid := c.Param("cid")
	//Add POC here
	cos, err := companies.GetById(db, cid)
	if err != nil {
		log.Printf("GetCompanyById: %s", err)
		return
	}
	c.JSON(200, cos)
}

func GetPersonById(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)
	pid := c.Param("pid")
	//Add POC here
	p, err := people.GetById(db, pid)
	if err != nil {
		log.Printf("GetPersonById: %s", err)
		return
	}
	c.JSON(200, p)
}

func CreateCompany(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)
	var co companies.Company
	c.Bind(&co)
	err := co.Save(db)
	if err != nil {
		log.Printf("CreateCompany: %s", err)
		return
	}
}

func AddUser(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)
	var u users.User
	var login credentials.Credential
	c.Bind(&u)
	if exists, _ := users.EmailExists(db, u.Email); exists {
		c.JSON(409, gin.H{"status": "User already exists"})
		return
	}
	if u.Email == "" {
		c.JSON(412, gin.H{"status": "Email is required"})
		return
	}
	serr := u.Save(db)
	if serr != nil {
		log.Println(serr)
		c.AbortWithError(500, errors.New("User could not be created"))
		return
	}
	login.Active = true
	login.Username = u.Email
	login.UserId = u.Id
	password := utils.GeneratePassword(MIN_PASSWORD_LENGTH)
	log.Println("Password created: ", password)
	perr := login.HashPassword(password)
	if perr != nil {
		log.Println(perr)
		c.AbortWithError(500, errors.New("Password could not be created"))
		return
	}
	login.PasswordResetToken, perr = utils.GenerateRandomString(APP_SECRET_LENGTH)
	if perr != nil {
		log.Println(perr)
		c.AbortWithError(500, errors.New("Could not Generate token"))
		return
	}
	log.Println("Password reset token: ", login.PasswordResetToken)
	serr = login.Save(db)
	if serr != nil {
		log.Println(serr)
		c.AbortWithError(500, errors.New("User credentials could not be created"))
		return
	}
	log.Println("Created User: ", u.Email, " with password: ", password)
	//Add Emailer here to email user with link
	c.JSON(201, u)
}
