package handlers

import (
	"fmt"
	"log"
	"time"

	"../companies"
	"../memberships"
	"../people"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

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
	query := fmt.Sprintf(`
		SELECT id, name, address1, address2, zip, city, state, created from companies where id='%s'
	`, cid)
	var cos []struct {
		Id       string
		Name     string
		Address1 string
		Address2 string
		Zip      string
		City     string
		State    string
		Created  time.Time
	}
	err := db.Select(&cos, query)
	if err != nil {
		log.Printf("GetMemberships: %s", err)
		return
	}
	c.JSON(200, cos)
}
